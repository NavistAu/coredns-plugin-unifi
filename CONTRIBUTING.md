# Contributing

Thanks for your interest. Bug reports, documentation fixes, and small focused
changes are all welcome.

Read [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) first. To report a
vulnerability, follow [`SECURITY.md`](SECURITY.md) instead of opening an issue.

## Before you start

This plugin is deliberately small. It fetches DHCP clients and networks from a
UniFi controller on a timer and serves matching A and PTR records from an
in-memory map. It tries to do nothing else.

If your change adds configuration surface, open an issue first. The bar for
new options is "CoreDNS's own plugins and the controller's settings cannot
express this", not "this would be convenient".

## Set up

```sh
git clone https://github.com/NavistAu/coredns-plugin-unifi.git
cd coredns-plugin-unifi
mise install     # installs the Go toolchain pinned in mise.toml
go vet ./...
go test ./...
```

You need Docker with the Compose plugin for the integration suite.

## The two test suites, and why one of them decides

```sh
go test ./...                             # unit tests against a mock API
go test -tags integration -timeout 300s . # real CoreDNS build + mock controller
```

The unit suite is fast and covers the refresh and lookup logic. The
integration suite compiles the plugin into a real CoreDNS binary, runs it in
Docker against a mock controller, and asserts over actual DNS queries — it is
the gate that counts, because only it exercises `plugin.cfg` registration,
Corefile parsing, and the wire format.

## Branch model

Feature branches (`feat/…`, `fix/…`, `docs/…`) start from and target
`develop`. Only the release PR merges `develop` into `main`; every push to
`main` is a release.

Every PR needs a `CHANGELOG.md` entry under `## [Unreleased]`. Any change to
Corefile directive behaviour gets an entry regardless of size.
