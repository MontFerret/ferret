package encoding_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	failingWriter struct {
		failure error
		bytes.Buffer
		limit  int
		calls  int
		closed bool
	}
	chunkReader struct {
		failure error
		cancel  context.CancelFunc
		data    []byte
		chunk   int
		calls   int
		closed  bool
	}
	cancelWriter struct {
		cancel context.CancelFunc
		bytes.Buffer
	}
	contextList struct {
		runtime.List
		expected context.Context
		seen     bool
	}
)

func (w *failingWriter) Write(data []byte) (int, error) {
	w.calls++
	size := min(len(data), max(0, w.limit-w.Len()))
	n, _ := w.Buffer.Write(data[:size])
	if size < len(data) {
		return n, w.failure
	}
	return n, nil
}
func (w *failingWriter) Close() error { w.closed = true; return nil }
func (r *chunkReader) Read(dst []byte) (int, error) {
	r.calls++
	if r.cancel != nil {
		r.cancel()
		r.cancel = nil
	}
	size := min(len(dst), r.chunk, len(r.data))
	copy(dst, r.data[:size])
	r.data = r.data[size:]
	if len(r.data) == 0 {
		if r.failure != nil {
			return size, r.failure
		}
		return size, io.EOF
	}
	return size, nil
}
func (r *chunkReader) Close() error { r.closed = true; return nil }
func (w *cancelWriter) Write(data []byte) (int, error) {
	n, err := w.Buffer.Write(data)
	if w.cancel != nil {
		w.cancel()
		w.cancel = nil
	}
	return n, err
}
func (v *contextList) Length(ctx context.Context) (runtime.Int, error) {
	if ctx != v.expected {
		return 0, errors.New("Length lost caller context")
	}
	v.seen = true
	return v.List.Length(ctx)
}
func (v *contextList) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if ctx != v.expected {
		return runtime.None, errors.New("At lost caller context")
	}
	return v.List.At(ctx, index)
}

