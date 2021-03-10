package zeebe

import (
	"context"
	"fmt"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
)

type Utils struct {
	logger log.CLoggerFunc
}

func NewUtils(logger log.CLoggerFunc) *Utils {
	return &Utils{logger: logger}
}

func (u *Utils) l() log.CLogger {
	return u.logger().Cmp("zeebe")
}

func (u *Utils) FailJob(client worker.JobClient, job entities.Job, err error) {
	u.l().Mth("fail-job").F(log.FF{"job": job.GetKey()}).E(err).St().Err()
	_, _ = client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).ErrorMessage(err.Error()).Send(context.Background())
}

func (u *Utils) CompleteJob(client worker.JobClient, job entities.Job, vars map[string]interface{}) error {

	cmd := client.NewCompleteJobCommand().JobKey(job.GetKey())

	l := u.l().Mth("complete-job").F(log.FF{"job": job.GetKey()})

	if vars != nil && len(vars) > 0 {
		rq, err := cmd.VariablesFromMap(vars)
		if err != nil {
			return err
		}
		_, err = rq.Send(context.Background())
		if err != nil {
			return err
		}
		l.F(log.FF{"vars": vars}).Trc("ok")
	} else {
		_, err := cmd.Send(context.Background())
		if err != nil {
			return err
		}
		l.Trc("ok")
	}

	return nil

}

func (u *Utils) GetVarsAndCtx(job entities.Job) (map[string]interface{}, context.Context, error){
	variables, err := job.GetVariablesAsMap()
	if err != nil {
		return nil, nil, err
	}
	ctx, err := u.CtxFromVars(variables)
	return variables, ctx, err
}

func (u *Utils) CtxFromVars(vars map[string]interface{}) (context.Context, error) {
	if mp, ok := vars["_ctx"].(map[string]interface{}); ok {
		ctx, err := kitContext.FromMap(context.Background(), mp)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	}
	return nil, fmt.Errorf("variable _ctx not found or invalid")
}

func (u *Utils) CtxToVars(ctx context.Context, vars map[string]interface{}) error {
	if r, ok := kitContext.Request(ctx); ok {
		vars["_ctx"] = r
		return nil
	}
	return fmt.Errorf("request context not found")
}


