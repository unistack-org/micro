package micro

import (
	"context"
	"reflect"
	"testing"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/config"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/router"
	"go.unistack.org/micro/v3/server"
	"go.unistack.org/micro/v3/store"
	"go.unistack.org/micro/v3/tracer"
)

func TestClient(t *testing.T) {
	c1 := client.NewClient(client.Name("test1"))
	c2 := client.NewClient(client.Name("test2"))

	svc := NewService(Client(c1, c2))

	if err := svc.Init(); err != nil {
		t.Fatal(err)
	}

	x1 := svc.Client("test2")
	if x1.Name() != "test2" {
		t.Fatalf("invalid client %#+v", svc.Options().Clients)
	}
}

type testItem struct {
	name string
}

func (ti *testItem) Name() string {
	return ti.name
}

func Test_getNameIndex(t *testing.T) {
	items := []*testItem{{name: "test1"}, {name: "test2"}}
	idx := getNameIndex("test2", items)
	if items[idx].Name() != "test2" {
		t.Fatal("getNameIndex wrong")
	}
}

func TestRegisterHandler(t *testing.T) {
	type args struct {
		s    server.Server
		h    interface{}
		opts []server.HandlerOption
	}
	h := struct{}{}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "RegisterHandler",
			args: args{
				s:    server.DefaultServer,
				h:    h,
				opts: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterHandler(tt.args.s, tt.args.h, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("RegisterHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterSubscriber(t *testing.T) {
	type args struct {
		topic string
		s     server.Server
		h     interface{}
		opts  []server.SubscriberOption
	}
	h := func(_ context.Context, _ interface{}) error {
		return nil
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "RegisterSubscriber",
			args: args{
				topic: "test",
				s:     server.DefaultServer,
				h:     h,
				opts:  nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterSubscriber(tt.args.topic, tt.args.s, tt.args.h, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("RegisterSubscriber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "NewService",
			args: args{
				opts: []Option{Name("test")},
			},
			want: NewService(Name("test")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got.Options().Name, tt.want.Options().Name)
			}
		})
	}
}

func Test_service_Name(t *testing.T) {
	type fields struct {
		opts Options
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test_service_Name",
			fields: fields{
				opts: Options{Name: "test"},
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Name(); got != tt.want {
				t.Errorf("service.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Init(t *testing.T) {
	type fields struct {
		opts Options
	}
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "service.Init()",
			fields: fields{
				opts: Options{},
			},
			args: args{
				opts: []Option{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if err := s.Init(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("service.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_Options(t *testing.T) {
	opts := Options{Name: "test"}
	type fields struct {
		opts Options
	}
	tests := []struct {
		name   string
		fields fields
		want   Options
	}{
		{
			name: "service.Options",
			fields: fields{
				opts: opts,
			},
			want: opts,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Broker(t *testing.T) {
	b := broker.NewBroker()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   broker.Broker
	}{
		{
			name: "service.Broker",
			fields: fields{
				opts: Options{Brokers: []broker.Broker{b}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: b,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Broker(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Broker() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func TestServiceBroker(t *testing.T) {
	b := broker.NewBroker(broker.Name("test"))

	srv := server.NewServer()

	svc := NewService(Server(srv),Broker(b))

	if err := svc.Init(); err != nil {
		t.Fatalf("failed to init service")
	}

	if brk := svc.Server().Options().Broker; brk.Name() != "test" {
		t.Fatalf("server broker not set: %v", svc.Server().Options().Broker)
	}

}
*/

func Test_service_Tracer(t *testing.T) {
	tr := tracer.NewTracer()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   tracer.Tracer
	}{
		{
			name: "service.Tracer",
			fields: fields{
				opts: Options{Tracers: []tracer.Tracer{tr}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: tr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Tracer(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Tracer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Config(t *testing.T) {
	c := config.NewConfig()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   config.Config
	}{
		{
			name: "service.Config",
			fields: fields{
				opts: Options{Configs: []config.Config{c}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: c,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Config(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Config() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Client(t *testing.T) {
	c := client.NewClient()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   client.Client
	}{
		{
			name: "service.Client",
			fields: fields{
				opts: Options{Clients: []client.Client{c}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: c,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Client(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Client() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Server(t *testing.T) {
	s := server.NewServer()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   server.Server
	}{
		{
			name: "service.Server",
			fields: fields{
				opts: Options{Servers: []server.Server{s}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Server(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Server() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Store(t *testing.T) {
	s := store.NewStore()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   store.Store
	}{
		{
			name: "service.Store",
			fields: fields{
				opts: Options{Stores: []store.Store{s}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Store(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Store() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Register(t *testing.T) {
	r := register.NewRegister()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   register.Register
	}{
		{
			name: "service.Register",
			fields: fields{
				opts: Options{Registers: []register.Register{r}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: r,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Register(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Logger(t *testing.T) {
	l := logger.NewLogger()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   logger.Logger
	}{
		{
			name: "service.Logger",
			fields: fields{
				opts: Options{Loggers: []logger.Logger{l}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: l,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Logger(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Logger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Router(t *testing.T) {
	r := router.NewRouter()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   router.Router
	}{
		{
			name: "service.Router",
			fields: fields{
				opts: Options{Routers: []router.Router{r}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: r,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Router(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Router() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_Meter(t *testing.T) {
	m := meter.NewMeter()
	type fields struct {
		opts Options
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   meter.Meter
	}{
		{
			name: "service.Meter",
			fields: fields{
				opts: Options{Meters: []meter.Meter{m}},
			},
			args: args{
				names: []string{"noop"},
			},
			want: m,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.Meter(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Meter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_String(t *testing.T) {
	type fields struct {
		opts Options
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "service.String",
			fields: fields{
				opts: Options{Name: "noop"},
			},
			want: "noop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				opts: tt.fields.opts,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("service.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func Test_service_Start(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		opts    Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				RWMutex: tt.fields.RWMutex,
				opts:    tt.fields.opts,
			}
			if err := s.Start(); (err != nil) != tt.wantErr {
				t.Errorf("service.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_Stop(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		opts    Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				RWMutex: tt.fields.RWMutex,
				opts:    tt.fields.opts,
			}
			if err := s.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("service.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_Run(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		opts    Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				RWMutex: tt.fields.RWMutex,
				opts:    tt.fields.opts,
			}
			if err := s.Run(); (err != nil) != tt.wantErr {
				t.Errorf("service.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getNameIndex(t *testing.T) {
	type args struct {
		n      string
		ifaces interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNameIndex(tt.args.n, tt.args.ifaces); got != tt.want {
				t.Errorf("getNameIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
