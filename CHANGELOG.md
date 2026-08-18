# Changelog

All notable changes to this project are documented in this file. Any change to
Corefile directive behaviour is called out regardless of size.

The format is [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and
this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Dependencies updated: CoreDNS 1.11.3 → 1.14.6 (library and the version the
  integration suite builds against), coredns/caddy 1.1.4, miekg/dns 1.1.72,
  testcontainers-go 0.44.0, x/net in the mock controller. Resolves the
  outstanding Dependabot security alerts.
- Toolchain: Go 1.24.9 → 1.25.9 (CoreDNS 1.14.6 requires it); Docker builds
  use `golang:1.25-alpine`.
- CI: actions/checkout v7, actions/setup-go v7, golangci-lint-action v9;
  `.golangci.yml` migrated to the v2 config format; workflows now also run on
  pushes to `develop`.

## [0.4.0] - 2026-08-19

### Changed

- Module re-homed to `github.com/navistau/coredns-plugin-unifi` following the
  repository transfer from `jhogendorn` to the NavistAu org. Update your
  `plugin.cfg` entry; the old module path keeps working for old tags via
  GitHub's redirect, but new versions publish under the new path only.

### Fixed

- Install instructions: the step-by-step build was missing the `go get` step,
  and the Dockerfile example carried leftover local-build lines.

### Added

- Documented credentials via environment variables (`{$UNIFI_USERNAME}` /
  `{$UNIFI_PASSWORD}`), proven by the integration harness.

## [0.3.1] - 2026-06-10

### Fixed

- Setup no longer contacts the controller; the session initializes lazily on
  first refresh, so CoreDNS starts cleanly when the controller is unreachable.

## [0.3.0] - 2026-05-15

### Added

- Multiple A records for clients present on multiple interfaces.

## [0.2.0] - 2026-05-15

### Added

- PTR record support for IPv4 reverse lookups.

## [0.1.0] - 2026-04-20

### Added

- Initial release: A records for UniFi DHCP clients, built from controller
  client names and network domain names, with periodic refresh, site
  filtering, TTL and refresh interval directives, and fallthrough.
