package api

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/pion/ion-log"
	sdk "github.com/pion/ion-sdk-go"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/vpx"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/ivfwriter"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"strings"
	//_ "github.com/pion/mediadevices/pkg/driver/audiotest"
	_ "github.com/pion/mediadevices/pkg/driver/camera"     // This is required to register camera adapter
	_ "github.com/pion/mediadevices/pkg/driver/microphone" // This is required to register microphone adapter
	// Note: If you don't have a camera or microphone or your adapters are not supported,
	//       you can always swap your adapters with our dummy adapters below.
	//_ "github.com/pion/mediadevices/pkg/driver/videotest"
)

func (h *TestHelper) WebrtcWs(sessionId, roomId string) (*websocket.Conn, chan struct{}, error) {

	wsConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s?session=%s&room=%s", WEBRTC_URL, sessionId, roomId), nil)
	if err != nil {
		return nil, nil, err
	}
	done := make(chan struct{})
	go h.ListenWs(wsConn, done)
	return wsConn, done, nil
}

type peer struct {
	id string
	c *sdk.Client
}

func (p *peer) save(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
	codec := track.Codec()
	var fileWriter media.Writer
	var err error
	if strings.EqualFold(codec.MimeType, webrtc.MimeTypeOpus) {
		log.Infof("Got Opus track, saving to disk as ogg (48 kHz, 2 channels)")
		fileWriter, err = oggwriter.New(fmt.Sprintf("%s_%d_%d.ogg", p.id, codec.PayloadType, track.SSRC()), 48000, 2)
	} else if strings.EqualFold(codec.MimeType, webrtc.MimeTypeVP8) {
		log.Infof("Got VP8 track, saving to disk as ivf")
		fileWriter, err = ivfwriter.New(fmt.Sprintf("%s_%d_%d.ivf", p.id, codec.PayloadType, track.SSRC()))
	}

	if err != nil {
		log.Errorf("err=%v", err)
		fileWriter.Close()
		return
	}

	for {
		rtpPacket, _, err := track.ReadRTP()
		if err != nil {
			panic(err)
		}
		if err := fileWriter.WriteRTP(rtpPacket); err != nil {
			panic(err)
		}
	}
}

var tracks []mediadevices.Track

func getTracks() []mediadevices.Track {

	if tracks != nil {
		return tracks
	}

	vpxParams, err := vpx.NewVP8Params()
	if err != nil {
		panic(err)
	}
	vpxParams.BitRate = 500_000 // 500kbps

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&vpxParams),
	)

	fmt.Println(mediadevices.EnumerateDevices())

	s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormat(frame.FormatYUY2)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
		},
		Codec: codecSelector,
	})

	if err != nil {
		panic(err)
	}

	tracks = s.GetTracks()

	return tracks

}

func (h *TestHelper) NewPeer(sid, peerId string, streamToFile bool) *peer {

	p := &peer{id: peerId}

	webrtcCfg := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			webrtc.ICEServer{
				URLs: []string{"stun:stun.stunprotocol.org:3478", "stun:stun.l.google.com:19302"},
			},
		},
	}

	config := sdk.Config{
		Log: log.Config{
			Level: "debug",
		},
		WebRTC: sdk.WebRTCTransportConfig{
			Configuration: webrtcCfg,
		},
	}
	// new sdk engine
	e := sdk.NewEngine(config)

	// get a client from engine
	p.c = e.AddClient(WEBRTC_GRPC, sid, p.id)

	// subscribe rtp from session
	if streamToFile {
		p.c.OnTrack = p.save
	}

	// client join a session
	if err := p.c.Join(sid); err != nil {
		panic(err)
	}

	for _, t := range getTracks() {
		if _, err := p.c.Publish(t); err != nil {
			panic(err)
		}
	}

	return p
}

func receiver(sid, peer string) {

	webrtcCfg := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			webrtc.ICEServer{
				URLs: []string{"stun:stun.stunprotocol.org:3478", "stun:stun.l.google.com:19302"},
			},
		},
	}

	config := sdk.Config{
		Log: log.Config{
			Level: "debug",
		},
		WebRTC: sdk.WebRTCTransportConfig{
			Configuration: webrtcCfg,
		},
	}
	// new sdk engine
	e := sdk.NewEngine(config)

	// get a client from engine
	c := e.AddClient(WEBRTC_GRPC, sid, peer)

	// subscribe rtp from sessoin
	// comment this if you don't need save to file
	//c.OnTrack = save

	// client join a session
	if err := c.Join(sid); err != nil {
		panic(err)
	}

}

func sender(sid, peer string) {

	// add stun servers
	webrtcCfg := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			webrtc.ICEServer{
				URLs: []string{"stun:stun.stunprotocol.org:3478", "stun:stun.l.google.com:19302"},
			},
		},
	}

	config := sdk.Config{
		Log: log.Config{
			Level: "debug",
		},
		WebRTC: sdk.WebRTCTransportConfig{
			Configuration: webrtcCfg,
		},
	}
	// new sdk engine
	e := sdk.NewEngine(config)

	// get a client from engine
	c := e.AddClient(WEBRTC_GRPC, sid, peer)

	// client join a session
	err := c.Join(sid)

	if err != nil {
		log.Errorf("err=%v", err)
		panic(err)
	}

	vpxParams, err := vpx.NewVP8Params()
	if err != nil {
		panic(err)
	}
	vpxParams.BitRate = 500_000 // 500kbps

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithVideoEncoders(&vpxParams),
	)

	fmt.Println(mediadevices.EnumerateDevices())

	s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormat(frame.FormatYUY2)
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
		},
		Codec: codecSelector,
	})

	if err != nil {
		panic(err)
	}

	for _, track := range s.GetTracks() {
		track.OnEnded(func(err error) {
			fmt.Printf("Track (ID: %s) ended with error: %v\n",
				track.ID(), err)
		})
		_, err = c.Publish(track)
		if err != nil {
			panic(err)
		} else {
			break // only publish first track, thanks
		}
	}

}
