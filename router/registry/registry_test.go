// +build ignore

package registry

import (
	"os"
	"testing"

	"github.com/unistack-org/micro/v3/registry/memory"
	"github.com/unistack-org/micro/v3/router"
)

func routerTestSetup() router.Router {
	r := memory.NewRegistry()
	return NewRouter(router.Registry(r))
}

func TestRouterClose(t *testing.T) {
	r := routerTestSetup()

	if err := r.Close(); err != nil {
		t.Errorf("failed to stop router: %v", err)
	}
	if len(os.Getenv("INTEGRATION_TESTS")) == 0 {
		t.Logf("TestRouterStartStop STOPPED")
	}
}
