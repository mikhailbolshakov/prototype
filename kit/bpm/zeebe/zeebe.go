package zeebe

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"path/filepath"
	"strings"
	"time"
)

const ERR_EXHAUSTED_RESOURCES_MAX_RETRY = 10

type engineImpl struct {
	params *bpm.Params
	client zbc.Client
	jobWorkers []worker.JobWorker
}

func NewEngine() bpm.Engine {

	zeebe := &engineImpl{
		params: &bpm.Params{},
		jobWorkers: []worker.JobWorker{},
	}

	return zeebe

}

func (z *engineImpl) Open(params *bpm.Params) error {

	z.params = params

	// if already opened, close it
	if err := z.Close(); err != nil {
		return err
	}

	if z.params.Port == "" || z.params.Host == "" {
		return errors.New("[zeebe] cannot open connection, params are invalid")
	}

	if zc, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         fmt.Sprintf("%s:%s", z.params.Host, z.params.Port),
		UsePlaintextConnection: true,
	}); err == nil {
		z.client = zc
		log.DbgF("[zeebe] connection is opened on %s:%s", z.params.Host, z.params.Port)
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

	for _, p := range paths {

		absPath, _ := filepath.Abs(p)

		go func(path string){

			errRetryCount := 0

			for {
				rs, err := z.client.NewDeployWorkflowCommand().AddResourceFile(path).Send(context.Background())
				if err != nil {
					if strings.Contains(err.Error(), "ResourceExhausted") && errRetryCount <= ERR_EXHAUSTED_RESOURCES_MAX_RETRY {
						time.Sleep(time.Millisecond * 100)
						errRetryCount++
					} else {
						log.Err(fmt.Errorf("[zeebe] deployment %s. err: %v", path, rs), true)
						return
					}
				} else {
					log.DbgF("[zeebe] %s deployed. details: %v", path, rs)
					return
				}
			}
		}(absPath)

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

	log.DbgF("[zeebe] process %s started. details = %v", processId, pRs.String())

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
	log.DbgF("[zeebe] message %s published, response: %v", messageId, rs)

	return nil
}

func (z *engineImpl) SendError(jobId int64, errCode, errMessage string) error {

	m := z.client.NewThrowErrorCommand().JobKey(jobId).ErrorCode(errCode).ErrorMessage(errMessage)

	rs, err := m.Send(context.Background())
	if err != nil {
		return err
	}
	log.DbgF("[zeebe] error sent. code: %s, message: %s, response: %v", errCode, errMessage, rs)

	return nil
}
