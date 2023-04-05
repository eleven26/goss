package goss

type Goss struct {
	config *Config
	Store
}

func New(opts ...Option) (*Goss, error) {
	goss := &Goss{}

	for _, option := range opts {
		option(goss)
	}

	err := goss.config.validate()
	if err != nil {
		return nil, err
	}

	goss.Store, err = newStorage(goss.config)
	if err != nil {
		return nil, err
	}

	return goss, nil
}
