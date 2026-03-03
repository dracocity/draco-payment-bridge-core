# Payment Bridge Server

다중 결제 제공자(NOWPayments, BitPay, CoinGate, CoinPayments)를 통합하는 Go 서버입니다.

## 기능

- **통합 API**: 여러 결제 제공자를 하나의 API로 통합
- **결제 생성**: 각 제공자별 결제 생성
- **결제 상태 조회**: 결제 상태 실시간 조회
- **환불 처리**: 전체 환불 및 부분 환불 지원
- **웹훅 처리**: 각 제공자의 웹훅 수신 및 처리
- **플러그인 시스템**: 동적 플러그인 로딩으로 필요한 제공자만 포함

## 지원 제공자

1. **NOWPayments** - 다양한 암호화폐 결제
3. **BitPay** - 비트코인 결제 처리
4. **CoinGate** - 암호화폐 결제 게이트웨이
5. **CoinPayments** - 글로벌 암호화폐 결제

## 설치 및 실행

### 1. 의존성 설치

```bash
go mod download
```

### 2. 플러그인 빌드

플러그인 시스템을 사용하므로 먼저 플러그인을 빌드해야 합니다:

**Linux/macOS:**
```bash
chmod +x scripts/build-plugins.sh
./scripts/build-plugins.sh
```

**Windows:**
```cmd
scripts\build-plugins.bat
```

또는 개별 플러그인 빌드:
```bash
go build -buildmode=plugin -o plugins/dist/nowpayments.so plugins/src/nowpayments
go build -buildmode=plugin -o plugins/dist/bitpay.so plugins/src/bitpay
go build -buildmode=plugin -o plugins/dist/coingate.so plugins/src/coingate
go build -buildmode=plugin -o plugins/dist/coinpayments.so plugins/src/coinpayments
```

### 3. 환경 변수 설정

`.env` 파일을 생성하고 필요한 API 키를 설정하세요:

```bash
# .env 파일 예시
PORT=8080
PLUGIN_DIR=./plugins/dist

# NOWPayments
NOWPAYMENTS_API_KEY=your_api_key

# BitPay
BITPAY_API_TOKEN=your_api_token

# CoinGate
COINGATE_API_TOKEN=your_api_token

# CoinPayments
COINPAYMENTS_PUBLIC_KEY=your_public_key
COINPAYMENTS_PRIVATE_KEY=your_private_key
```

### 4. 서버 실행

```bash
go run ./cmd/dpbc
```

또는 빌드 후 실행:

```bash
go build -o payment-bridge ./cmd/dpbc
./payment-bridge
```

서버는 기본적으로 `8080` 포트에서 실행됩니다.

## API 엔드포인트

### 1. 헬스 체크

```
GET /health
```

응답:
```json
{
  "status": "ok",
  "providers": ["nowpayments", "bitpay", "coingate", "coinpayments"]
}
```

### 2. 결제 생성

```
POST /api/payment
Content-Type: application/json
```

요청 본문:
```json
{
  "provider": "nowpayments",
  "amount": 100.00,
  "currency": "USD",
  "crypto": "BTC",
  "description": "Test payment",
  "order_id": "ORDER-123",
  "return_url": "https://example.com/success",
  "cancel_url": "https://example.com/cancel",
  "notify_url": "https://example.com/webhook"
}
```

응답:
```json
{
  "success": true,
  "payment_id": "payment-id-123",
  "payment_url": "https://nowpayments.io/payment?id=...",
  "qr_code": "...",
  "amount": "100.00",
  "currency": "USD",
  "status": "new",
  "expires_at": "2024-01-01T12:00:00Z"
}
```

### 3. 결제 상태 조회

```
GET /api/payment/status?provider=nowpayments&payment_id=payment-id-123
```

