package keycloak

import (
	"github.com/Nerzal/gocloak/v7"
	domain "gitlab.medzdrav.ru/prototype/bp/domain"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
)

type Adapter interface {
	Init(c *kitConfig.Config) error
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

func (c *adapterImpl) Init(cfg *kitConfig.Config) error {
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
