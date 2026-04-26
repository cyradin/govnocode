# govnocode

## Testing

This project separates unit and integration tests using Go build tags.

### Unit tests

Run normally:

```bash
go test ./...
```

Or via VS Code test runner.

### Integration tests

Marked with:

```go
//go:build integration
```

Run via CLI:

```bash
go test -tags=integration ./...
```

### VS Code setup

To enable running integration tests from VS Code, add to `.vscode/settings.json`:

```json
{
  "go.testTags": "integration",
  "go.buildTags": "integration",
  "gopls": {
    "buildFlags": ["-tags=integration"]
  }
}
```
