package bp_engine

import (
	"fmt"
	domain "gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/bp/logger"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/proto/config"
	"os"
	"path/filepath"
)

type Adapter interface {
	Init(c *config.Config, bps []domain.BusinessProcess, queueListener listener.QueueListener) error
	GetEngine() bpm.Engine
	Close()
}

type adapterImpl struct {
	Bpm bpm.Engine
}

func NewAdapter() Adapter {
	c := &adapterImpl{
		Bpm: zeebe.NewEngine(logger.LF()),
	}
	return c
}

func (c *adapterImpl) Init(cfg *config.Config, bps []domain.BusinessProcess, queueListener listener.QueueListener) error {

	l := logger.L().Cmp("bp-engine").Mth("init")

	err := c.Bpm.Open(&bpm.Params{
		Port: cfg.Zeebe.Port,
		Host: cfg.Zeebe.Host,
	})
	if err != nil {
		return err
	}

	if _, err := os.Stat(cfg.Bpmn.SrcFolder); err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("bpmn folder not found %s", cfg.Bpmn.SrcFolder)
			l.E(err).St().Err()
		}
		return err
	}

	f, _ := filepath.Abs(cfg.Bpmn.SrcFolder)
	l.DbgF("bpmn loading from folder %s", f)

	var BPMNs []string
	for _, bp := range bps {
		if err := bp.Init(); err != nil {
			return err
		}

		path := filepath.Join(cfg.Bpmn.SrcFolder, bp.GetBPMNFileName())

		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				err := fmt.Errorf("bpmn file not found %s", path)
				l.E(err).St().Warn()
			}
			continue
		}

		BPMNs = append(BPMNs, path)
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
