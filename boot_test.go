package unifi

import (
	"testing"

	"github.com/coredns/caddy"
)

// The plugin must never contact the controller during construction. A
// controller that is down (e.g. its host lost power) must not be able to
// fail plugin setup and crash-loop CoreDNS. The unpoller session is
// established lazily on the first refresh instead.

func TestNewUnifiClientDoesNoNetworkAtConstruction(t *testing.T) {
	// 192.0.2.1 is RFC5737 TEST-NET-1 — guaranteed unroutable. If the
	// constructor did network I/O it would error here.
	cfg := &UnifiConfig{
		controllerUrl:   "https://192.0.2.1:8443/",
		username:        "u",
		password:        "p",
		ttl:             defaultTTL,
		refreshInterval: defaultRefreshInterval,
	}
	c, err := NewUnifiClient(cfg)
	if err != nil {
		t.Fatalf("NewUnifiClient must not fail on an unreachable controller, got: %v", err)
	}
	if c == nil {
		t.Fatal("expected a client, got nil")
	}
	if c.api != nil {
		t.Fatal("api session must be created lazily (nil until first successful refresh)")
	}
}

func TestSetupSucceedsWhenControllerUnreachable(t *testing.T) {
	c := caddy.NewTestController("dns", `unifi 29c.sh {
		controllerurl https://192.0.2.1:8443/
		username u
		password p
	}`)
	if err := setup(c); err != nil {
		t.Fatalf("setup must succeed (CoreDNS must boot) when the controller is unreachable, got: %v", err)
	}
}
