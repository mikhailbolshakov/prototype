package ion

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/config"
	"strings"
)

func endpoint(c *config.Webrtc) string {
	port := strings.Split(c.Signal.HTTPAddr, ":")[1]
	if c.Signal.Key != "" && c.Signal.Cert != "" {
		return fmt.Sprintf("wss://%v:%v/webrtc", c.Signal.FQDN, port)
	}
	return fmt.Sprintf("ws://%v:%v/webrtc", c.Signal.FQDN, port)
}


