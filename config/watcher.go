package config

import (
	"reflect"

	"github.com/unistack-org/micro/v3/util/jitter"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

type defaultWatcher struct {
	opts  Options
	wopts WatchOptions
	done  chan struct{}
	vchan chan map[string]interface{}
	echan chan error
}

func (w *defaultWatcher) run() {
	ticker := jitter.NewTicker(w.wopts.MinInterval, w.wopts.MaxInterval)
	defer ticker.Stop()

	src := w.opts.Struct
	if w.wopts.Struct != nil {
		src = w.wopts.Struct
	}

	for {
		select {
		case <-w.done:
			return
		case <-ticker.C:
			dst, err := rutil.Zero(src)
			if err == nil {
				err = fillValues(reflect.ValueOf(dst), w.opts.StructTag)
			}
			if err != nil {
				w.echan <- err
				return
			}
			srcmp, err := rutil.StructFieldsMap(src)
			if err != nil {
				w.echan <- err
				return
			}
			dstmp, err := rutil.StructFieldsMap(dst)
			if err != nil {
				w.echan <- err
				return
			}
			for sk, sv := range srcmp {
				if reflect.DeepEqual(dstmp[sk], sv) {
					delete(dstmp, sk)
				}
			}
			if len(dstmp) > 0 {
				w.vchan <- dstmp
				src = dst
			}
		}
	}
}

func (w *defaultWatcher) Next() (map[string]interface{}, error) {
	select {
	case <-w.done:
		break
	case err := <-w.echan:
		return nil, err
	case v, ok := <-w.vchan:
		if !ok {
			break
		}
		return v, nil
	}
	return nil, ErrWatcherStopped
}

func (w *defaultWatcher) Stop() error {
	close(w.done)
	return nil
}
