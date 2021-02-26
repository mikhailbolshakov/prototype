package ion

import (
	"context"
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"sync"

	"github.com/pborman/uuid"
	"github.com/pion/ion-sfu/pkg/sfu"
)

var (
	errNonLocalSession = errors.New("session is not located on this node")
)

type roomMeta struct {
	RoomID       string `json:"room_id"`
	NodeID       string `json:"node_id"`
	NodeEndpoint string `json:"node_endpoint"`
	Redirect     bool   `json:"redirect"`
}

// Coordinator is responsible for managing rooms
// and providing connections to other nodes
type coordinator interface {
	getOrCreateRoom(ctx context.Context, roomId string) (*roomMeta, error)
	sfu.SessionProvider
}

// newCoordinator configures coordinator for this node
func newCoordinator(ctx context.Context, cfg *config.Config) (coordinator, error) {

	webrtcCfg := cfg.Webrtc
	if webrtcCfg.Coordinator.Etcd != nil {
		return newCoordinatorEtcd(ctx, webrtcCfg, cfg.Etcd)
	}
	if webrtcCfg.Coordinator.Local != nil {
		return newCoordinatorLocal(webrtcCfg)
	}
	return nil, fmt.Errorf("error no coodinator configured")
}

type localCoordinator struct {
	nodeID       string
	nodeEndpoint string
	mu    sync.Mutex
	w     sfu.WebRTCTransportConfig
	rooms map[string]*sfu.Session
}

func newCoordinatorLocal(conf *config.Webrtc) (coordinator, error) {
	w := sfu.NewWebRTCTransportConfig(*conf.SFU)
	return &localCoordinator{
		nodeID:       uuid.New(),
		nodeEndpoint: endpoint(conf),
		rooms:        make(map[string]*sfu.Session),
		w:            w,
	}, nil
}

func (c *localCoordinator) ensureSession(roomId string) *sfu.Session {
	c.mu.Lock()
	defer c.mu.Unlock()

	if s, ok := c.rooms[roomId]; ok {
		return s
	}

	s := sfu.NewSession(roomId)
	s.OnClose(func() {
		c.onRoomClosed(roomId)
	})
	prometheusGaugeRooms.Inc()

	c.rooms[roomId] = s
	return s
}

func (c *localCoordinator) GetSession(roomId string) (*sfu.Session, sfu.WebRTCTransportConfig) {
	return c.ensureSession(roomId), c.w
}

func (c *localCoordinator) getOrCreateRoom(ctx context.Context, roomId string) (*roomMeta, error) {
	c.ensureSession(roomId)

	return &roomMeta{
		RoomID:       roomId,
		NodeID:       c.nodeID,
		NodeEndpoint: c.nodeEndpoint,
		Redirect:     false,
	}, nil
}

func (c *localCoordinator) onRoomClosed(sessionID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	log.L().Cmp("ion").Mth("session-closed").DbgF("%v closed", sessionID)
	delete(c.rooms, sessionID)
	prometheusGaugeRooms.Dec()
}
