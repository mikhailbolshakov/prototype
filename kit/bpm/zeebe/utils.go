package zeebe

import (
	"context"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"log"
)

func FailJob(client worker.JobClient, job entities.Job, err error) {

	log.Printf("failed to complete job %s error %v", job.GetKey(), err)

	_, _ = client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).ErrorMessage(err.Error()).Send(context.Background())

}

func CompleteJob(client worker.JobClient, job entities.Job, vars map[string]interface{}) error {

	cmd := client.NewCompleteJobCommand().JobKey(job.GetKey())

	if vars != nil && len(vars) > 0 {
		rq, err := cmd.VariablesFromMap(vars)
		if err != nil {
			return err
		}
		_, err = rq.Send(context.Background())
		if err != nil {
			return err
		}
		return nil

	} else {

		_, err := cmd.Send(context.Background())
		if err != nil {
			return err
		}

		return nil
	}

}


