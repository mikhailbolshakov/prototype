package zeebe

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"log"
	"path/filepath"
)

type engineImpl struct {
	params *Params
	client zbc.Client
	jobWorkers []worker.JobWorker
}

type Params struct {
	Port string
	Host string
}

func NewEngine(params *Params) bpm.Engine {

	zeebe := &engineImpl{
		params: params,
		jobWorkers: []worker.JobWorker{},
	}

	return zeebe

}

func (z *engineImpl) Open() error {

	// if already opened, close it
	if err := z.Close(); err != nil {
		return err
	}

	if z.params.Port == "" || z.params.Host == "" {
		return errors.New("cannot open zeebe connection, params are invalid")
	}

	if zc, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         fmt.Sprintf("%s:%s", z.params.Host, z.params.Port),
		UsePlaintextConnection: true,
	}); err == nil {
		z.client = zc
		log.Printf("zeebe connetion is opened on %s:%s", z.params.Host, z.params.Port)
	} else {
		return err
	}

	return nil
}

func (z *engineImpl) IsOpened() bool {
	return z.client != nil
}

func (z *engineImpl) Close() error {
	if z.client != nil {
		err := z.client.Close()
		z.client = nil
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *engineImpl) DeployBPMNs(paths []string) error {

	ctx := context.Background()

	for _, p := range paths {

		absPath, _ := filepath.Abs(p)
		rs, err := z.client.NewDeployWorkflowCommand().AddResourceFile(absPath).Send(ctx)
		if err != nil {
			return err
		}

		log.Printf("zeebe: %s deployed. details: %v", p, rs)

	}

	return nil
}

func (z *engineImpl) RegisterTaskHandlers(handlers map[string]interface{}) error {

	for task, handlerFunc := range handlers {
		if f, ok := handlerFunc.(func(client worker.JobClient, job entities.Job)); ok {
			z.jobWorkers = append(z.jobWorkers, z.client.NewJobWorker().JobType(task).Handler(f).Open())
		} else {
			return errors.New("no valid handler passed")
		}
	}

	return nil

}

func (z *engineImpl) StartProcess(processId string, vars map[string]interface{}) (string, error) {

	p, err := z.client.NewCreateInstanceCommand().BPMNProcessId(processId).LatestVersion().VariablesFromMap(vars)
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	pRs, err := p.Send(ctx)
	if err != nil {
		return "", err
	}

	log.Printf("zeebe has started BPMN process %s. details = %v", processId, pRs.String())

	return fmt.Sprintf("%d", pRs.WorkflowInstanceKey), nil

}

func (z *engineImpl) SendMessage(messageId string, correlationId string, vars map[string]interface{}) error {

	ctx := context.Background()

	m := z.client.NewPublishMessageCommand().MessageName(messageId).CorrelationKey(correlationId)

	if vars != nil && len(vars) > 0 {
		m, _ = m.VariablesFromMap(vars)
	}

	rs, err := m.Send(ctx)
	if err != nil {
		return err
	}
	log.Printf("zeebe publish message, response: %v", rs)

	return nil
}
