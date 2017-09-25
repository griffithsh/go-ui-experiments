# UI experiments in Go

## Goals

I want to be able to run a command to compile a static binary that includes ...

- html
- javascript
- images

... so that distribution is simple.

I want a command to launch the UI in Chrome so that I can debug complicated UI javascript issues.

## aop-level-editor

Install with:

```bash
go install ./tools/... && go generate ./... && go install ./cmd/aop-level-editor/...
```

Run the dev version with:

```bash
dev-server -db /path/to/database/file -root /path/to/client/src
```