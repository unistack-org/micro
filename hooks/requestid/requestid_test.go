package requestid

import (
	"context"
	"testing"

	"go.unistack.org/micro/v4/metadata"
)

func TestDefaultMetadataFunc(t *testing.T) {
	ctx := context.TODO()

	nctx, err := DefaultMetadataFunc(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

	imd, ok := metadata.FromIncomingContext(nctx)
	if !ok {
		t.Fatalf("md missing in incoming context")
	}
	omd, ok := metadata.FromOutgoingContext(nctx)
	if !ok {
		t.Fatalf("md missing in outgoing context")
	}

	_, iok := imd.Get(DefaultMetadataKey)
	_, ook := omd.Get(DefaultMetadataKey)

	if !iok || !ook {
		t.Fatalf("missing metadata key value")
	}
}
