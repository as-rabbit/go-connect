package option

type Config struct {
	Address []string
}

type Option func(c Config)

func WithAddress(address []string) Option {
	return func(c Config) {
		c.Address = address
	}
}
