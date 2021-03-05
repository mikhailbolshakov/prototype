package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/webrtc/domain"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"time"
)

type etcdRoomCoordImpl struct {
	c *container
}

func (r *etcdRoomCoordImpl) GetOrCreate(ctx context.Context, meta *domain.RoomMeta) (bool, error) {

	roomId := meta.Id
	found := false

	l := log.L().Cmp("webrtc").Mth("etcd-coordination").C(ctx).F(log.FF{"room": roomId}).Dbg().TrcF("meta:%v", *meta)

	// This operation is only allowed 5 seconds to complete
	etcdCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	s, err := concurrency.NewSession(r.c.Etcd.Client, concurrency.WithContext(etcdCtx))
	if err != nil {
		return found, err
	}

	// Acquire the lock for this room
	key := fmt.Sprintf("/webrtc/%v", roomId)
	mu := concurrency.NewMutex(s, key)
	if err := mu.Lock(ctx); err != nil {
		l.E(err).Err("could not acquire session lock")
		return found, err
	}
	defer mu.Unlock(ctx)

	// Check to see if roomMeta exists for this key
	gr, err := r.c.Etcd.Client.Get(ctx, key)
	if err != nil {
		l.E(err).Err("error looking up session")
		return found, err
	}
	l.TrcF("get. key=%s, found=%d", key, gr.Count)

	// Session already exists somewhere in the cluster
	// return the existing meta to the user
	if gr.Count > 0 {
		l.Dbg("found")
		found = true

		if err := json.Unmarshal(gr.Kvs[0].Value, meta); err != nil {
			l.E(err).St().Err("unmarshal")
			return found, err
		}
		// return meta for session
		return found, nil
	}

	// Room does not already exist, so lets take it
	// @todo load balance here / be smarter

	// First lets create a lease for the sessionKey
	lease, err := r.c.Etcd.Client.Grant(ctx, 1)
	if err != nil {
		l.E(err).Err("acquiring lease")
		return found, err
	}
	l.F(log.FF{"leaseId": lease.ID}).TrcF("lease")

	leaseKeepAlive, err := r.c.Etcd.Client.KeepAlive(ctx, lease.ID)
	if err != nil {
		l.E(err).St().Err("activating keepAlive")
		return found, err
	}

	go func() {
		for {
			select { case <-leaseKeepAlive: }
		}
	}()

	payload, _ := json.Marshal(meta)
	_, err = r.c.Etcd.Client.Put(ctx, key, string(payload), clientv3.WithLease(lease.ID))
	if err != nil {
		l.E(err).St().Err("storing room meta")
		return found, err
	}
	l.TrcF("put key=%s", key)

	return found, nil

}

func (r *etcdRoomCoordImpl) Close(ctx context.Context, roomId string) {

	l := log.L().Cmp("webrtc").Mth("etcd-coordination.close-room").F(log.FF{"room": roomId})

	// Acquire the lock for this roomID
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	s, err := concurrency.NewSession(r.c.Etcd.Client, concurrency.WithContext(ctx))
	if err != nil {
		l.E(err).St().Err("couldn't acquire room")
		return
	}
	key := fmt.Sprintf("/webrtc/%v", roomId)
	mu := concurrency.NewMutex(s, key)
	if err := mu.Lock(ctx); err != nil {
		l.E(err).St().Err("couldn't acquire room lock")
		return
	}
	defer mu.Unlock(ctx)

	// Delete session meta
	_, err = r.c.Etcd.Client.Delete(ctx, key)
	if err != nil {
		l.E(err).St().Err("deleting roomMeta")
		return
	}

}
