package avp

import (
	"fmt"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"sync"
	"time"

	avp "github.com/pion/ion-avp/pkg"
	log "github.com/pion/ion-log"
	"github.com/pion/rtcp"

	"github.com/pion/webrtc/v3"
)

// this is ion-avp/pkg/webrtctransport.go fork

const (
	publisher  = 0
	subscriber = 1
)

type webRTCTransportConfig struct {
	configuration webrtc.Configuration
	setting       webrtc.SettingEngine
}

type pendingProcess struct {
	pid string
	fn  func() avp.Element
}

type AvpPeer struct {
	id  string
	pub *Publisher
	sub *Subscriber
	mu  sync.RWMutex

	builders  map[string]*avp.Builder     // one builder per track
	pending   map[string][]pendingProcess // maps track id to pending element constructors
	processes map[string]avp.Element      // existing processes
	onCloseFn func()

	webrtcService domain.WebrtcService
}

// NewWebRTCTransport creates a new webrtc transport
func NewAvpPeer(roomId string, c *avp.Config, defaultAvpElementFn avp.ElementFun) *AvpPeer {

	conf := webrtc.Configuration{}
	se := webrtc.SettingEngine{}

	var icePortStart, icePortEnd uint16

	if len(c.WebRTC.ICEPortRange) == 2 {
		icePortStart = c.WebRTC.ICEPortRange[0]
		icePortEnd = c.WebRTC.ICEPortRange[1]
	}

	if icePortStart != 0 || icePortEnd != 0 {
		if err := se.SetEphemeralUDPPortRange(icePortStart, icePortEnd); err != nil {
			panic(err)
		}
	}

	var iceServers []webrtc.ICEServer
	for _, iceServer := range c.WebRTC.ICEServers {
		s := webrtc.ICEServer{
			URLs:       iceServer.URLs,
			Username:   iceServer.Username,
			Credential: iceServer.Credential,
		}
		iceServers = append(iceServers, s)
	}

	conf.ICEServers = iceServers

	config := webRTCTransportConfig{
		setting:       se,
		configuration: conf,
	}

	pub, err := NewPublisher(config)
	if err != nil {
		log.Errorf("Error creating peer connection: %s", err)
		return nil
	}

	sub, err := NewSubscriber(config)
	if err != nil {
		log.Errorf("Error creating peer connection: %s", err)
		return nil
	}

	t := &AvpPeer{
		id:        roomId,
		pub:       pub,
		sub:       sub,
		builders:  make(map[string]*avp.Builder),
		pending:   make(map[string][]pendingProcess),
		processes: make(map[string]avp.Element),
	}

	sub.OnTrack(func(track *webrtc.TrackRemote, recv *webrtc.RTPReceiver) {

		id := track.ID()
		log.Debugf("Got track: %s", id)

		maxlate := c.SampleBuilder.AudioMaxLate
		if track.Kind() == webrtc.RTPCodecTypeVideo {
			maxlate = c.SampleBuilder.VideoMaxLate
		}

		if maxlate == 0 {
			log.Warnf("maxlate should not be 0. Using 100.")
			maxlate = 100
		}

		builder := avp.NewBuilder(track, maxlate)
		builder.AttachElement(defaultAvpElementFn(roomId, id, id, []byte{}))

		t.mu.Lock()
		defer t.mu.Unlock()
		t.builders[id] = builder

		// If there is a pending pipeline for this track,
		// initialize the pipeline.
		if pending := t.pending[id]; len(pending) != 0 {
			for _, p := range pending {
				process := t.processes[p.pid]
				if process == nil {
					process = p.fn()
					t.processes[p.pid] = process
				}
				builder.AttachElement(process)
			}
			delete(t.pending, id)
		}

		if track.Kind() == webrtc.RTPCodecTypeVideo {
			err := sub.pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{SenderSSRC: uint32(track.SSRC()), MediaSSRC: uint32(track.SSRC())}})
			if err != nil {
				log.Errorf("error writing pli %s", err)
			}
		}

		builder.OnStop(func() {
			t.mu.Lock()
			b := t.builders[id]
			if b != nil {
				log.Debugf("stop builder %s", id)
				delete(t.builders, id)
			}
			t.mu.Unlock()

			if t.isEmpty() {
				// No more tracks, cleanup transport
				t.Close()
			}
		})
	})

	go t.pliLoop(c.WebRTC.PLICycle)

	return t
}

func (t *AvpPeer) pliLoop(cycle uint) {
	if cycle == 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(cycle) * time.Millisecond)
	for range ticker.C {
		t.mu.RLock()
		builders := t.builders
		t.mu.RUnlock()

		if len(builders) == 0 {
			return
		}

		var pkts []rtcp.Packet
		for _, b := range builders {
			pkts = append(pkts, &rtcp.PictureLossIndication{SenderSSRC: uint32(b.Track().SSRC()), MediaSSRC: uint32(b.Track().SSRC())})
		}

		err := t.sub.pc.WriteRTCP(pkts)
		if err != nil {
			log.Errorf("error writing pli %s", err)
		}
	}
}

func (t *AvpPeer) isEmpty() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.builders) == 0 && len(t.pending) == 0
}

// OnClose sets a handler that is called when the webrtc transport is closed
func (t *AvpPeer) OnClose(f func()) {
	t.onCloseFn = f
}

// Close the webrtc transport
func (t *AvpPeer) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.onCloseFn != nil {
		t.onCloseFn()
	}

	err := t.sub.Close()
	if err != nil {
		return err
	}
	return t.pub.Close()
}

// CreateOffer starts the PeerConnection and generates the localDescription
func (t *AvpPeer) CreateOffer() (webrtc.SessionDescription, error) {
	return t.pub.CreateOffer()
}

// SetRemoteDescription sets the SessionDescription of the remote peer
func (t *AvpPeer) SetRemoteDescription(desc webrtc.SessionDescription) error {
	return t.pub.SetRemoteDescription(desc)
}

// Answer starts the PeerConnection and generates the localDescription
func (t *AvpPeer) Answer(offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {
	return t.sub.Answer(offer)
}

// AddICECandidate accepts an ICE candidate string and adds it to the existing set of candidates
func (t *AvpPeer) AddICECandidate(candidate webrtc.ICECandidateInit, target int) error {
	switch target {
	case publisher:
		if err := t.pub.AddICECandidate(candidate); err != nil {
			return fmt.Errorf("error setting ice candidate: %w", err)
		}
	case subscriber:
		if err := t.sub.AddICECandidate(candidate); err != nil {
			return fmt.Errorf("error setting ice candidate: %w", err)
		}
	}
	return nil
}

// OnICECandidate sets an event handler which is invoked when a new ICE candidate is found.
// Take note that the handler is gonna be called with a nil pointer when gathering is finished.
func (t *AvpPeer) OnICECandidate(f func(c *webrtc.ICECandidate, target int)) {
	t.pub.OnICECandidate(func(c *webrtc.ICECandidate) {
		f(c, publisher)
	})
	t.sub.OnICECandidate(func(c *webrtc.ICECandidate) {
		f(c, subscriber)
	})
}

func (t *AvpPeer) OnIceConnectionStateChanged(c webrtc.ICEConnectionState) {

}
