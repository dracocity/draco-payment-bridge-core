package plugin

import "github.com/dracocity/draco-payment-bridge-core/internal/bridge"

type Plugin interface {
	New() error
	Load() error
	Unload() error
	GetName() string
	Bridge() bridge.Bridge
}
