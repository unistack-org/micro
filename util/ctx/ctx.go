package ctx

import (
	"context"
	"net/http"
	"strings"

	"github.com/unistack-org/micro/v3/metadata"
)

func FromRequest(r *http.Request) context.Context {
	ctx := r.Context()
	md, ok := metadata.FromContext(ctx)
	if !ok {
		// create needed map with specific len
		md = make(metadata.Metadata, len(r.Header)+2)
	}
	for key, val := range r.Header {
		md.Set(key, strings.Join(val, ","))
	}
	// pass http host
	md["Host"] = r.Host
	// pass http method
	md["Method"] = r.Method
	return metadata.NewContext(ctx, md)
}
