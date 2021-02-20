package impl

import (
	"github.com/pion/ion-log"
	isfu "github.com/pion/ion-sfu/pkg"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
)

type sfu struct {
	ionSfu *isfu.SFU
}

func toSfuCfg(cfg *kitConfig.Config) isfu.Config {

	sfuCfg := isfu.Config{
		WebRTC: isfu.WebRTCConfig{
			ICEPortRange: cfg.Webrtc.PortRange,
			ICEServers:   []isfu.ICEServerConfig{},
			Candidates:   isfu.Candidates{
				IceLite:    cfg.Webrtc.Candidates.IceLite,
				NAT1To1IPs: cfg.Webrtc.Candidates.NAT1To1IPs,
			},
			SDPSemantics: cfg.Webrtc.SdpSemantics,
		},
		Log:    log.Config{
			Level: cfg.Services["webrtc"].Log.Level,
		},
		Router: isfu.RouterConfig{
			MaxBandwidth:  cfg.Webrtc.Sfu.Router.MaxBandWidth,
			MaxBufferTime: cfg.Webrtc.Sfu.Router.MaxBufferTime,
			Simulcast:     isfu.SimulcastConfig{
				BestQualityFirst:    cfg.Webrtc.Sfu.Router.Simulcast.BestQualityFirst,
				EnableTemporalLayer: false,
			},
		},
	}
	sfuCfg.SFU.Ballast = cfg.Webrtc.Sfu.Ballast

	if cfg.Webrtc.IceServers != nil {
		for _, ic := range cfg.Webrtc.IceServers {

			sfuIc := isfu.ICEServerConfig{
				URLs:       ic.URLs,
				Username:   ic.Username,
				Credential: ic.Credential,
			}

			sfuCfg.WebRTC.ICEServers = append(sfuCfg.WebRTC.ICEServers, sfuIc)

		}
	}

	return sfuCfg
}

func newSfu(cfg *kitConfig.Config) *sfu {
	return &sfu{
		ionSfu: isfu.NewSFU(toSfuCfg(cfg)),
	}
}
