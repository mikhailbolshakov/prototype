module gitlab.medzdrav.ru/prototype/api

go 1.15

//must be substitute with an external dependency once splpitted by repositories
replace gitlab.medzdrav.ru/prototype/kit => ../kit

replace gitlab.medzdrav.ru/prototype/proto => ../proto

require (
	github.com/gorilla/mux v1.8.0
	gitlab.medzdrav.ru/prototype/kit v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/proto v0.0.0-00010101000000-000000000000
)
