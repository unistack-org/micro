package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unistack-org/micro/v3/registry"
)

func TestStoreRegex(t *testing.T) {
	t.Skip()
	router, err := newRouter()
	if err != nil {
		t.Fatal(err)
	}
	router.store([]*registry.Service{
		{
			Name:    "Foobar",
			Version: "latest",
			Endpoints: []*registry.Endpoint{
				{
					Name: "foo",
					Metadata: map[string]string{
						"endpoint":    "FooEndpoint",
						"description": "Some description",
						"method":      "POST",
						"path":        "^/foo/$",
						"handler":     "rpc",
					},
				},
			},
			Metadata: map[string]string{},
		},
	},
	)

	assert.Len(t, router.ceps["Foobar.foo"].pcreregs, 1)
}
