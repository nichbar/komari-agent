# Project Contract

## Build And Test

- Install: `go mod download`
- Dev / Run: `go run . [flags]`
- Build: `go build -trimpath -o komari-agent .`
- Multi-arch build: `./build_all.sh`
- Test: `go test -v ./update ./monitoring/unit/...`
- Vet / Typecheck: `go vet ./...`
- Lint / Format: `gofmt -s -w .`

## Architecture Boundaries

- CLI entrypoint and commands live in `cmd/` (flags in `cmd/flags/`)
- System metric collectors live in `monitoring/unit/`
- Report aggregation and JSON formatting live in `monitoring/`
- Server communication, remote task execution, and latency pinging live in `server/`
- PTY / Web SSH session handling lives in `terminal/`
- Thread-safe WebSocket synchronization lives in `ws/`
- Self-update logic and version management live in `update/`
- Custom DNS and network transport overrides live in `dnsresolver/`
- Platform-specific code MUST live in platform-suffixed files (e.g., `*_linux.go`, `*_darwin.go`, `*_windows.go`, `*_freebsd.go`) with proper build constraints

## Coding Conventions

- Ensure cross-platform compatibility across Linux, macOS, Windows, FreeBSD, and Android
- Always synchronize WebSocket writes through `ws.SafeConn` to prevent concurrent write panics
- Check and respect security flags (`--disable-web-ssh`, `--disable-auto-update`) before invoking remote operations
- Propagate authentication tokens and Cloudflare Access headers (`CF-Access-Client-Id`, `CF-Access-Client-Secret`) across all server HTTP/WS calls
- Gracefully handle missing host dependencies (e.g., fallback if `vnstat` is unavailable) without panicking
- Maintain backward compatibility for deprecated CLI flags in `cmd/root.go`

## Safety Rails

## NEVER

- Modify `go.mod` or `go.sum` without explicit approval
- Perform raw unsynchronized writes on `websocket.Conn`
- Bypass `flags.DisableWebSsh` checks when handling terminal or remote exec requests
- Commit without running `go vet ./...` and verifying the build
- Hardcode platform-dependent filesystem paths or shell invocations in shared files
- Expose tokens, auto-discovery keys, or Cloudflare Access credentials in logs or commit history

## ALWAYS

- Show diff before committing
- Verify cross-platform compilation (`GOOS=linux`, `GOOS=windows`, `GOOS=darwin`)
- Use `ws.SafeConn` for all WebSocket communication
- Keep command flags organized within `cmd/flags/`

## Verification

- Lint and static analysis: `gofmt -s -d .` + `go vet ./...`
- Local build: `CGO_ENABLED=0 go build -trimpath .`
- Cross-platform build check:
  - Linux: `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath .`
  - Windows: `GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath .`
  - macOS: `GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath .`
- Unit tests: `go test -v ./update`
- CLI sanity check: `./komari-agent --help` and `./komari-agent list-disk`

## Compact Instructions

Preserve:

1. Architecture decisions (NEVER summarize)
2. Modified files and key changes
3. Current verification status (pass/fail commands)
4. Open risks, TODOs, rollback notes
