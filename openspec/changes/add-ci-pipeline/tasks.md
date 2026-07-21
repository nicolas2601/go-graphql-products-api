## 1. Implementation

- [ ] 1.1 `.golangci.yml` (set estandar, excluye el codigo generado)
- [ ] 1.2 `.github/workflows/ci.yml` (lint; test con race + tidy check + govulncheck; docker build)

## 2. Verification

- [ ] 2.1 `golangci-lint run` local en verde
- [ ] 2.2 El workflow corre VERDE en GitHub Actions (confirmado con `gh run`/`gh pr checks`)
