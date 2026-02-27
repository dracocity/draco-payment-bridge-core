| 내부 Status          | BitPay           | CoinGate           | CoinPayments | NOWPayments                      |
| ------------------ | ---------------- | ------------------ | ------------ | -------------------------------- |
| **CREATED**        | —                | new                | —            | —                                |
| **PENDING**        | new              | pending            | 0            | waiting                          |
| **CONFIRMING**     | paid / confirmed | confirming         | 1 / 2        | confirming / confirmed / sending |
| **PARTIALLY_PAID** | —                | —                  | (amount 비교) | partially_paid                   |
| **COMPLETED**      | complete         | paid               | 100          | finished                         |
| **FAILED**         | invalid          | invalid / canceled | -1 / -3      | failed                           |
| **EXPIRED**        | expired          | expired            | -2           | expired                          |
| **REFUNDED**       | —                | refunded           | —            | refunded                         |
