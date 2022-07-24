package kafka

func WithTracing(enabled bool) Option {
	return func(p *producer) {
		p.enableTracing = enabled
	}
}
