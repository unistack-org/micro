package build

// Options struct
type Options struct {
	// local path to download source
	Path string
}

// Option func
type Option func(o *Options)

// Path is the Local path for repository
func Path(p string) Option {
	return func(o *Options) {
		o.Path = p
	}
}
