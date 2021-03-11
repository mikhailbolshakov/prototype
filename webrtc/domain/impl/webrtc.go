package impl

import (
	"context"
	"fmt"
	"github.com/pion/ion-sfu/pkg/buffer"
	"github.com/pion/ion-sfu/pkg/middlewares/datachannel"
	"github.com/pion/ion-sfu/pkg/sfu"
	"gitlab.medzdrav.ru/prototype/kit/common"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"gitlab.medzdrav.ru/prototype/webrtc/logger"
	"gitlab.medzdrav.ru/prototype/webrtc/meta"
	"sync"
)

type room struct {
	id         string
	sfuSession *sfu.Session
	meta       *domain.RoomMeta
	cancel     context.CancelFunc
	recorder   domain.RoomRecorder
}

type webrtcImpl struct {
	//sfu.SessionProvider
	sync.RWMutex
	common.BaseService
	sfu             *sfu.SFU
	roomService     domain.RoomService
	sfuTransportCfg sfu.WebRTCTransportConfig
	rooms           map[string]*room
	roomLeases      map[string]context.CancelFunc
	roomCoord       domain.RoomCoordinator
	cfg             *config.Config
	factory         *buffer.Factory
	dc              *sfu.Datachannel
	recording       domain.Recording
}

func NewWebrtcService(roomCoord domain.RoomCoordinator, roomService domain.RoomService, queue queue.Queue) domain.WebrtcService {

	s := &webrtcImpl{
		roomService: roomService,
		rooms:       make(map[string]*room),
		roomCoord:   roomCoord,
		roomLeases:  make(map[string]context.CancelFunc),
	}
	s.BaseService = common.BaseService{Queue: queue}

	return s
}

func (w *webrtcImpl) l() log.CLogger {
	return logger.L().Cmp("webrtc")
}

func (w *webrtcImpl) Init(ctx context.Context, cfg *config.Config) error {
	w.cfg = cfg
	w.sfu = sfu.NewSFU(*cfg.Webrtc.Pion)
	w.sfuTransportCfg = sfu.NewWebRTCTransportConfig(*cfg.Webrtc.Pion)
	w.factory = buffer.NewBufferFactory(cfg.Webrtc.Pion.Router.MaxPacketTrack)
	w.dc = w.sfu.NewDatachannel(sfu.APIChannelLabel)
	w.dc.Use(datachannel.SubscriberAPI)

	if w.cfg.Webrtc.Recording.File.Enabled {
		w.recording = NewRecording()
		return w.recording.Init(ctx, cfg, w)
	}
	return nil

}

func (w *webrtcImpl) GetSFU() *sfu.SFU {
	return w.sfu
}

func (w *webrtcImpl) NewPeer(ctx context.Context, userId, username string) domain.Peer {
	return newPeer(ctx, w.sfu, w.roomService, userId, username)
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

func (w *webrtcImpl) createLocal(ctx context.Context, roomId string, meta *domain.RoomMeta) (*room, error) {

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

	return r, nil

}

func (w *webrtcImpl) GetOrCreateRoom(ctx context.Context, roomId string) (*domain.RoomMeta, error) {

	l := w.l().C(ctx).Mth("get-or-create-room").F(log.FF{"room": roomId}).Dbg()

	// get persistent room
	rm := w.roomService.Get(ctx, roomId)
	if rm == nil || rm.ClosedAt != nil {
		err := fmt.Errorf("no opened room found")
		l.E(err).St().Err()
		return nil, err
	}

	// check if there is a local room
	if r, ok := w.getLocalRoom(roomId); ok {
		l.Dbg("local room found")
		return r.meta, nil
	}

	//var err error
	roomMeta := w.createMeta(roomId)

	// check if there is a remote room on another node
	ctx, cancelFunc := context.WithCancel(ctx)
	foundOnAnotherNode, err := w.roomCoord.GetOrCreate(ctx, roomMeta)
	if err != nil {
		return nil, err
	}

	if foundOnAnotherNode {
		l.DbgF("found in cluster. meta=%v", *roomMeta)
		return roomMeta, nil
	}

	// otherwise create a new local room
	w.Lock()
	w.roomLeases[roomId] = cancelFunc
	w.Unlock()

	r, err := w.createLocal(ctx, roomId, roomMeta)
	if err != nil {
		return nil, err
	}

	l.DbgF("local created. meta=%v", *roomMeta)

	// create a new recorder
	// TODO: to support MINIO recorder
	if w.cfg.Webrtc.Recording.File.Enabled {
		rec, err := w.recording.NewRoomRecorder(ctx, roomId)
		if err != nil {
			return nil, err
		}
		r.recorder = rec
	}

	return roomMeta, nil
}

//func (w *webrtcImpl) GetSession(roomId string) (*sfu.Session, sfu.WebRTCTransportConfig) {
//	return w.ensureLocal(roomId).sfuSession, w.sfuTransportCfg
//}

//TODO: handle room closed correctly
func (w *webrtcImpl) onRoomClosed(roomId string) {

	l := w.l().Mth("room-closed").F(log.FF{"room": roomId})

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

	ctx := kitContext.NewRequestCtx().WithNewRequestId().Webrtc().ToContext(nil)
	_, err := w.roomService.Close(ctx, roomId)
	if err != nil {
		l.E(err).St().Err("persistence error")
		return
	}

	l.Dbg("closed")
}
