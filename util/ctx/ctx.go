package ctx

import "context"

// Mixin mixes a typically shorter-lived (service, etc.) context passed in
// “shortctx” into a long-living context “longctx”, returning a derived “mixed”
// context. The “mixed” context is derived from the long-living context and
// additionally gets the deadline (if any) and cancellation of the mixed-in
// “shortctx”.
//
// Please note that it is essential to always call the additionally returned
// cancel function in order to not leak go routines. Calling this cancel
// function won't cancel the long-living and short-lived contexts, just clean
// up. This follows the established context pattern of [context.WithCancel],
// [context.WithDeadline] and [context.WithTimeout].
func Mixin(longctx, shortctx context.Context) (mixedctx context.Context, mixedcancel context.CancelFunc) {
	if longctx == nil {
		panic("wye.Mixin: cannot mix into nil context")
	}
	if shortctx == nil {
		panic("wye.Mixin: cannot mix-in nil context")
	}
	mixedctx = longctx
	shortDone := shortctx.Done()
	if shortDone == nil {
		// In case the shorter-living context isn't cancellable at all, then we
		// cannot cancel it and the cancel function returned must be a "no
		// operation". There's no need to mix in something, so we can pass on
		// the long-lived context.
		mixedcancel = func() {}
		return
	}
	// Nota bene: cancelled contexts "trickle down", so if a context higher up
	// the hierarchy was cancelled, this will automatically propagate down to
	// all child contexts, and so on.
	//
	// In case the shorter-lived context has a deadline, we need to carry it
	// over into the final mixed context.
	if deadline, ok := shortctx.Deadline(); ok {
		mixedctx, mixedcancel = context.WithDeadline(mixedctx, deadline)
	}
	// As the shorter-living context can be cancelled, we will need to supervise
	// it so we notice when it gets cancelled and then cancel the mixed context.
	mixedctx, mixedcancel = context.WithCancel(mixedctx)
	go func(ctx context.Context, cancel context.CancelFunc) {
		select {
		case <-shortDone:
			// In case the shorter-living context was cancelled (but it did
			// not pass a deadline) then we need to propagate this mixed
			// context. Please note this correctly won't cancel the original
			// longer context, because that's a long-living context we
			// shouldn't interfere with.
			if shortctx.Err() == context.Canceled {
				cancel()
			}
		case <-ctx.Done():
			// The final mixed context was either cancelled itself or its parent
			// deadline context met its fate; here, do not touch short-lived
			// context.
		}
	}(mixedctx, mixedcancel)
	return
}
