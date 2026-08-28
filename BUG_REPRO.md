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
?   	courseattachments/audit	[no test files]
?   	courseattachments/cmd/attachmentctl	[no test files]
--- FAIL: TestCancelledAttachmentSummaryStops (0.02s)
    cancellation_regression_test.go:35: expected unprocessed after cancellation, got complete
FAIL
FAIL	courseattachments	0.027s
ok  	courseattachments/cli	0.016s
ok  	courseattachments/domain	0.001s
ok  	courseattachments/fixture	0.001s
ok  	courseattachments/ingest	0.011s
ok  	courseattachments/review	0.011s
ok  	courseattachments/search	0.010s
ok  	courseattachments/store	0.012s
--- FAIL: TestCancelledAttachmentSummaryStops (0.02s)
    summary_test.go:26: cancelled task completed
FAIL
FAIL	courseattachments/summary	0.017s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/attachmentctl): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/attachmentctl): exit `0`
- Frontend build (web): exit `0`
