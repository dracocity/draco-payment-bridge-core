package plugin

import (
	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
)

type Plugin interface {
	Load(config.PGConfig) error
	Unload() error
	Name() string
	PaymentGateway() pg.PaymentGateway
}
