package config

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

type Http struct {
	Host string
	Port string
}

type Mattermost struct {
	Url             string
	Ws              string
	AdminUsername   string `config:"admin-username"`
	AdminPassword   string `config:"admin-password"`
	DefaultPassword string `config:"default-password"`
	Team            string
	BotUsername     string `config:"bot-username"`
	BotAccessToken  string `config:"bot-access-token"`
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
}

type Config struct {
	Redis      *Redis
	Zeebe      *Zeebe
	Keycloak   *Keycloak
	Mattermost *Mattermost
	Services   map[string]*SrvCfg
	Http       *Http
	Test       string
}
