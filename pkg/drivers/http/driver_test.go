package http

import (
	"crypto/tls"
	"net/http"
	"reflect"
	"testing"
	"unsafe"

	"github.com/smartystreets/goconvey/convey"
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
				Proxy:         "http://0.0.0.|",
				HTTPTransport: httpTransport,
			}},
		},
		{
			name: "check transport exist with pester.NewExtendedClient()",
			args: args{options: &Options{
				Proxy:         "http://0.0.0.0",
				HTTPTransport: httpTransport,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			convey.Convey(tt.name, t, func() {
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

				convey.So(verify, convey.ShouldBeTrue)
			})
		})
	}
}

func Test_newHTTPClient(t *testing.T) {

	convey.Convey("pester.New()", t, func() {
		var (
			client = newHTTPClient(&Options{
				Proxy: "http://0.0.0.|",
			})

			rValue = reflect.ValueOf(client).Elem()
			rField = rValue.Field(0)
		)

		rField = reflect.NewAt(rField.Type(), unsafe.Pointer(rField.UnsafeAddr())).Elem()
		hc := rField.Interface().(*http.Client)

		convey.So(hc, convey.ShouldBeNil)
	})

	convey.Convey("pester.NewExtend()", t, func() {
		var (
			client = newHTTPClient(&Options{
				Proxy: "http://0.0.0.0",
			})

			rValue = reflect.ValueOf(client).Elem()
			rField = rValue.Field(0)
		)

		rField = reflect.NewAt(rField.Type(), unsafe.Pointer(rField.UnsafeAddr())).Elem()
		hc := rField.Interface().(*http.Client)

		convey.So(hc, convey.ShouldNotBeNil)
	})

}
