package logging

// WithExclusions - Excludes multiple paths from generating logs
func WithExclusions(names ...string) Option {
	return func(r *Logger) {
		for _, name := range names {
			r.notLogged[name] = struct{}{}
		}
	}
}
