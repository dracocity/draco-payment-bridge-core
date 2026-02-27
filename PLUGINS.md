# Plugin System

This project loads payment provider integrations as Go plugins (`.so`). Each plugin must:

- Be built with `-buildmode=plugin`.
- Use `package main`.
- Export a `New() plugin.Plugin` function.
- Implement `plugin.Plugin` and return a `bridge.Bridge`.

## Interface Overview

```go
// internal/plugin/plugin.go
// type Plugin interface {
//     New() error
//     Load() error
//     Unload() error
//     GetName() string
//     Bridge() bridge.Bridge
// }

// internal/bridge/interface.go
// type Bridge interface {
//     Name() string
//     CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error)
//     GetStatus(ctx context.Context, paymentID string) (*models.PaymentStatusResponse, error)
//     Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error)
//     HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error)
// }
```

## Example Plugin Skeleton

```go
package main

import (
    "os"

    "github.com/dracocity/draco-payment-bridge-core/internal/bridge"
    "github.com/dracocity/draco-payment-bridge-core/internal/models"
    "github.com/dracocity/draco-payment-bridge-core/internal/plugin"
)

type samplePlugin struct {
    apiKey string
    bridge *sampleBridge
}

type sampleBridge struct{}

func New() plugin.Plugin { return &samplePlugin{} }

func (p *samplePlugin) Load() error {
    p.apiKey = os.Getenv("SAMPLE_API_KEY")
    p.bridge = &sampleBridge{}
    return nil
}

func (p *samplePlugin) Unload() error { return nil }

func (p *samplePlugin) GetName() string { return "sample" }

func (p *samplePlugin) Bridge() bridge.Bridge { return p.bridge }

func (b *sampleBridge) Name() string { return "sample" }

func (b *sampleBridge) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
    // TODO: implement provider API call
    return nil, nil
}

// Implement GetStatus, Refund, HandleWebhook...
```

## Build Plugins

```bash
./scripts/build-plugins.sh
```

Plugins are loaded from `PLUGIN_DIR` (default `./plugins/dist`).
