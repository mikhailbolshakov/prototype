package dentist_online_consultation

import (
	"context"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
	"log"
	"testing"
	"time"
)

const zeebeAddr = "localhost:26500"

func Test_Deploy(t *testing.T) {

	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         zeebeAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	response, err := zbClient.NewDeployWorkflowCommand().AddResourceFile("../bpmn/expert_online_consultation.bpmn").Send(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(response.String())

}

func Test_RunProcess_ConsultationFinished(t *testing.T) {

	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         zeebeAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	deliveryId := "123"

	variables := make(map[string]interface{})
	variables["deliveryId"] = deliveryId

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId("p-expert-online-consultation").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	msg, err := request.Send(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(msg.String())

	_ = zbClient.NewJobWorker().JobType("st-create-task").Handler(createTaskJobHandler).Open()
	//defer createTaskWorker.Close()
	//defer createTaskWorker.AwaitClose()

	_ = zbClient.NewJobWorker().JobType("st-task-in-progress").Handler(taskInProgressJobHandler).Open()
	//defer taskInProgressWorker.Close()
	//defer taskInProgressWorker.AwaitClose()

	_ = zbClient.NewJobWorker().JobType("ts-complete-consultation").Handler(completeConsultationJobHandler).Open()
	//defer taskCompleteConsultationWorker.Close()
	//defer taskCompleteConsultationWorker.AwaitClose()

	time.Sleep(time.Second * 5)

	rs, err := zbClient.NewPublishMessageCommand().MessageName("msg-consultation-time").CorrelationKey(deliveryId).Send(ctx)
	log.Printf("publish message response: %v", rs)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	cmd, err := zbClient.NewPublishMessageCommand().MessageName("msg-task-finished").CorrelationKey(deliveryId).VariablesFromMap(map[string]interface{}{"taskCompleted": true})
	if err != nil {
		t.Fatal(err)
	}
	rs, err = cmd.Send(ctx)
	log.Printf("publish message response: %v", rs)
	if err != nil {
		t.Fatal(err)
	}

	select {
		case <-processFinished:
		case <-time.After(time.Second * 30): log.Println("timeout")
	}

}

var processFinished = make(chan struct{})

func createTaskJobHandler(client worker.JobClient, job entities.Job) {

	jobKey := job.GetKey()

	log.Println("task created")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	variables["taskCompleted"] = false
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")

}

func taskInProgressJobHandler(client worker.JobClient, job entities.Job) {

	jobKey := job.GetKey()

	log.Println("task in progress")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	variables["taskCompleted"] = false
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")

}

func completeConsultationJobHandler(client worker.JobClient, job entities.Job) {

	jobKey := job.GetKey()

	log.Println("consultation is completed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")

	close(processFinished)

}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())
}