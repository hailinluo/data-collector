package fundcompany

type options struct {
	spec     string
	homePage string
	resUrl   string
}

type Option func(opts *options)

var defaultOptions = options{
	//spec: "0 * * * * MON-FRI",
	spec:   "0 * * * * *",
	homePage: "",
	resUrl: "",
}

func WithSpec(spec string) Option {
	return func(opts *options) {
		opts.spec = spec
	}
}

func WithHomePage(url string) Option {
	return func(opts *options) {
		opts.homePage = url
	}
}

func WithResUrl(url string) Option {
	return func(opts *options) {
		opts.resUrl = url
	}
}
