package impl

import (
	"context"
	"github.com/pion/ion-sfu/pkg/buffer"
	"github.com/pion/ion-sfu/pkg/middlewares/datachannel"
	"github.com/pion/ion-sfu/pkg/sfu"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/meta"
	"sync"
)

type room struct {
	id         string
	sfuSession *sfu.Session
	meta       *domain.RoomMeta
	cancel     context.CancelFunc
}

type webrtcImpl struct {
	sfu.SessionProvider
	sync.RWMutex
	common.BaseService
	sfu             *sfu.SFU
	storage         domain.WebrtcStorage
	sfuTransportCfg sfu.WebRTCTransportConfig
	rooms           map[string]*room
	roomLeases      map[string]context.CancelFunc
	roomCoord       domain.RoomCoordinator
	cfg             *config.Webrtc
	factory         *buffer.Factory
	dc              *sfu.Datachannel
}

func NewWebrtcService(cfg *config.Config, roomCoord domain.RoomCoordinator, storage domain.WebrtcStorage, queue queue.Queue) domain.WebrtcService {

	s := &webrtcImpl{
		cfg:             cfg.Webrtc,
		storage:         storage,
		sfu:             sfu.NewSFU(*cfg.Webrtc.Pion),
		sfuTransportCfg: sfu.NewWebRTCTransportConfig(*cfg.Webrtc.Pion),
		rooms:           make(map[string]*room),
		roomCoord:       roomCoord,
		roomLeases:      make(map[string]context.CancelFunc),
		factory:         buffer.NewBufferFactory(cfg.Webrtc.Pion.Router.MaxPacketTrack),
	}
	s.BaseService = common.BaseService{Queue: queue}

	s.dc = s.sfu.NewDatachannel(sfu.APIChannelLabel)
	s.dc.Use(datachannel.SubscriberAPI)

	return s
}

func (w *webrtcImpl) GetSFU() *sfu.SFU {
	return w.sfu
}

func (w *webrtcImpl) NewPeer(ctx context.Context) domain.Peer {
	return newPeer(ctx, w.sfu)
	//return newPeer(ctx, w)
}

func (w *webrtcImpl) createMeta(roomId string) *domain.RoomMeta {
	return &domain.RoomMeta{
		Id:       roomId,
		Endpoint: meta.Endpoint,
		NodeId:   meta.NodeId,
	}
}

func (w *webrtcImpl) getLocalRoom(roomId string) (*room, bool) {
	w.RLock()
	defer w.RUnlock()
	r, ok := w.rooms[roomId]
	return r, ok
}

func (w *webrtcImpl) createLocal(roomId string, meta *domain.RoomMeta) *room {

	//otherwise create a new room

	//bufferFactory := buffer.NewBufferFactory(w.cfg.Pion.Router.MaxPacketTrack)
	//dc := w.sfu.NewDatachannel(sfu.APIChannelLabel)
	//dc.Use(datachannel.SubscriberAPI)
	//
	//w.cfg.Pion.BufferFactory = bufferFactory
	//t := sfu.NewWebRTCTransportConfig(*w.cfg.Pion)
	//
	//sfuS := sfu.NewSession(roomId, nil, []*sfu.Datachannel { dc }, t)
	//sfuS.OnClose(func() {
	//	w.onRoomClosed(roomId)
	//})

	r := &room{
		id: roomId,
		//sfuSession: sfuS,
		meta: meta,
	}

	w.Lock()
	w.rooms[roomId] = r
	defer w.Unlock()

	return r

}

func (w *webrtcImpl) ensureLocal(roomId string) *room {

	if r, ok := w.getLocalRoom(roomId); ok {
		return r
	}

	return w.createLocal(roomId, w.createMeta(roomId))

}

func (w *webrtcImpl) GetOrCreateRoom(ctx context.Context, roomId string) (*domain.RoomMeta, error) {

	// check if there is a local room
	if r, ok := w.getLocalRoom(roomId); ok {
		return r.meta, nil
	}

	//var err error
	roomMeta := w.createMeta(roomId)

	// check if there is a remote room on another node
	//ctx, cancelFunc := context.WithCancel(ctx)
	//foundOnAnotherNode, err := w.roomCoord.GetOrCreate(ctx, roomMeta)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if foundOnAnotherNode {
	//	return roomMeta, nil
	//}
	//
	//// otherwise create a new local room
	//w.Lock()
	//w.roomLeases[roomId] = cancelFunc
	//w.Unlock()

	_ = w.createLocal(roomId, roomMeta)

	return roomMeta, nil

}

func (w *webrtcImpl) GetSession(roomId string) (*sfu.Session, sfu.WebRTCTransportConfig) {
	return w.ensureLocal(roomId).sfuSession, w.sfuTransportCfg
}

func (w *webrtcImpl) onRoomClosed(roomId string) {

	l := log.L().Cmp("webrtc").Mth("room-closed").F(log.FF{"room": roomId})

	w.roomCoord.Close(context.Background(), roomId)

	w.Lock()
	defer w.Unlock()

	if cancel, ok := w.roomLeases[roomId]; ok {
		delete(w.roomLeases, roomId)
		cancel()
	}

	if _, ok := w.rooms[roomId]; ok {
		delete(w.rooms, roomId)
	}

	l.Dbg("closed")
}
