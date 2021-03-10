module gitlab.medzdrav.ru/prototype/services

go 1.15

//must be substitute with an external dependency once splpitted by repositories
replace gitlab.medzdrav.ru/prototype/kit => ../kit

replace gitlab.medzdrav.ru/prototype/proto => ../proto

require (
	github.com/armon/go-metrics v0.3.4 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/google/uuid v1.2.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c
	gitlab.medzdrav.ru/prototype/kit v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/proto v0.0.0-00010101000000-000000000000
	go.uber.org/multierr v1.6.0 // indirect
	google.golang.org/genproto v0.0.0-20201007142714-5c0e72c5e71e // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
)
