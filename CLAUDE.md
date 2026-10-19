# coredns-plugin-unifi

CoreDNS plugin (`unifi`) that queries a UniFi Controller for DHCP clients and resolves their hostnames using network
domain names. Built on [unpoller/unifi](https://github.com/unpoller/unifi). Not a standalone binary — it compiles into
CoreDNS via `plugin.cfg`.

## Build & Test

Into CoreDNS: add `unifi:github.com/navistau/coredns-plugin-unifi` to CoreDNS's `plugin.cfg`, then `go generate && go build`
in the CoreDNS source directory.

```sh
go vet ./...
go test ./...
go test -run TestSetup                            # single test
go test -tags integration -v -timeout 180s .      # integration, needs Docker
golangci-lint run ./...
```

Go 1.25.9 (`mise.toml`).

## Architecture

- **setup.go** — registration (`init()`) and Corefile parsing (`controllerurl`, `username`, `password`, `ttl`, `refreshinterval`, `sites`, `fallthrough`). Zone origins via `plugin.OriginsFromArgsOrServerBlock`. Starts the background refresh goroutine.
- **unifi.go** — `ServeDNS()` filters by zone (`Origins`), answers A lookups from the in-memory `UnifiConfigEntryMap`. `refresh()` fetches sites/clients/networks, optionally filters by `sites`, rebuilds hostname→IP. `sanitizeHostname()` lowercases, maps spaces/underscores to hyphens, strips invalid chars. UI alias (`Name`) wins over DHCP `Hostname`. `sync.RWMutex`. Hostname collisions are logged and skipped.
- **client.go** — wraps `unpoller/unifi` behind the `UnifiAPI` interface; `NewUnifiClient()`.
- **ready.go** — readiness (always true).
- **metrics.go** — Prometheus counter `coredns_unifi_request_count_total`.
- **mock_test.go** — `mockUnifiAPI` for unit tests.
- **integration_test.go** / **integration/** — compose stack, Dockerfiles, mock controller, Corefile.

## Corefile

```
unifi {
  controllerurl http://controller:port/
  username svc_coredns
  password mysecretpassword
  refreshInterval 30
  ttl 30
  sites default,branch-office
  fallthrough
}
```

## Integration tests

- `integration_test.go` shells out to `docker compose`; `TestMain` runs `up -d --build --wait` (both services have healthchecks; CoreDNS's uses the `health` plugin), reads the DNS port via `compose port`, waits for the initial refresh, tears down after.
- Assertions use the `miekg/dns` client.
- Mock controller (`integration/mock-controller/`, own `go.mod`) uses unpoller/unifi types (`Client`, `Network`, `Site`, `ServerStatus`) and API path constants; serves 5 clients across 2 networks.
- Dockerfiles: base `golang:1.25-alpine` builder, `alpine:3.20` runtime; `go mod download` layer before source.
- Cases: named clients resolve, hostname fallback when Name empty, multi-network domain (`iot.lan`), unknown client → NXDOMAIN, nonexistent host → NXDOMAIN, case-insensitive lookup.
- unpoller's `GetClients` fills `Hostname` from the MAC when both `Name` and `Hostname` are empty, so the "nameless client" fixture resolves by its MAC.
