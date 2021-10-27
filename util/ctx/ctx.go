package ctx // import "go.unistack.org/micro/v3/util/ctx"

import (
	"context"
	"net/http"
	"strings"

	"go.unistack.org/micro/v3/metadata"
)

// FromRequest creates context with metadata from http.Request
func FromRequest(r *http.Request) context.Context {
	ctx := r.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(len(r.Header) + 2)
	}
	for key, val := range r.Header {
		md.Set(key, strings.Join(val, ","))
	}
	// pass http host
	md["Host"] = r.Host
	// pass http method
	md["Method"] = r.Method
	return metadata.NewIncomingContext(ctx, md)
}
