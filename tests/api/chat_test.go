package api

import (
	"fmt"
	"testing"
	"time"
)

func Test_ClientPost_TaskCreated_Success(t *testing.T) {

	testHelper := NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	// create a client
	userClient, err := testHelper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	if _, _, err := testHelper.Login(userClient.Username); err != nil {
		t.Fatal(err)
	}

	if err := testHelper.MyPost(userClient.ClientDetails.CommonChannelId, "hahaha"); err != nil {
		t.Fatal(err)
	}

	task, err := testHelper.TaskAwaitByChannel(userClient.ClientDetails.CommonChannelId, time.Second * 20)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("test passed. task=%s status=%s", task.Num, task.Status.SubStatus)

}
