<!--
Feature branches target `develop`. Only the release PR targets `main`.
See CONTRIBUTING.md for the branch model.
-->

## What this changes

<!-- One or two sentences. What behaviour is different after this merge? -->

## Why

<!-- The problem this solves. Link the issue if there is one. -->

Closes #

## Checklist

- [ ] `go vet ./...` and `go test ./...` pass.
- [ ] `go test -tags integration -timeout 300s .` passes. **This is the gate
      that counts** — it runs a real CoreDNS build against a mock controller.
- [ ] `golangci-lint run ./...` reports nothing.
- [ ] `CHANGELOG.md` has an entry under `## [Unreleased]`.
- [ ] Documentation is updated if the configuration surface changed.

## Configuration compatibility

<!--
Delete if not applicable. Does an existing Corefile keep working
unchanged? If not, say exactly what an operator must change, and confirm
the version bump reflects it.
-->
