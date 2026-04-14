# Draco Payment Bridge Core

Go 기반 암호화폐 결제 브리지 서버입니다. 결제 제공자 플러그인(`.so`)을 동적으로 로드해 단일 HTTP API로 결제 링크 생성, 결제 생성/조회, 환불, 웹훅 처리를 제공합니다.

## 주요 기능

- 다중 결제 제공자 통합 (`nowpayments`, `bitpay(예정)`, `coingate(예정)`, `coinpayments(예정)`)
- 플러그인 기반 로딩 (`plugin_dir` 내 `.so` 자동 스캔)
- 통합 결제 API (`/api/v1/...`)
- TOML + 환경변수 기반 설정 (`DPBC_*`)
- TCP / Unix socket 멀티 리스닝 지원

## 요구 사항

- Go `1.25+`
- Linux/macOS (Go plugin 사용)

## 빠른 시작

### 1. 의존성 설치

```bash
go mod download
```

### 2. 플러그인 빌드

```bash
chmod +x scripts/build-plugins.sh
./scripts/build-plugins.sh
```

기본 출력 경로:

- `dist/release/plugins/nowpayments.so`

### 3. 설정 파일 생성

```bash
cp config.toml.sample config.toml
```

플러그인을 위 스크립트 기본 경로로 빌드했다면 `config.toml`의 `plugin_dir`를 다음처럼 변경하세요.

```toml
plugin_dir = "./dist/release/plugins"
```

### 4. 서버 실행

```bash
go run ./cmd/dpbc -c ./config.toml
```

기본 `-c` 값은 `./config.toml` 입니다.  
`-c ""`(빈 문자열)로 실행하면 파일 로딩 없이 환경변수/기본값만 사용합니다.

## 설정

### TOML

- 샘플: `config.toml.sample`
- 로드 우선순위: `config.toml` > 환경변수 > 기본값
- `plugin_dir`가 상대 경로면 실행 시 절대 경로로 변환되며, 디렉터리가 없으면 자동 생성됩니다.

### 핵심 설정 항목

```toml
plugin_dir = "./dist/release/plugins"
listen = "tcp@127.0.0.1:8080,unix@/tmp/dpbc.sock"

[providers.nowpayments]
mode = "sandbox"
api_key = "NOWPAYMENTS-API-KEY"
ipn_secret = "NOWPAYMENTS-IPN-SECRET"
```

### 환경변수 오버라이드

- 리슨 주소: `DPBC_LISTEN="tcp@127.0.0.1:8080,unix@/tmp/dpbc.sock"`
- 플러그인 경로: `DPBC_PLUGIN_DIR`
- 로그: `DPBC_LOG_LEVEL`, `DPBC_LOG_FORMAT`, `DPBC_LOG_STDOUT`, `DPBC_LOG_FILE` 등
- 제공자별: `DPBC_PROVIDERS_<PROVIDER>_<KEY>`
예: `DPBC_PROVIDERS_NOWPAYMENTS_API_KEY=...`

## API

Base Path: `/api/v1`

### 1) Health Check

`GET /api/v1/health`

응답 예시:

```json
{
  "status": "ok",
  "providers": ["nowpayments"]
}
```

### 2) 결제 링크 생성

`POST /api/v1/providers/:provider/payment-links`

요청 예시:

```json
{
  "amount": "100.00",
  "currency": "USD",
  "receive_currency": "BTC",
  "order_id": "ORDER-123",
  "description": "Test payment",
  "customer_email": "buyer@example.com",
  "webhook_url": "https://example.com/webhook",
  "success_url": "https://example.com/success",
  "cancel_url": "https://example.com/cancel"
}
```

### 3) 결제 생성

`POST /api/v1/providers/:provider/payments`

요청 예시:

```json
{
  "invoice_id": "inv-123",
  "receive_currency": "BTC",
  "description": "Test payment",
  "customer_email": "buyer@example.com",
  "provider_payload": {}
}
```

### 4) 결제 조회

`GET /api/v1/providers/:provider/payments/:paymentId`

### 5) 환불

`POST /api/v1/providers/:provider/refunds`

요청 예시:

```json
{
  "payment_id": "payment-id-123",
  "amount": 50,
  "reason": "customer request"
}
```

### 6) 웹훅

`POST /api/v1/providers/:provider/webhooks`

## 프로젝트 구조

```text
draco-payment-bridge-core/
├── cmd/dpbc/                  # 애플리케이션 진입점
├── internal/
│   ├── config/                # TOML/env 설정 로딩
│   ├── handlers/              # HTTP 핸들러
│   ├── models/                # API 모델
│   ├── pg/                    # PaymentGateway 인터페이스/레지스트리
│   ├── plugin/                # Go plugin 로더
│   └── services/              # 비즈니스 서비스
├── plugins/                   # 각 제공자 플러그인 소스
├── scripts/build-plugins.sh   # 플러그인 빌드 스크립트
├── config.toml.sample
└── README.md
```

## 참고

- 플러그인 구현 상세는 `PLUGINS.md`를 참고하세요.
- 코드 스타일 가이드는 `docs/code-style.md`를 참고하세요.
