package config

import "github.com/pion/ion-sfu/pkg/sfu"

type Log struct {
	Level string
}

type Redis struct {
	Port     string
	Host     string
	Password string
	Ttl      int
}

type Zeebe struct {
	Port string
	Host string
}

type Keycloak struct {
	Url           string
	AdminUsername string `config:"admin-username"`
	AdminPassword string `config:"admin-password"`
	ClientId      string `config:"client-id"`
	ClientSecret  string `config:"client-secret"`
	ClientRealm   string `config:"client-realm"`
}

type Tls struct {
	Cert string
	Key  string
}

type Http struct {
	Host string
	Port string
	Tls  *Tls
}

type Mattermost struct {
	Url              string
	Ws               string
	AdminUsername    string `config:"admin-username"`
	AdminPassword    string `config:"admin-password"`
	AdminAccessToken string `config:"admin-access-token"`
	DefaultPassword  string `config:"default-password"`
	Team             string
	BotUsername      string `config:"bot-username"`
	BotAccessToken   string `config:"bot-access-token"`
}

type Database struct {
	Dbname   string
	User     string
	Password string
	Port     string
	HostRw   string `config:"host-rw"`
	HostRo   string `config:"host-ro"`
}

type Grpc struct {
	Port string
	Host string
}

type SrvCfg struct {
	Database *Database
	Grpc     *Grpc
	Log      *Log
}

type Etcd struct {
	Hosts []string
}

type Es struct {
	Url   string
	Trace bool
}

type AuthConfig struct {
	Enabled bool
	Key     string
	KeyType string
}

type SignalConfig struct {
	FQDN     string
	Key      string
	Cert     string
	HTTPAddr string
	GRPCAddr string
	Auth     *AuthConfig
}

type CoordinatorConfig struct {
	Local *struct {
		Enabled bool
	}
	Etcd *struct {
		Enabled bool
	}
}

type Webrtc struct {
	Signal      *SignalConfig
	SFU         *sfu.Config
	Coordinator *CoordinatorConfig
}

type Nats struct {
	Url       string
	ClusterId string
}

type Config struct {
	Redis      *Redis
	Zeebe      *Zeebe
	Keycloak   *Keycloak
	Mattermost *Mattermost
	Services   map[string]*SrvCfg
	Http       *Http
	Es         *Es
	Webrtc     *Webrtc
	Nats       *Nats
	Etcd       *Etcd
	Test       string
}
