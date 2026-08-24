# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	campgoods/cmd/campgoods	[no test files]
?   	campgoods/workflow	[no test files]
--- FAIL: TestProductPageKeepsContiguousItems (0.02s)
    regression_test.go:29: position 0: got product-21 want product-11
FAIL
FAIL	campgoods	0.076s
ok  	campgoods/catalog	0.002s
ok  	campgoods/inventory	0.033s
ok  	campgoods/pricing	0.005s
ok  	campgoods/query	0.005s
ok  	campgoods/reporting	0.029s
ok  	campgoods/store	0.013s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/campgoods): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/campgoods): exit `0`
