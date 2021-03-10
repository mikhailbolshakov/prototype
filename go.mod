module gitlab.medzdrav.ru/prototype

go 1.15

replace gitlab.medzdrav.ru/prototype/proto => ./proto

replace gitlab.medzdrav.ru/prototype/kit => ./kit

replace gitlab.medzdrav.ru/prototype/api => ./api

replace gitlab.medzdrav.ru/prototype/bp => ./bp

replace gitlab.medzdrav.ru/prototype/chat => ./chat

replace gitlab.medzdrav.ru/prototype/config => ./config

replace gitlab.medzdrav.ru/prototype/services => ./services

replace gitlab.medzdrav.ru/prototype/sessions => ./sessions

replace gitlab.medzdrav.ru/prototype/tasks => ./tasks

replace gitlab.medzdrav.ru/prototype/users => ./users

replace gitlab.medzdrav.ru/prototype/webrtc => ./webrtc

replace github.com/coreos/etcd => go.etcd.io/etcd/v3 v3.5.0-alpha.0

require (
	github.com/gorilla/websocket v1.4.2
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/pion/ion-log v1.0.0
	github.com/pion/ion-sdk-go v0.4.0
	github.com/pion/mediadevices v0.1.17
	github.com/pion/webrtc/v3 v3.0.11
	gitlab.medzdrav.ru/prototype/api v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/bp v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/chat v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/config v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/kit v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/services v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/sessions v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/tasks v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/users v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/webrtc v0.0.0-00010101000000-000000000000
)
