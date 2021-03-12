package service

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"go.uber.org/atomic"
)

type Service interface {
	GetCode() string
	Init(ctx context.Context) error
	ListenAsync(ctx context.Context) error
	Close(ctx context.Context)
}

type MetaInfo interface {
	ServiceCode() string
	NodeId() string
	Leader() bool
	SetMeAsLeader(l bool)
	InstanceId() string
}

type metaInfo struct {
	svcCode    string
	nodeId     string
	instanceId string
	leader     *atomic.Bool
}

func NewMetaInfo(svcCode, nodeId string) MetaInfo {
	return &metaInfo{
		svcCode:    svcCode,
		nodeId:     nodeId,
		instanceId: fmt.Sprintf("%s-%s", svcCode, nodeId),
		leader:     atomic.NewBool(true),
	}
}

func (m *metaInfo) ServiceCode() string {
	return m.svcCode
}

func (m *metaInfo) NodeId() string {
	return m.nodeId
}

func (m *metaInfo) InstanceId() string {
	return m.instanceId
}

func (m *metaInfo) Leader() bool {
	return m.leader.Load()
}

func (m *metaInfo) SetMeAsLeader(l bool) {
	m.leader.Store(l)
}

type Cluster struct {
	Raft      Raft
	Meta      MetaInfo
	logger    log.CLoggerFunc
	isCluster bool
}

func NewCluster(logger log.CLoggerFunc, meta MetaInfo) Cluster {
	return Cluster{Raft: NewRaft(logger), Meta: meta, logger: logger}
}

func (c *Cluster) Init(size int, natsUrl string, ev OnLeaderChangedEvent) error {

	l := c.logger().Cmp("cluster").Mth("init")

	if size <= 1 {
		// no cluster needed
		l.Warn("no cluster needed for the given size")
		return nil
	}

	if size % 2 == 0 {
		err := fmt.Errorf("cannot start cluster with odd size")
		l.E(err).Err()
		return err
	}

	err := c.Raft.Init(&Options{
		ClusterName: c.Meta.ServiceCode(),
		ClusterSize: size,
		NatsUrl:     natsUrl,
		LogPath:     "/tmp/raft.log",
	}, func(l bool) {
		c.Meta.SetMeAsLeader(l)
		if ev != nil {
			ev(l)
		}
	})
	if err != nil {
		return err
	}

	c.isCluster = true

	l.Inf("ok")

	return nil

}

func (c *Cluster) Start() error {

	if !c.isCluster {
		return nil
	}

	if err := c.Raft.Start(); err != nil {
		return err
	}
	c.Meta.SetMeAsLeader(c.Raft.AmILeader())
	c.logger().Cmp("cluster").Mth("start").Inf("ok")
	return nil
}

func (c *Cluster) Close() {

	if !c.isCluster {
		return
	}

	c.Raft.Close()
	c.logger().Cmp("cluster").Mth("close").Inf("ok")
}
