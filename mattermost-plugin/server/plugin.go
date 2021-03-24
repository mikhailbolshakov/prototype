package main

import (
	"github.com/adacta-ru/mattermost-server/v6/model"
	"github.com/adacta-ru/mattermost-server/v6/plugin"
	"github.com/nats-io/stan.go"
	"sync"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// stan connection
	stanConn stan.Conn

	// chan to close all goroutines
	close chan struct{}

	// configurationLock synchronizes access to the cfg.
	configurationLock sync.RWMutex

	// cfg is the active plugin cfg. Consult getConfiguration and
	// setConfiguration for usage.
	cfg *configuration
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
// It also creates a demo bot account
func (p *Plugin) OnActivate() error {

	p.close = make(chan struct{})

	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	if err := p.ConnectStan(); err != nil {
		return err
	}

	return nil
}

// OnDeactivate is invoked when the plugin is deactivated. This is the plugin's last chance to use
// the API, and the plugin will be terminated shortly after this invocation.
//
// This demo implementation logs a message to the demo channel whenever the plugin is deactivated.
func (p *Plugin) OnDeactivate() error {

	close(p.close)

	if err := p.StanClose(); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {

	p.API.LogDebug("[STAN] MessageHasBeenPosted")

	configuration := p.getConfiguration()

	if configuration.Disabled {
		return
	}

	p.StanPublishPost(post)


}


