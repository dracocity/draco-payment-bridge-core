# DPBC Code Style Guide

이 문서는 `draco-payment-bridge-core` 저장소에서 일관된 코드 품질을 유지하기 위한 기준입니다.

## 1. 기본 원칙

- 명확성 우선: 짧더라도 의미가 분명한 이름을 사용합니다.
- 작은 단위 유지: 함수는 한 가지 책임만 갖도록 분리합니다.
- 조기 반환 선호: 중첩 `if` 대신 early return으로 흐름을 단순화합니다.
- 기존 패턴 존중: 동일한 레이어(`handlers`, `services`, `plugins`) 내 구현 스타일을 맞춥니다.

## 2. Go 포맷/정렬

- 모든 Go 파일은 `gofmt` 기준을 따릅니다.
- import 그룹은 표준 라이브러리 / 내부 모듈 / 외부 모듈 순서를 유지합니다.
- 수동 정렬보다 자동 포맷 결과를 신뢰합니다.

권장 명령:

```bash
gofmt -w ./...
```

## 3. 패키지 구조 규칙

- `cmd/dpbc`: 애플리케이션 진입점과 런타임 조립 로직.
- `internal/handlers`: HTTP 요청/응답 처리, 상태 코드 매핑.
- `internal/services`: 비즈니스 오케스트레이션, provider 호출 흐름.
- `internal/pg`, `internal/plugin`: 플러그인 인터페이스/로딩 책임.
- `plugins/*`: 각 provider 구현(외부 API 연동 책임).

의존성 방향:

- `handlers -> services -> pg/plugin` 방향으로만 의존합니다.
- `plugins/*`는 `internal/models`, `internal/httpx` 등을 사용할 수 있지만, `handlers`에 의존하지 않습니다.

## 4. 네이밍 규칙

- 공개 심볼: `PascalCase` (`CreatePayment`, `PaymentService`).
- 비공개 심볼: `camelCase` (`loadConfig`, `bindListeners`).
- 인터페이스 이름은 역할 중심 명사로 작성 (`PaymentService`, `PaymentGateway`).
- 축약어는 관례를 따르되 과도한 축약은 피합니다 (`cfg`, `ctx`는 허용).

## 5. 에러 처리

- 에러는 즉시 처리하고 상위로 명확히 전달합니다.
- 에러 래핑은 `%w`를 사용해 원인 체인을 보존합니다.
- 레이어별 메시지 접두어를 붙여 맥락을 남깁니다.

예시:

```go
if err := logger.Init(cfg.Log); err != nil {
    return fmt.Errorf("failed to initialize logger: %w", err)
}
```

## 6. 로깅 규칙

- 구조화 로깅을 사용하고 key-value 형태를 유지합니다.
- 민감정보(API key, secret, 원문 webhook payload 전체)는 로그에 남기지 않습니다.
- 모듈별 로거를 사용합니다 (`logger.WithModule("payment-service")`).

## 7. HTTP 핸들러 규칙

- 입력 검증 실패: `400 Bad Request`.
- provider/외부 시스템 오류: `502 Bad Gateway`.
- 서비스 레이어의 sentinel error(`services.ErrUnknownProvider`)는 핸들러에서 HTTP 상태로 매핑합니다.
- 응답 포맷은 기존 엔드포인트와 동일한 JSON 형태를 유지합니다.

## 8. 테스트 규칙

- 테스트 함수명: `Test<대상>_<시나리오>`.
- 실패 메시지는 `got/want` 형태로 명확히 작성합니다.
- 테이블 테스트(`[]struct{...}` + `t.Run`)를 우선 사용합니다.
- 외부 네트워크 의존 테스트는 기본적으로 mock 기반으로 작성합니다.
- 서브테스트에서 독립성이 보장되면 `t.Parallel()`을 사용합니다.

권장 명령:

```bash
go test ./...
```

## 9. 플러그인 구현 체크리스트

- `Load`에서 필수 설정값 검증 (`api_key`, `ipn_secret` 등).
- timeout/base URL 등 런타임 옵션은 설정 기반으로 주입.
- provider 응답을 `internal/models`로 매핑할 때 필드 정규화(통화코드 대문자 등)를 일관되게 적용.
- 웹훅 서명 검증 로직은 우회 경로 없이 실패 시 즉시 에러 반환.

## 10. 문서/주석 규칙

- 주석은 "무엇"보다 "왜"가 필요한 경우에만 작성합니다.
- exported 타입/함수는 GoDoc 스타일 설명을 유지합니다.
- 동작 변경 시 README 또는 `docs/` 문서 업데이트를 함께 수행합니다.
