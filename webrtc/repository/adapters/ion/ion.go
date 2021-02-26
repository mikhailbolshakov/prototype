package ion

import "github.com/pion/ion-sfu/pkg/sfu"

type ionImpl struct {
	sfu         *sfu.SFU
	coordinator coordinator
	signal      *Signal
}