응답:
```json
{
  "payment_id": "payment-id-123",
  "status": "confirmed",
  "amount": "100.00",
  "currency": "USD",
  "transaction": "tx-hash-123",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### 4. 환불 처리

```
POST /api/payment/refund
Content-Type: application/json
```

요청 본문 (전체 환불):
```json
{
  "provider": "nowpayments",
  "payment_id": "payment-id-123",
  "reason": "Customer request"
}
```

요청 본문 (부분 환불):
```json
{
  "provider": "nowpayments",
  "payment_id": "payment-id-123",
  "amount": 50.00,
  "reason": "Partial refund"
}
```

응답:
```json
{
  "success": true,
  "refund_id": "refund-id-123",
  "payment_id": "payment-id-123",
  "amount": "50.00",
  "currency": "USD",
  "status": "pending",
  "refund_type": "partial"
}
```

### 5. 웹훅 수신

```
POST /api/webhook?provider=nowpayments
Content-Type: application/json
X-Signature: <signature>
```

## 프로젝트 구조

```
draco-payment-bridge-core/
├── internal/         # 내부 패키지 (외부에서 직접 import 불가)
│   ├── bridge/      # 각 결제 제공자 브리지 구현
│   │   └── interface.go
│   ├── plugin/      # 플러그인 로더
│   │   └── plugin_loader.go
│   ├── config/      # 설정 관리
│   │   ├── config.go
│   │   └── load_config.go
│   ├── handlers/    # HTTP 핸들러
│   │   └── handlers.go
│   └── models/      # 데이터 모델
│       └── payment.go
├── plugins/         # 플러그인 관련
│   ├── src/         # 플러그인 소스 코드
│   │   ├── nowpayments/
│   │   │   └── nowpayments.go
│   │   ├── bitpay/
│   │   │   └── bitpay.go
│   │   ├── coingate/
│   │   │   └── coingate.go
│   │   └── coinpayments/
│   │       └── coinpayments.go
│   └── dist/        # 빌드된 플러그인 파일들 (.so/.dll)
├── scripts/         # 빌드 스크립트
│   ├── build-plugins.sh  # 플러그인 빌드 스크립트 (Linux/macOS)
│   └── build-plugins.bat # 플러그인 빌드 스크립트 (Windows)
├── bin/            # 빌드 산출물
├── cmd/            # 메인 진입점
├── PLUGINS.md      # 플러그인 시스템 상세 문서
├── go.mod
└── README.md
```

## 플러그인 시스템

Payment Bridge는 플러그인 기반 아키텍처를 사용합니다. 자세한 내용은 [PLUGINS.md](PLUGINS.md)를 참조하세요.

### 플러그인 추가/제거

**플러그인 추가:**
```bash
# 플러그인 파일을 plugins/dist 디렉토리에 복사
cp new-plugin.so plugins/dist/
```

**플러그인 제거:**
```bash
# 플러그인 파일 삭제
rm plugins/dist/nowpayments.so
```

애플리케이션 재시작 시 변경사항이 적용됩니다.

### 플러그인 디렉토리 변경

기본 플러그인 디렉토리는 `./plugins/dist`입니다. 다른 디렉토리를 사용하려면:

```bash
export PLUGIN_DIR="/path/to/plugins"
```

또는 `.env` 파일에 설정:
```env
PLUGIN_DIR=/path/to/plugins
```

## 주의사항

1. **서명 검증**: 프로덕션 환경에서는 각 제공자의 웹훅 서명을 올바르게 검증해야 합니다. 현재 구현은 기본 구조만 제공합니다.

2. **에러 처리**: 프로덕션 환경에서는 더 상세한 에러 처리와 로깅이 필요합니다.

3. **보안**: API 키는 절대 코드에 하드코딩하지 말고 환경 변수로 관리하세요.

4. **플러그인 호환성**: 플러그인은 빌드된 플랫폼과 동일한 플랫폼에서만 작동합니다. Linux에서 빌드한 플러그인은 Linux에서만, macOS에서 빌드한 플러그인은 macOS에서만 작동합니다.

5. **Windows 지원**: Go 플러그인은 Windows에서 지원되지 않습니다.

6. **Go 버전**: 플러그인과 메인 애플리케이션은 동일한 Go 버전으로 빌드되어야 합니다.

## 라이선스

MIT
