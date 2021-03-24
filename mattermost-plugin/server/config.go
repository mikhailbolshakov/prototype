package main

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration, as well as values computed from the configuration. Any public fields will be
// deserialized from the Mattermost server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// configuration can change at any time, access to the configuration must be synchronized. The
// strategy used in this plugin is to guard a pointer to the configuration, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
//
// If you add non-reference types to your configuration struct, be sure to rewrite Clone as a deep
// copy appropriate for your types.
type configuration struct {
	Disabled          bool
	NatsUrl           string
	NatsClusterId     string
	TopicNewPost      string
}

// Clone shallow copies the cfg. Your implementation may require a deep copy if
// your cfg has reference types.
func (c *configuration) Clone() *configuration {
	return &configuration{
		Disabled:          c.Disabled,
		NatsUrl:           c.NatsUrl,
		NatsClusterId:     c.NatsClusterId,
		TopicNewPost:      c.TopicNewPost,
	}
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.cfg == nil {
		return &configuration{}
	}

	return p.cfg
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing configuration. This almost
// certainly means that the configuration was modified without being cloned and may result in
// an unsafe access.
func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.cfg == configuration {
		// Ignore assignment if the cfg struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing cfg")
	}

	p.cfg = configuration
}

func (p *Plugin) ValidateConfiguration(cfg *configuration) error {
	if cfg.TopicNewPost == "" {
		return fmt.Errorf("[STAN] validation: topic is empty")
	}
	if cfg.NatsUrl == "" {
		return fmt.Errorf("[STAN] validation: nats url is empty")
	}
	if cfg.NatsClusterId == "" {
		return fmt.Errorf("[STAN] validation: cluster Id is empty")
	}
	return nil
}

// OnConfigurationChange is invoked when cfg changes may have been made.
func (p *Plugin) OnConfigurationChange() error {

	cfg := p.getConfiguration().Clone()

	// Load the public cfg fields from the Mattermost server cfg.
	if err := p.API.LoadPluginConfiguration(cfg); err != nil {
		return errors.Wrap(err, "[STAN] failed to load plugin cfg")
	}

	if err := p.ValidateConfiguration(cfg); err != nil {
		return err
	}

	p.setConfiguration(cfg)

	return nil
}
