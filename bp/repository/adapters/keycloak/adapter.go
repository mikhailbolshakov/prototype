package keycloak

import (
	"github.com/Nerzal/gocloak/v7"
	domain "gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/proto/config"
)

type Adapter interface {
	Init(c *config.Config) error
	GetProvider() domain.KeycloakProvider
	Close()
}

type adapterImpl struct {
	keycloak gocloak.GoCloak
}

func NewAdapter() Adapter {
	c := &adapterImpl{}
	return c
}

func (c *adapterImpl) Init(cfg *config.Config) error {
	c.keycloak = gocloak.NewClient(cfg.Keycloak.Url)
	return nil
}

func (c *adapterImpl) GetProvider() domain.KeycloakProvider {
	return func() gocloak.GoCloak {
		return c.keycloak
	}
}

func (c *adapterImpl) Close() {
}
