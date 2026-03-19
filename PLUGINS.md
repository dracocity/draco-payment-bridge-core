# Plugin System

This project loads payment provider integrations as Go plugins (`.so`). Each plugin must:

- Be built with `-buildmode=plugin`.
- Use `package main`.
- Export a `New() plugin.Plugin` function.
- Implement `internal/plugin.Plugin` and return `pg.PaymentGateway`.

## Interface Overview

```go
// internal/plugin/plugin.go
// type Plugin interface {
//     Load(cfg config.PGConfig) error
//     Unload() error
//     Name() string
//     PaymentGateway() pg.PaymentGateway
// }

// internal/pg/interface.go
// type PaymentGateway interface {
//     Name() string
//     CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error)
//     CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error)
//     GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error)
//     Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error)
//     HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error)
// }
```

## Example Plugin Skeleton

```go
package main

import (
    "context"

    "github.com/dracocity/draco-payment-bridge-core/internal/config"
    "github.com/dracocity/draco-payment-bridge-core/internal/models"
    "github.com/dracocity/draco-payment-bridge-core/internal/pg"
    "github.com/dracocity/draco-payment-bridge-core/internal/plugin"
)

type samplePlugin struct {
    gateway *samplePG
}

type samplePG struct{}

func New() plugin.Plugin { return &samplePlugin{} }

func (p *samplePlugin) Load(cfg config.PGConfig) error {
    // cfg["api_key"], cfg["mode"] ...
    p.gateway = &samplePG{}
    return nil
}

func (p *samplePlugin) Unload() error { return nil }

func (p *samplePlugin) Name() string { return "sample" }

func (p *samplePlugin) PaymentGateway() pg.PaymentGateway { return p.gateway }

func (g *samplePG) Name() string { return "sample" }

func (g *samplePG) CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
    return nil, nil
}

func (g *samplePG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
    return nil, nil
}

func (g *samplePG) GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error) {
    return nil, nil
}

func (g *samplePG) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
    return nil, nil
}

func (g *samplePG) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
    return nil, nil
}
```

## Build Plugins

```bash
./scripts/build-plugins.sh
```

The build script outputs plugins to `./dist/release/plugins`.

## Runtime Plugin Directory

Runtime plugin loading path is configured by `plugin_dir` in `config.toml`.

- Default: `./plugins`
- If you use `./scripts/build-plugins.sh`, set `plugin_dir = "./dist/release/plugins"`
