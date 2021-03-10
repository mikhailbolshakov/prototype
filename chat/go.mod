module gitlab.medzdrav.ru/prototype/chat

go 1.15

//must be substitute with an external dependency once splpitted by repositories
replace gitlab.medzdrav.ru/prototype/kit => ../kit

replace gitlab.medzdrav.ru/prototype/proto => ../proto

require (
	github.com/adacta-ru/mattermost-server/v6 v6.0.0
	gitlab.medzdrav.ru/prototype/kit v0.0.0-00010101000000-000000000000
	gitlab.medzdrav.ru/prototype/proto v0.0.0-00010101000000-000000000000
	go.uber.org/atomic v1.7.0
)
