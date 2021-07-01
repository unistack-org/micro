package fn

type Initer interface {
	Init(opts ...interface{}) error
}

func Init(ifaces ...Initer) error {
	var err error
	for _, iface := range ifaces {
		if err = iface.Init(); err != nil {
			return err
		}
	}
	return nil
}
