package register

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

var testData = map[string][]*Service{
	"foo": {
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*Node{
				{
					ID:      "foo-1.0.0-123",
					Address: "localhost:9999",
				},
				{
					ID:      "foo-1.0.0-321",
					Address: "localhost:9999",
				},
			},
		},
		{
			Name:    "foo",
			Version: "1.0.1",
			Nodes: []*Node{
				{
					ID:      "foo-1.0.1-321",
					Address: "localhost:6666",
				},
			},
		},
		{
			Name:    "foo",
			Version: "1.0.3",
			Nodes: []*Node{
				{
					ID:      "foo-1.0.3-345",
					Address: "localhost:8888",
				},
			},
		},
	},
	"bar": {
		{
			Name:    "bar",
			Version: "default",
			Nodes: []*Node{
				{
					ID:      "bar-1.0.0-123",
					Address: "localhost:9999",
				},
				{
					ID:      "bar-1.0.0-321",
					Address: "localhost:9999",
				},
			},
		},
		{
			Name:    "bar",
			Version: "latest",
			Nodes: []*Node{
				{
					ID:      "bar-1.0.1-321",
					Address: "localhost:6666",
				},
			},
		},
	},
}

//nolint:gocyclo
func TestMemoryRegistry(t *testing.T) {
	ctx := context.TODO()
	m := NewRegister()

	fn := func(k string, v []*Service) {
		services, err := m.LookupService(ctx, k)
		if err != nil {
			t.Errorf("Unexpected error getting service %s: %v", k, err)
		}

		if len(services) != len(v) {
			t.Errorf("Expected %d services for %s, got %d", len(v), k, len(services))
		}

		for _, service := range v {
			var seen bool
			for _, s := range services {
				if s.Version == service.Version {
					seen = true
					break
				}
			}
			if !seen {
				t.Errorf("expected to find version %s", service.Version)
			}
		}
	}

	// register data
	for _, v := range testData {
		serviceCount := 0
		for _, service := range v {
			if err := m.Register(ctx, service); err != nil {
				t.Errorf("Unexpected register error: %v", err)
			}
			serviceCount++
			// after the service has been registered we should be able to query it
			services, err := m.LookupService(ctx, service.Name)
			if err != nil {
				t.Errorf("Unexpected error getting service %s: %v", service.Name, err)
			}
			if len(services) != serviceCount {
				t.Errorf("Expected %d services for %s, got %d", serviceCount, service.Name, len(services))
			}
		}
	}

	// using test data
	for k, v := range testData {
		fn(k, v)
	}

	services, err := m.ListServices(ctx)
	if err != nil {
		t.Errorf("Unexpected error when listing services: %v", err)
	}

	totalServiceCount := 0
	for _, testSvc := range testData {
		for range testSvc {
			totalServiceCount++
		}
	}

	if len(services) != totalServiceCount {
		t.Errorf("Expected total service count: %d, got: %d", totalServiceCount, len(services))
	}

	// deregister
	for _, v := range testData {
		for _, service := range v {
			if err := m.Deregister(ctx, service); err != nil {
				t.Errorf("Unexpected deregister error: %v", err)
			}
		}
	}

	// after all the service nodes have been deregistered we should not get any results
	for _, v := range testData {
		for _, service := range v {
			services, err := m.LookupService(ctx, service.Name)
			if err != ErrNotFound {
				t.Errorf("Expected error: %v, got: %v", ErrNotFound, err)
			}
			if len(services) != 0 {
				t.Errorf("Expected %d services for %s, got %d", 0, service.Name, len(services))
			}
		}
	}
}

