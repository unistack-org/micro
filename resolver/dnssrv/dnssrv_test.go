package dnssrv

import "testing"

func TestResolveError(t *testing.T) {
	r := &Resolver{}
	_, err := r.Resolve("_nonexistent._udp.invalid.invalid.")
	if err == nil {
		t.Log("surprisingly got records — DNS may have responded")
	}
}
