package zeebe

import (
	"context"
	"fmt"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	kitContext "gitlab.medzdrav.ru/prototype/kit/context"
	"gitlab.medzdrav.ru/prototype/kit/log"
)

func FailJob(client worker.JobClient, job entities.Job, err error) {
	log.L().Cmp("zeebe").Mth("fail-job").F(log.FF{"job": job.GetKey()}).E(err).St().Err()
	_, _ = client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).ErrorMessage(err.Error()).Send(context.Background())
}

func CompleteJob(client worker.JobClient, job entities.Job, vars map[string]interface{}) error {

	cmd := client.NewCompleteJobCommand().JobKey(job.GetKey())

	l := log.L().Cmp("zeebe").Mth("complete-job").F(log.FF{"job": job.GetKey()})

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

func GetVarsAndCtx(job entities.Job) (map[string]interface{}, context.Context, error){
	variables, err := job.GetVariablesAsMap()
	if err != nil {
		return nil, nil, err
	}
	ctx, err := CtxFromVars(variables)
	return variables, ctx, err
}

func CtxFromVars(vars map[string]interface{}) (context.Context, error) {
	if mp, ok := vars["_ctx"].(map[string]interface{}); ok {
		ctx, err := kitContext.FromMap(context.Background(), mp)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	}
	return nil, fmt.Errorf("variable _ctx not found or invalid")
}

func CtxToVars(ctx context.Context, vars map[string]interface{}) error {
	if r, ok := kitContext.Request(ctx); ok {
		vars["_ctx"] = r
		return nil
	}
	return fmt.Errorf("request context not found")
}


