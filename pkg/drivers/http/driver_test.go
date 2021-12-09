package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"unsafe"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/text/encoding/charmap"

	"github.com/MontFerret/ferret/pkg/drivers"
)

func Test_newHTTPClientWithTransport(t *testing.T) {
	httpTransport := (http.DefaultTransport).(*http.Transport)
	httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	type args struct {
		options *Options
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "check transport exist with pester.New()",
			args: args{options: &Options{
				Options: &drivers.Options{
					Proxy: "http://0.0.0.|",
				},
				HTTPTransport: httpTransport,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Convey(tt.name, t, func() {
				var (
					transport *http.Transport
					client    = newHTTPClient(tt.args.options)
					rValue    = reflect.ValueOf(client).Elem()
					rField    = rValue.Field(0)
				)

				rField = reflect.NewAt(rField.Type(), unsafe.Pointer(rField.UnsafeAddr())).Elem()
				hc := rField.Interface().(*http.Client)

				if hc != nil {
					transport = hc.Transport.(*http.Transport)
				} else {
					transport = client.Transport.(*http.Transport)
				}

				verify := transport.TLSClientConfig.InsecureSkipVerify

				So(verify, ShouldBeTrue)
			})
		})
	}
}

func Test_newHTTPClient(t *testing.T) {

	Convey("pester.New()", t, func() {
		var (
			client = newHTTPClient(&Options{
				Options: &drivers.Options{
					Proxy: "http://0.0.0.|",
				},
			})

			rValue = reflect.ValueOf(client).Elem()
			rField = rValue.Field(0)
		)

		rField = reflect.NewAt(rField.Type(), unsafe.Pointer(rField.UnsafeAddr())).Elem()
		hc := rField.Interface().(*http.Client)

		So(hc, ShouldBeNil)
	})
}

func TestDriver_convertToUTF8(t *testing.T) {
	type args struct {
		inputData  string
		srcCharset string
	}
	tests := []struct {
		name     string
		args     args
		wantData io.Reader
		expected string
		wantErr  bool
	}{
		{
			name: "should convert to expected state",
			args: args{
				inputData:  `<!DOCTYPE html><html><head><meta charset="windows-1251"/></head><body>феррет</body></html>`,
				srcCharset: "windows-1251",
			},
			wantErr:  false,
			expected: `<!DOCTYPE html><html><head><meta charset="windows-1251"/></head><body>феррет</body></html>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drv := &Driver{}

			Convey(tt.name, t, func() {

				data, err := ioutil.ReadAll(bytes.NewBufferString(tt.args.inputData))
				if err != nil {
					panic(err)
				}

				encodedData := make([]byte, len(data)*2)

				dec := charmap.Windows1251.NewEncoder()
				nDst, _, err := dec.Transform(encodedData, data, false)
				if err != nil {
					panic(err)
				}

				encodedData = encodedData[:nDst]

				gotData, err := drv.convertToUTF8(bytes.NewReader(encodedData), tt.args.srcCharset)
				So(err, ShouldBeNil)

				outData, err := ioutil.ReadAll(gotData)
				So(err, ShouldBeNil)

				So(string(outData), ShouldEqual, tt.expected)

			})

		})
	}
}

func TestDriver_Concurrency(t *testing.T) {
	Convey("Should make only 1 request", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "http://localhost:1111",
			httpmock.NewStringResponder(200, `<!DOCTYPE html><html><head></head><body></body></html>`))

		drv := NewDriver()

		page, err := drv.Open(context.Background(), drivers.Params{
			URL: "http://localhost:1111",
		})

		So(err, ShouldBeNil)
		So(page, ShouldNotBeNil)
		So(httpmock.GetTotalCallCount(), ShouldEqual, 1)
	})
}
