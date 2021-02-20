package bp_engine

import (
	domain "gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	kitConfig "gitlab.medzdrav.ru/prototype/kit/config"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
)

type Adapter interface {
	Init(c *kitConfig.Config, bps []domain.BusinessProcess, queueListener listener.QueueListener) error
	GetEngine() bpm.Engine
	Close()
}

type adapterImpl struct {
	Bpm bpm.Engine
}

func NewAdapter() Adapter {
	c := &adapterImpl{
		Bpm: zeebe.NewEngine(),
	}
	return c
}

func (c *adapterImpl) Init(cfg *kitConfig.Config, bps []domain.BusinessProcess, queueListener listener.QueueListener) error {

	l := log.L().Cmp("bp-engine").Mth("init")

	err := c.Bpm.Open(&bpm.Params{
		Port: cfg.Zeebe.Port,
		Host: cfg.Zeebe.Host,
	})
	if err != nil {
		return err
	}

	var BPMNs []string
	for _, bp := range bps {
		if err := bp.Init(); err != nil {
			return err
		}
		BPMNs = append(BPMNs, bp.GetBPMNPath())
		bp.SetQueueListeners(queueListener)
		l.F(log.FF{"bpmn": bp.GetId()}).Dbg("ok")
	}

	if len(BPMNs) > 0 {
		if err := c.Bpm.DeployBPMNs(BPMNs); err != nil {
			l.E(err).Err("deploy failed")
		}
	}

	return nil

}

func (c *adapterImpl) GetEngine() bpm.Engine {
	return c.Bpm
}

func (c *adapterImpl) Close() {
	_ = c.Bpm.Close()
}
