package json

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/goccy/go-json"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type frameKind uint8

const (
	frameObject frameKind = iota + 1
	frameArray
)

type (
	decoder struct {
		pre  []encoding.PreDecoderHook
		post []encoding.PostDecoderHook
	}

	decodeFrame struct {
		obj       *runtime.Object
		arr       *runtime.Array
		key       string
		kind      frameKind
		expectKey bool
	}
)

func (dec decoder) decodeValue(ctx context.Context, jdec *json.Decoder) (runtime.Value, error) {
	var stack []decodeFrame
	var root runtime.Value

	attach := func(val runtime.Value) error {
		if len(stack) == 0 {
			if root != nil {
				return fmt.Errorf("json: multiple root values")
			}

			root = val

			return nil
		}

		top := &stack[len(stack)-1]

		switch top.kind {
		case frameArray:
			return top.arr.Append(ctx, val)
		case frameObject:
			if top.expectKey {
				return fmt.Errorf("json: expected object key before value")
			}

			if err := top.obj.Set(ctx, runtime.NewString(top.key), val); err != nil {
				return err
			}

			top.expectKey = true
			top.key = ""
		default:
			return fmt.Errorf("json: invalid frame kind")
		}

		return nil
	}

	for {
		token, err := jdec.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return runtime.None, err
		}

		switch v := token.(type) {
		case json.Delim:
			switch v {
			case '{':
				stack = append(stack, decodeFrame{kind: frameObject, obj: runtime.NewObject(), expectKey: true})
			case '[':
				stack = append(stack, decodeFrame{kind: frameArray, arr: runtime.NewArray(0)})
			case '}':
				if len(stack) == 0 || stack[len(stack)-1].kind != frameObject {
					return runtime.None, fmt.Errorf("json: unexpected object end")
				}

				top := stack[len(stack)-1]
				if !top.expectKey {
					return runtime.None, fmt.Errorf("json: missing value for key %q", top.key)
				}

				stack = stack[:len(stack)-1]

				if err := attach(top.obj); err != nil {
					return runtime.None, err
				}
			case ']':
				if len(stack) == 0 || stack[len(stack)-1].kind != frameArray {
					return runtime.None, fmt.Errorf("json: unexpected array end")
				}

				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if err := attach(top.arr); err != nil {
					return runtime.None, err
				}
			default:
				return runtime.None, fmt.Errorf("json: unsupported delimiter %q", v)
			}
		case string:
			if len(stack) > 0 && stack[len(stack)-1].kind == frameObject && stack[len(stack)-1].expectKey {
				stack[len(stack)-1].key = v
				stack[len(stack)-1].expectKey = false

				continue
			}

			if err := attach(runtime.NewString(v)); err != nil {
				return runtime.None, err
			}
		case json.Number:
			raw := v.String()

			if !strings.ContainsAny(raw, ".eE") {
				parsed, err := strconv.ParseInt(raw, 10, 0)

				if err == nil {
					if err := attach(runtime.NewInt(int(parsed))); err != nil {
						return runtime.None, err
					}

					break
				}
			}

			f, err := v.Float64()
			if err != nil {
				return runtime.None, err
			}

			if err := attach(runtime.NewFloat(f)); err != nil {
				return runtime.None, err
			}
		case float64:
			if err := attach(runtime.NewFloat(v)); err != nil {
				return runtime.None, err
			}
		case bool:
			if err := attach(runtime.NewBoolean(v)); err != nil {
				return runtime.None, err
			}
		case nil:
			if err := attach(runtime.None); err != nil {
				return runtime.None, err
			}
		default:
			return runtime.None, fmt.Errorf("json: unsupported token %T", v)
		}
	}

	if root == nil {
		return runtime.None, fmt.Errorf("json: empty input")
	}

	return root, nil
}

func (dec decoder) runPreHooks(data []byte) error {
	if len(dec.pre) == 0 {
		return nil
	}

	for _, hook := range dec.pre {
		if err := hook(data); err != nil {
			return err
		}
	}

	return nil
}

func (dec decoder) runPostHooks(data []byte, err error) error {
	if len(dec.post) == 0 {
		return nil
	}

	for _, hook := range dec.post {
		if hookErr := hook(data, err); hookErr != nil {
			return hookErr
		}
	}

	return nil
}