func TestCodecReaderWriterContract(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			value := runtime.NewArrayWith(runtime.NewInt64(42), runtime.NewString("x\nλ"), runtime.True, runtime.None)
			golden, err := encoding.EncodeBytes(t.Context(), codec, value)
			if err != nil {
				t.Fatal(err)
			}

			t.Run("direct_writer_and_fragmented_readers", func(t *testing.T) {
				writer := &failingWriter{limit: 1 << 20}
				if err := codec.Encode(t.Context(), writer, value); err != nil {
					t.Fatal(err)
				}
				if writer.closed || !bytes.Equal(writer.Bytes(), golden) {
					t.Fatal("writer bytes or ownership changed")
				}
				for _, chunk := range []int{1, 3, 17} {
					reader := &chunkReader{data: bytes.Clone(golden), chunk: chunk}
					decoded, err := codec.Decode(t.Context(), reader)
					if err != nil {
						t.Fatalf("chunk %d: %v", chunk, err)
					}
					equal, err := runtime.EqualValues(t.Context(), value, decoded)
					if err != nil || !bool(equal) || reader.closed {
						t.Fatalf("value or ownership changed: %v", err)
					}
				}
			})

			t.Run("write_failures_and_ancestor_hooks", func(t *testing.T) {
				failure := errors.New("writer failed")
				for _, test := range []struct {
					failure error
					want    error
					limit   int
				}{
					{limit: 3, failure: failure, want: failure},
					{limit: 3, want: io.ErrShortWrite},
					{limit: 0, want: io.ErrShortWrite},
				} {
					writer := &failingWriter{limit: test.limit, failure: test.failure}
					var postErrors []error
					enc := codec.EncodeWith().PostHook(func(_ context.Context, _ runtime.Value, err error) error {
						postErrors = append(postErrors, err)
						return nil
					}).Encoder()
					err := enc.Encode(t.Context(), writer, value)
					if !errors.Is(err, test.want) {
						t.Fatalf("error %v, want %v", err, test.want)
					}
					if writer.closed || !bytes.Equal(writer.Bytes(), golden[:test.limit]) {
						t.Fatal("lost output prefix or closed writer")
					}
					if len(postErrors) == 0 || !errors.Is(postErrors[len(postErrors)-1], test.want) {
						t.Fatal("ancestor post-hook lost write error")
					}
				}
			})

			t.Run("reader_error_with_and_without_data", func(t *testing.T) {
				failure := errors.New("reader failed")
				for _, readerErr := range []error{failure, errors.Join(io.EOF, failure)} {
					for _, data := range [][]byte{nil, golden, golden[:len(golden)/2]} {
						reader := &chunkReader{data: bytes.Clone(data), chunk: len(golden), failure: readerErr}
						hookFailure := errors.New("post-hook failed")
						dec := codec.DecodeWith().PostHook(func(ctx context.Context, value runtime.Value, err error) error {
							if ctx != t.Context() || value != runtime.None || !errors.Is(err, failure) {
								t.Errorf("post-hook lost failure/context: %v", err)
							}
							return hookFailure
						}).Decoder()
						value, err := dec.Decode(t.Context(), reader)
						if value != runtime.None || !errors.Is(err, failure) || !errors.Is(err, hookFailure) || reader.closed {
							t.Fatalf("decode lost error or ownership: %v", err)
						}
					}
				}
			})

			t.Run("configuration_and_context", func(t *testing.T) {
				ctx := context.WithValue(t.Context(), struct{}{}, "marker")
				list := &contextList{List: runtime.NewArrayWith(runtime.NewInt(1)), expected: ctx}
				pre, post := 0, 0
				enc := codec.EncodeWith().PreHook(func(actual context.Context, _ runtime.Value) error {
					if actual != ctx {
						t.Fatal("encoder hook lost context")
					}
					pre++
					return nil
				}).PostHook(func(actual context.Context, _ runtime.Value, err error) error {
					if actual != ctx || err != nil {
						t.Fatalf("encoder post-hook: %v", err)
					}
					post++
					return nil
				}).Encoder()
				data, err := encoding.EncodeBytes(ctx, enc, list)
				if err != nil || !list.seen || pre != 2 || post != 2 {
					t.Fatalf("encoder config/context: %v", err)
				}
				reader := &chunkReader{data: data, chunk: 1}
				dec := codec.DecodeWith().PreHook(func(actual context.Context) error {
					if actual != ctx || reader.calls != 0 {
						t.Fatal("decoder pre-hook consumed input or lost context")
					}
					return nil
				}).PostHook(func(actual context.Context, value runtime.Value, err error) error {
					equal, equalErr := runtime.EqualValues(ctx, list.List, value)
					if actual != ctx || err != nil || equalErr != nil || !equal {
						t.Fatalf("decoder post-hook: %v", err)
					}
					return nil
				}).Decoder()
				if _, err := dec.Decode(ctx, reader); err != nil {
					t.Fatal(err)
				}
				if _, err := encoding.DecodeBytes(ctx, codec, data); err != nil {
					t.Fatal(err)
				}
			})

			t.Run("truncated_and_trailing_documents", func(t *testing.T) {
				for _, data := range [][]byte{golden[:len(golden)-1], append(bytes.Clone(golden), golden...)} {
					reader := &chunkReader{data: data, chunk: 1}
					if value, err := codec.Decode(t.Context(), reader); err == nil || value != runtime.None {
						t.Fatalf("accepted incomplete or multiple roots: value=%v err=%v", value, err)
					}
				}
			})

			t.Run("pre_hook_failure", func(t *testing.T) {
				failure := errors.New("pre-hook failed")
				writer := &failingWriter{limit: 100}
				enc := codec.EncodeWith().PreHook(func(context.Context, runtime.Value) error { return failure }).
					PostHook(func(context.Context, runtime.Value, error) error { t.Fatal("post-hook ran"); return nil }).Encoder()
				if err := enc.Encode(t.Context(), writer, value); !errors.Is(err, failure) || writer.calls != 0 {
					t.Fatal(err)
				}
				reader := &chunkReader{data: golden, chunk: 1}
				dec := codec.DecodeWith().PreHook(func(context.Context) error { return failure }).
					PostHook(func(context.Context, runtime.Value, error) error { t.Fatal("post-hook ran"); return nil }).Decoder()
				if _, err := dec.Decode(t.Context(), reader); !errors.Is(err, failure) || reader.calls != 0 {
					t.Fatal(err)
				}
			})

			t.Run("validation_and_pre_cancellation", func(t *testing.T) {
				if err := codec.Encode(nil, io.Discard, value); !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatal(err)
				}
				if err := codec.Encode(t.Context(), nil, value); !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatal(err)
				}
				if _, err := codec.Decode(nil, bytes.NewReader(golden)); !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatal(err)
				}
				if _, err := codec.Decode(t.Context(), nil); !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatal(err)
				}
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				writer := &failingWriter{limit: 100}
				if err := codec.Encode(ctx, writer, value); !errors.Is(err, context.Canceled) || writer.calls != 0 {
					t.Fatal(err)
				}
				reader := &chunkReader{data: golden, chunk: 1}
				if _, err := codec.Decode(ctx, reader); !errors.Is(err, context.Canceled) || reader.calls != 0 {
					t.Fatal(err)
				}
			})

			t.Run("cancellation_during_large_work", func(t *testing.T) {
				large := runtime.NewRange(1, 100_000)
				ctx, cancel := context.WithCancel(t.Context())
				writer := &cancelWriter{cancel: cancel}
				err := codec.Encode(ctx, writer, large)
				if !errors.Is(err, context.Canceled) || writer.Len() == 0 || writer.Len() > 8192 {
					t.Fatalf("encode cancellation: length=%d err=%v", writer.Len(), err)
				}
				data, err := encoding.EncodeBytes(t.Context(), codec, large)
				if err != nil {
					t.Fatal(err)
				}
				ctx, cancel = context.WithCancel(t.Context())
				reader := &chunkReader{data: data, chunk: 64, cancel: cancel}
				if value, err := codec.Decode(ctx, reader); value != runtime.None || !errors.Is(err, context.Canceled) {
					t.Fatalf("decode cancellation: %v", err)
				}
			})

			t.Run("sampled_traversal_and_payload_cancellation", func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				calls := 0
				enc := codec.EncodeWith().PreHook(func(context.Context, runtime.Value) error {
					calls++
					cancel()
					return nil
				}).Encoder()
				if err := enc.Encode(ctx, io.Discard, runtime.NewRange(0, 1000)); !errors.Is(err, context.Canceled) || calls > 256 {
					t.Fatalf("traversal cancellation: hooks=%d err=%v", calls, err)
				}
				for _, value := range []runtime.Value{runtime.NewString(strings.Repeat("x", 1<<20)), runtime.NewBinary(bytes.Repeat([]byte{255}, 1<<20))} {
					ctx, cancel := context.WithCancel(t.Context())
					writer := &cancelWriter{cancel: cancel}
					if err := codec.Encode(ctx, writer, value); !errors.Is(err, context.Canceled) || writer.Len() > 8192 {
						t.Fatalf("payload cancellation: length=%d err=%v", writer.Len(), err)
					}
				}
			})

			t.Run("reuse_does_not_touch_prior_io", func(t *testing.T) {
				writer := &failingWriter{limit: 1000}
				if err := codec.Encode(t.Context(), writer, value); err != nil {
					t.Fatal(err)
				}
				retained := bytes.Clone(writer.Bytes())
				if _, err := encoding.EncodeBytes(t.Context(), codec, runtime.NewString("different")); err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(retained, writer.Bytes()) {
					t.Fatal("encoder reused previous writer")
				}
				reader := &chunkReader{data: bytes.Clone(golden), chunk: 3}
				decoded, err := codec.Decode(t.Context(), reader)
				if err != nil {
					t.Fatal(err)
				}
				calls := reader.calls
				if _, err := encoding.DecodeBytes(t.Context(), codec, golden); err != nil {
					t.Fatal(err)
				}
				equal, err := runtime.EqualValues(t.Context(), decoded, value)
				if reader.calls != calls || err != nil || !equal {
					t.Fatal("decoder reused I/O or mutated retained value")
				}
			})
		})
	}
}

