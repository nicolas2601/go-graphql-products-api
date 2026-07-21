## 1. Implementation

- [x] 1.1 `.golangci.yml` (set estandar, excluye el codigo generado)
- [x] 1.2 `.github/workflows/ci.yml` (lint; test con race + tidy check + govulncheck; docker build)

## 2. Verification

- [x] 2.1 `golangci-lint run` local en verde (0 issues)
- [x] 2.2 El workflow corre VERDE en GitHub Actions (run 29851523651: jobs lint, test y docker todos en success)
