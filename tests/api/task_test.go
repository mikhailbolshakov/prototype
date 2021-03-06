package api

import (
	"encoding/json"
	"fmt"
	taskApi "gitlab.medzdrav.ru/prototype/api/public/tasks"
	"strings"
	"testing"
	"time"
)

func Test_CreateTask_Success(t *testing.T) {

	// create user
	testHelper := NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	user, err := testHelper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	if user.ClientDetails == nil || user.ClientDetails.CommonChannelId == "" {
		t.Fatal("user must be a client")
	}

	rq := taskApi.NewTaskRequest{
		Type: &taskApi.Type{
			Type:    "client",
			SubType: "common-request",
		},
		Reported: &taskApi.Reported{
			Username: user.Username,
		},
		Assignee: &taskApi.Assignee{},
		ChannelId: user.ClientDetails.CommonChannelId,
	}

	rqJ, _ := json.Marshal(rq)

	task, err := testHelper.CreateTask(rqJ)
	if err != nil {
		t.Fatal(err)
	}

	if err := testHelper.Logout(TEST_USER); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("test passed. taskId %s", task.Id)

}

func Test_CreateTaskWithEmptyReporter_Error(t *testing.T) {

	// create user
	testHelper := NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	user, err := testHelper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	if user.ClientDetails == nil || user.ClientDetails.CommonChannelId == "" {
		t.Fatal("user must be a client")
	}

	rq := taskApi.NewTaskRequest{
		Type: &taskApi.Type{
			Type:    "client",
			SubType: "common-request",
		},
		Reported: &taskApi.Reported{},
		Assignee: &taskApi.Assignee{},
		ChannelId: user.ClientDetails.CommonChannelId,
	}

	rqJ, _ := json.Marshal(rq)

	_, err = testHelper.CreateTask(rqJ)
	if err != nil {
		if strings.Contains(err.Error(), "cannot find reporter user") {
			fmt.Println("test passed")
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("error expected")
	}

	if err := testHelper.Logout(TEST_USER); err != nil {
		t.Fatal(err)
	}

}

func Test_AutoAssign_Success(t *testing.T) {

	testHelper := NewTestHelper()

	if _, _, err := testHelper.Login(TEST_USER); err != nil {
		t.Fatal(err)
	}

	// create a client
	userClient, err := testHelper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	// create a consultant
	userConsultant, err := testHelper.CreateConsultant("consultant")
	if err != nil {
		t.Fatal(err)
	}

	// set consultant Online
	if err := testHelper.SetStatus(userConsultant.Username, "online"); err != nil {
		t.Fatal(err)
	}

	rq := taskApi.NewTaskRequest{
		Type: &taskApi.Type{
			Type:    "client",
			SubType: "common-request",
		},
		Reported: &taskApi.Reported{
			Username: userClient.Username,
		},
		Assignee: &taskApi.Assignee{},
		ChannelId: userClient.ClientDetails.CommonChannelId,
	}

	rqJ, _ := json.Marshal(rq)

	task, err := testHelper.CreateTask(rqJ)
	if err != nil {
		t.Fatal(err)
	}

	task, err = testHelper.MakeTransition(task.Id, "2")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("task created. taskId %s %s", task.Num, task.Status.SubStatus)

	task, err = testHelper.AssignTaskAwait(task.Id, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if err := testHelper.Logout(TEST_USER); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("test passed. assigned: %s\n", task.Assignee.Username)

}