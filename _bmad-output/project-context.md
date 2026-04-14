# Project Context - draco-payment-bridge-core

## 1. 프로젝트 개요
- 목적: 다수의 암호화폐 결제 제공자 API를 단일 HTTP 인터페이스로 브릿지하는 Go 서버.
- 핵심 방식: Go plugin(`.so`)을 동적 로딩해 provider 구현을 런타임에 주입.
- 현재 확인된 provider: `nowpayments`.
- API base path: `/api/v1`.

## 2. 기술 스택
- Language: Go `1.25` (`toolchain go1.25.6`)
- HTTP: `gin-gonic/gin`
- CLI Entrypoint: `urfave/cli/v2`
- Config: `go-toml/v2` + `envconfig`
- Logging: `zap` + `lumberjack`
- Numeric: `shopspring/decimal`
- Plugin: Go 표준 `plugin` 패키지

## 3. 실행/빌드/테스트
- 플러그인 빌드: `./scripts/build-plugins.sh`
- 서버 실행: `go run ./cmd/dpbc -c ./config.toml`
- 기본 테스트: `go test ./...`
- 주요 테스트 파일:
  - `internal/config/load_config_test.go`
  - `internal/handlers/payment_handler_test.go`
  - `internal/httpx/client_test.go`
  - `internal/services/payment_service_test.go`
  - `plugins/nowpayments/nowpayments_test.go`

## 4. 아키텍처 요약
- 진입점: `cmd/dpbc/main.go`
- 흐름:
  1. `config.Load`로 설정 로드/정규화
  2. logger 초기화
  3. `plugin_loader`가 `plugin_dir`에서 `.so` 스캔 및 로드
  4. plugin `Load(cfg)` 후 `PaymentGateway()`를 registry에 등록
  5. `handlers -> services -> registry -> provider` 호출 체인으로 요청 처리
  6. SIGINT/SIGTERM에서 graceful shutdown + unix socket 정리

## 5. 핵심 도메인 계약
### 5.1 PaymentGateway 인터페이스
`internal/pg/interface.go`
- `CreatePaymentLink`
- `CreatePayment`
- `GetPayment`
- `CreateRefund`
- `HandleWebhook`

### 5.2 Handler/Service 에러 매핑
- 서비스 sentinel error: `services.ErrUnknownProvider`
- 핸들러는 이를 HTTP 코드로 변환 (`400/404` 케이스 존재)
- provider 실패는 `502 Bad Gateway`로 응답

## 6. 설정 컨텍스트
- 우선순위: TOML > 환경변수 > 기본값 (`internal/config/load_config.go`)
- listen 포맷: `network@address[,network@address...]`
- 지원 네트워크: `tcp`, `unix`
- provider env 패턴: `DPBC_PROVIDERS_<PROVIDER>_<KEY>`

## 7. 코드 스타일/규칙 컨텍스트
- 스타일 가이드: `docs/code-style.md`
- 의존 방향: `handlers -> services -> pg/plugin`
- 에러 래핑: `%w` 사용
- 로깅: 구조화 key-value

## 8. 현재 알려진 구현 메모
- webhook payload 원문을 debug 로그로 남기는 코드가 있음 (`internal/handlers/payment_handler.go`).
- HTTP 클라이언트가 response payload 전체를 info 로그로 기록 (`internal/httpx/client.go`).
- 환불 입력 검증 메시지 문구에 오타성 표현 존재 (`payment_id are required`).

## 9. BMAD 워크플로우 라우팅 제안
이 저장소는 이미 구현된 브라운필드 서버이므로, 다음 순서가 적합함.
1. 완료: `bmad-generate-project-context (GPC)`
2. 기능 단위 작업 시작 시: `bmad-quick-dev (QQ)` 또는 `bmad-create-story (CS)`
3. 구현 후 품질 점검: `bmad-code-review (CR)`

## 10. 다음 액션 제안
- 즉시 작업 시작이 목적이면 `QQ`로 첫 개선 항목(로깅 민감정보 제거)을 스펙-구현-검토 사이클로 진행.
