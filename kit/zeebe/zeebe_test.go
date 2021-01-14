package zeebe

import (
	"context"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	"log"
	"testing"
	"github.com/zeebe-io/zeebe/clients/go/pkg/pb"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

const zeebeAddr = "localhost:26500"

func Test1(t *testing.T) {

	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         zeebeAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	topology, err := zbClient.NewTopologyCommand().Send(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, broker := range topology.Brokers {
		log.Println("Broker", broker.Host, ":", broker.Port)
		for _, partition := range broker.Partitions {
			log.Println("  Partition", partition.PartitionId, ":", roleToString(partition.Role))
		}
	}

}

func roleToString(role pb.Partition_PartitionBrokerRole) string {
	switch role {
	case pb.Partition_LEADER:
		return "Leader"
	case pb.Partition_FOLLOWER:
		return "Follower"
	default:
		return "Unknown"
	}
}

func Test2(t *testing.T) {

	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         zeebeAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	response, err := zbClient.NewDeployWorkflowCommand().AddResourceFile("order-process.bpmn").Send(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(response.String())
}

func Test3(t *testing.T) {

	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         zeebeAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	// After the workflow is deployed.
	variables := make(map[string]interface{})
	variables["orderId"] = "31243"

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId("order-process-2").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	msg, err := request.Send(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(msg.String())
}


var readyClose = make(chan struct{})

func Test4(t *testing.T) {

	zbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         zeebeAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	// deploy workflow
	ctx := context.Background()
	response, err := zbClient.NewDeployWorkflowCommand().AddResourceFile("order-process-4.bpmn").Send(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(response.String())

	// create a new workflow instance
	variables := make(map[string]interface{})
	variables["orderId"] = "31243"

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId("order-process-4").LatestVersion().VariablesFromMap(variables)
	if err != nil {
		t.Fatal(err)
	}

	result, err := request.Send(ctx)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(result.String())

	jobWorker := zbClient.NewJobWorker().JobType("payment-service").Handler(handleJob).Open()

	<-readyClose
	jobWorker.Close()
	jobWorker.AwaitClose()
}

func handleJob(client worker.JobClient, job entities.Job) {

	jobKey := job.GetKey()

	headers, err := job.GetCustomHeadersAsMap()
	if err != nil {
		// failed to handle job as we require the custom job headers
		failJob(client, job)
		return
	}

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	variables["totalPrice"] = 46.50
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)
	log.Println("Processing order:", variables["orderId"])
	log.Println("Collect money using payment method:", headers["method"])

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
	close(readyClose)
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
	if err != nil {
		panic(err)
	}
}