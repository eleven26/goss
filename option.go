package goss

type Option func(goss *Goss)

func WithConfig(config *Config) Option {
	return func(goss *Goss) {
		goss.config = config
	}
}
