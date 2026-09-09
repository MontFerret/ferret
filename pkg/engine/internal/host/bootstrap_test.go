package host

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestBootstrapBuildSnapshotsHostRegistries(t *testing.T) {
	resources := resource.NewManager()
	t.Cleanup(func() { _ = resources.Close() })
	params := runtime.Params{"value": runtime.Int(1)}
	codecs := encoding.NewRegistry(encodingjson.Default)
	hooks := NewHooks()
	boot, err := NewBootstrap(Config{
		Library: runtime.NewLibrary(), Params: params, Encoding: codecs,
		FSRoot: t.TempDir(), FSReadOnly: true,
	}, hooks, resources)
	if err != nil {
		t.Fatal(err)
	}

	called := false
	boot.Hooks().Engine().OnInit(func() error {
		called = true

		return nil
	})
	boot.Host().Params()["value"] = runtime.Int(2)
	host, err := boot.Build()
	if err != nil {
		t.Fatal(err)
	}

	if called {
		t.Fatal("building the host ran initialization hooks")
	}

	params["value"] = runtime.Int(3)
	if err := boot.Host().Encoding().Register(encodingmsgpack.Default); err != nil {
		t.Fatal(err)
	}

	if host.Params["value"] != runtime.Int(2) {
		t.Fatal("host parameters share the bootstrap map")
	}

	host.Params["value"] = runtime.Int(4)
	if params["value"] != runtime.Int(3) {
		t.Fatal("host parameter mutation changed bootstrap state")
	}

	if _, err := host.Encoding.Codec(encodingmsgpack.ContentType); !errors.Is(err, encoding.ErrCodecNotFound) {
		t.Fatalf("host encoding registry shares bootstrap registrations: %v", err)
	}

	if host.FileSystem != boot.Host().FileSystem() || host.Network != boot.Host().Network() || !host.FSReadOnly {
		t.Fatal("building the host changed its services or filesystem policy")
	}

	if err := hooks.EngineHooks.RunInit(); err != nil || !called {
		t.Fatalf("bootstrap hooks lost the supplied registry: called=%t error=%v", called, err)
	}
}
