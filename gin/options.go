package ginPlugin

type Option func(*options)

type options struct {
	excludePaths  []string
	fromBodyPaths []string
}

func WithExcludePaths(e []string) Option {
	return func(o *options) {
		o.excludePaths = e
	}
}

func WithFromBodyPaths(f []string) Option {
	return func(o *options) {
		o.fromBodyPaths = f
	}
}
