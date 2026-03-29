package app

import (
	"context"
	"golab/internal/platform/config"
)

type Container struct {
	Config config.Config
	Infra  *Infra
	User   UserModule
}

func NewContainer(ctx context.Context) (*Container, error) {
	return NewContainerWithConfig(ctx, config.Load())
}

func NewContainerWithConfig(ctx context.Context, cfg config.Config) (*Container, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	infra, err := bootstrapInfra(ctx, cfg)
	if err != nil {
		return nil, err
	}

	userModule := buildUserModule(infra, cfg)

	return &Container{
		Config: cfg,
		Infra:  infra,
		User:   userModule,
	}, nil
}

func (c *Container) Close() error {
	if c == nil || c.Infra == nil {
		return nil
	}

	if c.Infra.Rabbit != nil {
		_ = c.Infra.Rabbit.Close()
	}

	if c.Infra.DB != nil {
		c.Infra.DB.Close()
	}

	return nil
}
