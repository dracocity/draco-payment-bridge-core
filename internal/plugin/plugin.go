package plugin

import "github.com/dracocity/draco-payment-bridge-core/internal/pg"

type Plugin interface {
	New() error
	Load() error
	Unload() error
	GetName() string
	PaymentGateway() pg.PaymentGateway
}
