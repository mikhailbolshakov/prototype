package impl

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	isfu "github.com/pion/ion-sfu/pkg"
	"github.com/pion/webrtc/v3"
	"gitlab.medzdrav.ru/prototype/kit"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/proto"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	"sync"
)

type peer struct {
	userId  string
	ionPeer *isfu.Peer
}

type room struct {
	sync.RWMutex
	id              string
	initiatorUserId string
	chatChannelId   string
	peers           map[string]*peer
}

type hub struct {
	sync.RWMutex
	rooms   map[string]*room
	storage domain.RoomStorage
	cfg     *kitConfig.Config
	sfu     *sfu
}

func (h *hub) addRoom(r *room) {
	h.Lock()
	defer h.Unlock()
	h.rooms[r.id] = r
}

func (h *hub) getRoom(id string) (*room, bool) {
	h.RLock()
	defer h.RUnlock()
	r, ok := h.rooms[id]
	return r, ok
}

func (h *hub) createPeer(userId string) *peer {
	p := &peer{
		userId:  userId,
		ionPeer: isfu.NewPeer(h.sfu.ionSfu),
	}
	return p
}

func (r *room) addPeer(p *peer) {
	r.Lock()
	defer r.Unlock()
	r.peers[p.userId] = p
}

func (p *peer) join(roomId, SDPOffer string) (*webrtc.SessionDescription, error) {
	return p.ionPeer.Join(roomId, webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: SDPOffer})
}

func (p *peer) answer(roomId, SDPOffer string) (*webrtc.SessionDescription, error) {
	return p.ionPeer.Answer(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: SDPOffer})
}

//func (r *room) toDomain() *domain.Room {
//
//	r.RLock()
//	defer r.Unlock()
//
//	dr := &domain.Room{
//		Id:              r.id,
//		InitiatorUserId: r.initiatorUserId,
//		ChatChannelId:   r.chatChannelId,
//		UserIds:         []string{},
//	}
//
//	for _, r := range r.peers {
//		dr.UserIds = append(dr.UserIds, r.userId)
//	}
//
//	return dr
//
//}

func NewRoomHub(cfg *kitConfig.Config, storage domain.RoomStorage) domain.RoomHub {
	return &hub{
		rooms:   make(map[string]*room),
		storage: storage,
		sfu:     newSfu(cfg),
		cfg:     cfg,
	}
}

func (h *hub) createRoom(ctx context.Context, rq *domain.WebRtsSignalCreateRoomRequest) (*domain.WebRtcSignalCreateRoomResponse, error) {

	// create a room
	r := &room{
		id:              kit.NewId(),
		initiatorUserId: rq.UserId,
		chatChannelId:   rq.ChannelId,
		peers:           map[string]*peer{},
	}
	h.addRoom(r)

	// create and join peer
	p := h.createPeer(rq.UserId)
	sdpAnswer, err := p.join(r.id, rq.SDPOffer)
	if err != nil {
		return nil, err
	}

	// add peer to the room
	r.addPeer(p)

	// create domain response
	rs := &domain.WebRtcSignalCreateRoomResponse{
		RoomId:    r.id,
		SDPAnswer: sdpAnswer.SDP,
	}

	// save to storage
	//h.storage.SaveAsync(ctx, rs.Room)

	log.DbgF("[webrtc] room created %s", r.id)

	return rs, nil

}

func (h *hub) Join(ctx context.Context, rq *domain.JoinRoomRequest) (*domain.JoinRoomResponse, error) {

	rs := &domain.JoinRoomResponse{}

	r, ok := h.getRoom(rq.RoomId)
	if !ok {
		// for cluster mode a room is hosted on one node, so it's OK not find it on other nodes
		log.InfF("[webrtc] no room found %s", rq.RoomId)
		return rs, nil
	}

	p := h.createPeer(rq.UserId)
	sdpAnswer, err := p.answer(rq.RoomId, rq.SDPOffer)
	if err != nil {
		return nil, err
	}

	// add peer to the room
	r.addPeer(p)

	rs.SDPAnswer = sdpAnswer.SDP
	//rs.Room = r.toDomain()

	return rs, nil

}

func (h *hub) offer(ctx context.Context, msg map[string]interface{}) error {

	return nil
}

func (h *hub) answer(ctx context.Context, msg map[string]interface{}) error {
	return nil
}

func (h *hub) signal(ctx context.Context, signalMsg *domain.WebRtcSignal) (interface{}, error) {

	//var err error
	//
	//switch signalMsg.Method {
	//case domain.SIGNAL_METHOD_FROM_CLIENT_ROOM_RQ:
	//	var rq *domain.WebRtsSignalCreateRoomRequest
	//	if err := mapstructure.Decode(signalMsg.Payload, &rq); err != nil {
	//		return nil, err
	//	}
	//	return h.createRoom(ctx, rq)
	//
	//case domain.SIGNAL_METHOD_FROM_CLIENT_OFFER:
	//	var rq *domain.WebRtsSignalCreateRoomRequest
	//	if err := mapstructure.Decode(signalMsg.Payload, &rq); err != nil {
	//		return nil, err
	//	}
	//	err = h.offer(ctx, payload)
	//default:
	//	err = fmt.Errorf("[webrtc] unsupported method")
	//}
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil

}

func (h *hub) WebrtcWsMessageHandler(msg []byte) error {

	var wsMessage *proto.WsMessage

	ctx, err := queue.Decode(context.Background(), msg, &wsMessage)
	if err != nil {
		return err
	}

	if wsMessage.MessageType != proto.WS_MESSAGE_TYPE_WEBRTC {
		return fmt.Errorf("[webrtc] ws message incorrect routing")
	}

	var signalMsg *domain.WebRtcSignal
	if err := mapstructure.Decode(wsMessage.Data, &signalMsg); err != nil {
		return err
	}

	h.signal(ctx, signalMsg)

	return nil

}
