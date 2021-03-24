package main

import (
	"github.com/adacta-ru/mattermost-server/v6/plugin"
)

func main() {
	plugin.ClientMain(NewPlugin())
}
