package unifi

import (
	unpoller_unifi "github.com/unpoller/unifi"
)

// UnifiAPI defines the methods we use from the unpoller client,
// allowing us to mock the controller in tests.
type UnifiAPI interface {
	GetSites() ([]*unpoller_unifi.Site, error)
	GetClients(sites []*unpoller_unifi.Site) ([]*unpoller_unifi.Client, error)
	GetNetworks(sites []*unpoller_unifi.Site) ([]unpoller_unifi.Network, error)
}

type UnifiClient struct {
	controllerUrl string
	config        *UnifiConfig
	api           UnifiAPI
	// owned records that ensureAPI built this session, so invalidateAPI may
	// discard it. An api injected by a test is never owned and is never
	// discarded — the test supplied it and expects it to stay.
	owned bool
}

// NewUnifiClient builds the client WITHOUT contacting the controller.
//
// The underlying unpoller session logs in over HTTP at construction
// (checkNewStyleAPI/Login/GetServerData). Doing that here would make plugin
// setup — and therefore CoreDNS startup — depend on controller reachability:
// when the controller's host is down, NewUnifi blocks/errors, setup() returns
// the error, and CoreDNS crash-loops, taking ALL DNS down. So the session is
// established lazily on the first refresh via ensureAPI instead. A down
// controller now only costs the (optional) UniFi client records, never the
// resolver itself.
func NewUnifiClient(cfg *UnifiConfig) (*UnifiClient, error) {
	return &UnifiClient{
		controllerUrl: cfg.controllerUrl,
		config:        cfg,
	}, nil
}

// ensureAPI lazily establishes the authenticated unpoller session. It is a
// no-op once connected, and when an api has been injected (tests). On failure
// it returns the error so refresh() logs and retries on the next tick —
// without ever aborting startup. Only ever called from the single refresh
// goroutine, so no locking is needed around c.api.
func (c *UnifiClient) ensureAPI() error {
	if c.api != nil {
		return nil
	}
	unpollerConfig := &unpoller_unifi.Config{
		User:     c.config.username,
		Pass:     c.config.password,
		URL:      c.config.controllerUrl,
		ErrorLog: log.Warningf,
		DebugLog: log.Debugf,
	}
	client, err := unpoller_unifi.NewUnifi(unpollerConfig)
	if err != nil {
		return err
	}
	c.api = client
	c.owned = true
	return nil
}

// invalidateAPI drops the cached session so the next ensureAPI logs in again.
//
// The unpoller client authenticates once, in NewUnifi, and never re-logs in:
// its do() turns every non-200 into a bare "invalid status code from server"
// with no 401 branch. So when the controller expires the session cookie, the
// client is poisoned for the life of the process — ensureAPI short-circuits on
// a non-nil c.api, and every subsequent refresh 401s against a dead session.
//
// Observed on a live controller (UniFi Network 10.6.101): the session lasted
// roughly ten hours, then produced 2880 identical 401s a day — one per
// refreshinterval — for eleven days, until CoreDNS was restarted. The records
// the plugin already held stayed frozen at their last good refresh the whole
// time, so the failure is silent to anything but the log.
//
// Called only from the single refresh goroutine, same as ensureAPI, so c.api
// needs no locking.
func (c *UnifiClient) invalidateAPI() {
	if !c.owned {
		return
	}
	c.api = nil
	c.owned = false
}
