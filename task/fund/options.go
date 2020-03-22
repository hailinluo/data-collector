package fund

type options struct {
	spec     string
}

type Option func(opts *options)

var defaultOptions = options{
	//spec: "0 * * * * MON-FRI",
	spec:   "0 * * * * *",
}

func WithSpec(spec string) Option {
	return func(opts *options) {
		opts.spec = spec
	}
}

