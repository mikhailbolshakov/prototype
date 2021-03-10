module gitlab.medzdrav.ru/prototype/webrtc

go 1.15

//must be substitute with an external dependency once splpitted by repositories
replace gitlab.medzdrav.ru/prototype/kit => ../kit

replace gitlab.medzdrav.ru/prototype/proto => ../proto

replace github.com/coreos/etcd => go.etcd.io/etcd/v3 v3.5.0-alpha.0

require (
	github.com/armon/go-metrics v0.3.4 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/koding/websocketproxy v0.0.0-20181220232114-7ed82d81a28c
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pion/ion-avp v1.8.1
	github.com/pion/ion-log v1.0.0
	github.com/pion/ion-sfu v1.9.3
	github.com/pion/rtcp v1.2.6
	github.com/pion/webrtc/v3 v3.0.10
	github.com/sourcegraph/jsonrpc2 v0.0.0-20200429184054-15c2290dcb37
	gitlab.medzdrav.ru/prototype/kit v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/proto v0.0.0-00010101000000-000000000000
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0
	go.uber.org/multierr v1.6.0 // indirect
	google.golang.org/genproto v0.0.0-20201007142714-5c0e72c5e71e // indirect
	google.golang.org/grpc v1.36.0
	gopkg.in/ini.v1 v1.62.0 // indirect
)
