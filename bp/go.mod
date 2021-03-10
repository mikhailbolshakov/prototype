module gitlab.medzdrav.ru/prototype/bp

go 1.15

//must be substitute with an external dependency once splpitted by repositories
replace gitlab.medzdrav.ru/prototype/kit => ../kit

replace gitlab.medzdrav.ru/prototype/proto => ../proto

require (
	github.com/Nerzal/gocloak/v7 v7.11.0
	github.com/golang/protobuf v1.4.3
	github.com/zeebe-io/zeebe/clients/go v0.26.1
	gitlab.medzdrav.ru/prototype/kit v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/proto v0.0.0-00010101000000-000000000000
)
