package transport

type NoopTransport struct {
	opts Options
}

func (t *NoopTransport) Init(opts ...Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *NoopTransport) Options() Options {
	return t.opts
}

func (t *NoopTransport) Dial(addr string, opts ...DialOption) (Client, error) {
	options := NewDialOptions(opts...)
	return &noopClient{opts: options}, nil
}

func (t *NoopTransport) Listen(addr string, opts ...ListenOption) (Listener, error) {
	options := NewListenOptions(opts...)
	return &noopListener{opts: options}, nil
}

func (t *NoopTransport) String() string {
	return "noop"
}

type noopClient struct {
	opts DialOptions
}

func (c *noopClient) Close() error {
	return nil
}

func (c *noopClient) Local() string {
	return ""
}

func (c *noopClient) Remote() string {
	return ""
}

func (c *noopClient) Recv(*Message) error {
	return nil
}

func (c *noopClient) Send(*Message) error {
	return nil
}

type noopListener struct {
	opts ListenOptions
}

func (l *noopListener) Addr() string {
	return ""
}

func (l *noopListener) Accept(fn func(Socket)) error {
	return nil
}

func (l *noopListener) Close() error {
	return nil
}
