# Security policy

## Supported versions

The most recent minor release gets security fixes. This project is small and
has no long-term support branches.

| Version | Supported |
| ------- | --------- |
| 0.4.x   | Yes       |
| < 0.4   | No        |

## Report a vulnerability

Report vulnerabilities privately through
[GitHub Security Advisories](https://github.com/NavistAu/coredns-plugin-unifi/security/advisories/new).
Do not open a public issue for a vulnerability.

Include the plugin version, the CoreDNS version, a Corefile that shows the
problem, and the effect you can demonstrate. Expect an acknowledgement within
seven days.

## What is in scope

This plugin holds controller credentials and answers DNS inside the zones you
configure. Examples of in-scope reports:

- The plugin answers a query for a zone the server block does not declare.
- The controller password appears in logs, metrics, or DNS responses.
- A controller response that crashes CoreDNS or hangs the refresh loop.
- A client name that survives sanitization and produces a record outside the
  network's domain.

## What is out of scope

- **A controller you do not trust.** The plugin serves whatever client and
  network data the controller returns. Anyone who can write client aliases on
  the controller can choose the names this plugin serves.
- Vulnerabilities in CoreDNS, the UniFi controller, or dependencies — report
  those upstream.
- DNS-layer attacks the plugin cannot influence (spoofing of upstream
  resolvers, cache poisoning elsewhere in the chain).