func TestEncodeBytesPreservesPrefixAndHookErrors(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		preFailure, postFailure := errors.New("nested pre-hook failed"), errors.New("ancestor post-hook failed")
		value := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2))
		enc := codec.EncodeWith().PreHook(func(_ context.Context, value runtime.Value) error {
			if value == runtime.Int(2) {
				return preFailure
			}
			return nil
		}).PostHook(func(_ context.Context, value runtime.Value, err error) error {
			if errors.Is(err, preFailure) {
				return postFailure
			}
			return nil
		}).Encoder()
		data, err := encoding.EncodeBytes(t.Context(), enc, value)
		if len(data) == 0 || !errors.Is(err, preFailure) || !errors.Is(err, postFailure) {
			t.Fatalf("prefix/hooks: %s %v", data, err)
		}
	}
}

func TestJSONReaderRejectsIncompleteDocuments(t *testing.T) {
	for _, text := range []string{"", "{", "[1", "{\"a\":1", "1 [", "1 2", "invalid", "[1,]", "[1 2]"} {
		if _, err := ferretjson.Default.Decode(t.Context(), &chunkReader{data: []byte(text), chunk: 1}); err == nil {
			t.Errorf("accepted %q", text)
		}
	}
}

func TestJSONChunkedEscapingAndBinary(t *testing.T) {
	text := strings.Repeat("aλ\n\x01\\\"", 2000) + "\xff"
	data, err := encoding.EncodeBytes(t.Context(), ferretjson.Default, runtime.NewString(text))
	if err != nil {
		t.Fatal(err)
	}
	value, err := ferretjson.Default.Decode(t.Context(), &chunkReader{data: data, chunk: 1})
	if err != nil || value != runtime.NewString(strings.ReplaceAll(text, "\xff", "\ufffd")) {
		t.Fatalf("chunked string: %v", err)
	}
	binary := bytes.Repeat([]byte{0, 1, 255}, 4001)
	data, err = encoding.EncodeBytes(t.Context(), ferretjson.Default, runtime.NewBinary(binary))
	if err != nil {
		t.Fatal(err)
	}
	// Chunk boundaries retain the standard, padded Base64 representation.
	if string(data) != "\""+base64.StdEncoding.EncodeToString(binary)+"\"" {
		t.Fatal("base64 bytes changed")
	}
}
