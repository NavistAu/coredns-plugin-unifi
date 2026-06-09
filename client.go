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
	return nil
}
