module gitlab.medzdrav.ru/prototype

go 1.15

replace github.com/coreos/etcd => go.etcd.io/etcd/v3 v3.5.0-alpha.0

require (
	github.com/Nerzal/gocloak/v7 v7.11.0
	github.com/adacta-ru/mattermost-server/v6 v6.0.0
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/go-co-op/gocron v0.4.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/joho/godotenv v1.3.0
	github.com/koding/websocketproxy v0.0.0-20181220232114-7ed82d81a28c
	github.com/mitchellh/mapstructure v1.2.2
	github.com/nats-io/nats-streaming-server v0.19.0 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/stan.go v0.7.0
	github.com/olivere/elastic v6.2.35+incompatible
	github.com/olivere/elastic/v7 v7.0.22
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/pion/ion-avp v1.8.1
	github.com/pion/ion-log v1.0.0
	github.com/pion/ion-sdk-go v0.4.0
	github.com/pion/ion-sfu v1.9.3
	github.com/pion/mediadevices v0.1.17
	github.com/pion/rtcp v1.2.6
	github.com/pion/webrtc/v3 v3.0.11
	github.com/satori/go.uuid v1.2.0
	github.com/sherifabdlnaby/configuro v0.0.2
	github.com/sirupsen/logrus v1.7.0
	github.com/sourcegraph/jsonrpc2 v0.0.0-20200429184054-15c2290dcb37
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c
	github.com/zeebe-io/zeebe/clients/go v0.26.0
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0
	go.uber.org/atomic v1.7.0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
	gorm.io/driver/postgres v1.0.6
	gorm.io/gorm v1.20.8
)