func TestMemoryRegistryTTL(t *testing.T) {
	m := NewRegister()
	ctx := context.TODO()

	for _, v := range testData {
		for _, service := range v {
			if err := m.Register(ctx, service, RegisterTTL(time.Millisecond)); err != nil {
				t.Fatal(err)
			}
		}
	}

	time.Sleep(ttlPruneTime * 2)

	for name := range testData {
		svcs, err := m.LookupService(ctx, name)
		if err != nil {
			t.Fatal(err)
		}

		for _, svc := range svcs {
			if len(svc.Nodes) > 0 {
				t.Fatalf("Service %q still has nodes registered", name)
			}
		}
	}
}

func TestMemoryRegistryTTLConcurrent(t *testing.T) {
	concurrency := 1000
	waitTime := ttlPruneTime * 2
	m := NewRegister()
	ctx := context.TODO()
	for _, v := range testData {
		for _, service := range v {
			if err := m.Register(ctx, service, RegisterTTL(waitTime/2)); err != nil {
				t.Fatal(err)
			}
		}
	}

	//if len(os.Getenv("IN_TRAVIS_CI")) == 0 {
	//	t.Logf("test will wait %v, then check TTL timeouts", waitTime)
	//}

	errChan := make(chan error, concurrency)
	syncChan := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		go func() {
			<-syncChan
			for name := range testData {
				svcs, err := m.LookupService(ctx, name)
				if err != nil {
					errChan <- err
					return
				}

				for _, svc := range svcs {
					if len(svc.Nodes) > 0 {
						errChan <- fmt.Errorf("Service %q still has nodes registered", name)
						return
					}
				}
			}

			errChan <- nil
		}()
	}

	time.Sleep(waitTime)
	close(syncChan)

	for i := 0; i < concurrency; i++ {
		if err := <-errChan; err != nil {
			t.Fatal(err)
		}
	}
}

func TestMemoryWildcard(t *testing.T) {
	m := NewRegister()
	ctx := context.TODO()

	testSrv := &Service{Name: "foo", Version: "1.0.0"}

	if err := m.Register(ctx, testSrv, RegisterDomain("one")); err != nil {
		t.Fatalf("Register err: %v", err)
	}
	if err := m.Register(ctx, testSrv, RegisterDomain("two")); err != nil {
		t.Fatalf("Register err: %v", err)
	}

	if recs, err := m.ListServices(ctx, ListDomain("one")); err != nil {
		t.Errorf("List err: %v", err)
	} else if len(recs) != 1 {
		t.Errorf("Expected 1 record, got %v", len(recs))
	}

	if recs, err := m.ListServices(ctx, ListDomain("*")); err != nil {
		t.Errorf("List err: %v", err)
	} else if len(recs) != 2 {
		t.Errorf("Expected 2 records, got %v", len(recs))
	}

	if recs, err := m.LookupService(ctx, testSrv.Name, LookupDomain("one")); err != nil {
		t.Errorf("Lookup err: %v", err)
	} else if len(recs) != 1 {
		t.Errorf("Expected 1 record, got %v", len(recs))
	}

	if recs, err := m.LookupService(ctx, testSrv.Name, LookupDomain("*")); err != nil {
		t.Errorf("Lookup err: %v", err)
	} else if len(recs) != 2 {
		t.Errorf("Expected 2 records, got %v", len(recs))
	}
}

func TestWatcher(t *testing.T) {
	testSrv := &Service{Name: "foo", Version: "1.0.0"}

	ctx := context.TODO()
	m := NewRegister()
	m.Init()
	m.Connect(ctx)
	wc, err := m.Watch(ctx)
	if err != nil {
		t.Fatalf("cant watch: %v", err)
	}
	defer wc.Stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			_, err := wc.Next()
			if err != nil {
				t.Fatal("unexpected err", err)
			}
			// t.Logf("changes %#+v", ch.Service)
			wc.Stop()
			wg.Done()
			return
		}
	}()

	if err := m.Register(ctx, testSrv); err != nil {
		t.Fatalf("Register err: %v", err)
	}

	wg.Wait()
	if _, err := wc.Next(); err == nil {
		t.Fatal("expected error on Next()")
	}
}
