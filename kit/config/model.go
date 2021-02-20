package config

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

type Es struct {
	Url   string
	Trace bool
}

type Simulcast struct {
	BestQualityFirst bool
}

type Router struct {
	MaxBandWidth  uint64
	MaxBufferTime int
	Simulcast     *Simulcast
}

type Sfu struct {
	Ballast int64
	Router  *Router
}

type IceServer struct {
	URLs       []string `mapstructure:"urls"`
	Username   string
	Credential string
}

type Candidates struct {
	IceLite    bool
	NAT1To1IPs []string
}

type Webrtc struct {
	Sfu          *Sfu
	PortRange    []uint16
	SdpSemantics string
	Candidates   *Candidates
	IceServers   []*IceServer
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
	Test       string
}
