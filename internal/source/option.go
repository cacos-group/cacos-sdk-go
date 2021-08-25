package source

// Option is etcd config option.
type Option func(o *options)

type options struct {
	path   string
	prefix bool
}

// WithPrefix is config prefix
func WithPrefix(prefix bool) Option {
	return Option(func(o *options) {
		o.prefix = prefix
	})
}
