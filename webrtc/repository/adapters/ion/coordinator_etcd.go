package ion

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"sync"
	"time"

	"github.com/pborman/uuid"
	"github.com/pion/ion-sfu/pkg/sfu"
	"google.golang.org/grpc"

)

type etcdCoordinator struct {
	mu           sync.Mutex
	nodeID       string
	nodeEndpoint string
	client       *clientv3.Client
	w          sfu.WebRTCTransportConfig
	localRooms map[string]*sfu.Session
	roomLeases map[string]context.CancelFunc
}

func newCoordinatorEtcd(ctx context.Context, cfg *config.Webrtc, etcdCfg *config.Etcd) (*etcdCoordinator, error) {

	l := log.L().Cmp("ion").Mth("etcd-coord.new")

	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: time.Second * 3,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
		Endpoints:   etcdCfg.Hosts,
	})
	if err != nil {
		return nil, err
	}
	w := sfu.NewWebRTCTransportConfig(*cfg.SFU)

	e := endpoint(cfg)
	l.F(log.FF{"endpoint": e}).DbgF("created")
	return &etcdCoordinator{
		client:       cli,
		nodeID:       uuid.New(),
		nodeEndpoint: e,
		w:            w,
		roomLeases:   make(map[string]context.CancelFunc),
		localRooms:   make(map[string]*sfu.Session),
	}, nil
}

func (e *etcdCoordinator) getOrCreateRoom(ctx context.Context, roomID string) (*roomMeta, error) {

	l := log.L().Cmp("ion").Mth("etcd-coord.get-session").C(ctx).F(log.FF{"room": roomID})

	// This operation is only allowed 5 seconds to complete
	etcdCtx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	s, err := concurrency.NewSession(e.client, concurrency.WithContext(etcdCtx))
	if err != nil {
		return nil, err
	}

	// Acquire the lock for this room
	key := fmt.Sprintf("/webrtc/%v", roomID)
	mu := concurrency.NewMutex(s, key)
	if err := mu.Lock(ctx); err != nil {
		l.E(err).Err("could not acquire session lock")
		return nil, err
	}
	defer mu.Unlock(ctx)

	// Check to see if roomMeta exists for this key
	gr, err := e.client.Get(ctx, key)
	if err != nil {
		l.E(err).Err("error looking up session")
		return nil, err
	}

	// Session already exists somewhere in the cluster
	// return the existing meta to the user
	if gr.Count > 0 {
		l.Dbg("found")

		var meta roomMeta
		if err := json.Unmarshal(gr.Kvs[0].Value, &meta); err != nil {
			l.E(err).St().Err("unmarshal")
			return nil, err
		}
		meta.Redirect = (meta.NodeID != e.nodeID)

		// return meta for session
		return &meta, nil
	}

	// Room does not already exist, so lets take it
	// @todo load balance here / be smarter

	// First lets create a lease for the sessionKey
	lease, err := e.client.Grant(ctx, 1)
	if err != nil {
		l.E(err).Err("acquiring lease")
		return nil, err
	}
	l.F(log.FF{"leaseId": lease.ID})

	etcdCtx, leaseCancel := context.WithCancel(context.Background())
	leaseKeepAlive, err := e.client.KeepAlive(etcdCtx, lease.ID)
	if err != nil {
		l.E(err).St().Err("activating keepAlive")
	}
	<-leaseKeepAlive

	e.mu.Lock()
	e.roomLeases[roomID] = leaseCancel
	defer e.mu.Unlock()

	meta := roomMeta{
		RoomID:       roomID,
		NodeID:       e.nodeID,
		NodeEndpoint: e.nodeEndpoint,
	}
	payload, _ := json.Marshal(&meta)
	_, err = e.client.Put(ctx, key, string(payload), clientv3.WithLease(lease.ID))
	if err != nil {
		l.E(err).St().Err("storing room meta")
		return nil, err
	}

	return &meta, nil
}

func (e *etcdCoordinator) ensureRoom(roomID string) *sfu.Session {

	e.mu.Lock()
	defer e.mu.Unlock()

	if s, ok := e.localRooms[roomID]; ok {
		return s
	}

	s := sfu.NewSession(roomID)
	s.OnClose(func() {
		e.onRoomClosed(roomID)
	})
	prometheusGaugeRooms.Inc()

	e.localRooms[roomID] = s
	return s
}

func (e *etcdCoordinator) GetSession(roomId string) (*sfu.Session, sfu.WebRTCTransportConfig) {
	return e.ensureRoom(roomId), e.w
}

func (e *etcdCoordinator) onRoomClosed(roomID string) {

	l := log.L().Cmp("ion").Mth("etcd-coord.room-closed").F(log.FF{"room": roomID})

	e.mu.Lock()
	defer e.mu.Unlock()

	// Acquire the lock for this roomID
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	s, err := concurrency.NewSession(e.client, concurrency.WithContext(ctx))
	if err != nil {
		l.E(err).St().Err("couldn't acquire room")
		return
	}
	key := fmt.Sprintf("/webrtc/%v", roomID)
	mu := concurrency.NewMutex(s, key)
	if err := mu.Lock(ctx); err != nil {
		l.E(err).St().Err("couldn't acquire room lock")
		return
	}
	defer mu.Unlock(ctx)

	// Cancel our lease
	leaseCancel := e.roomLeases[roomID]
	delete(e.roomLeases, roomID)
	leaseCancel()

	// Delete session meta
	_, err = e.client.Delete(ctx, key)
	if err != nil {
		l.E(err).St().Err("deleting roomMeta")
		return
	}

	// Delete localRoom
	delete(e.localRooms, roomID)
	prometheusGaugeRooms.Dec()

	l.Dbg("closed")
}
