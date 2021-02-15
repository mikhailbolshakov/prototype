package api

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func Test_NewClientLogin_Success(t *testing.T) {

	helper := NewTestHelper()

	_, _, err := helper.Login(TEST_USER)
	if err != nil {
		t.Fatal(err)
	}

	user, err := helper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	sessionId, _, err := helper.Login(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("test passed", sessionId)

}

func Test_NewClientLogout_Success(t *testing.T) {

	helper := NewTestHelper()

	_, _,err := helper.Login(TEST_USER)
	if err != nil {
		t.Fatal(err)
	}

	user, err := helper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	sessionId, _, err := helper.Login(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	err = helper.Logout(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	_, err = helper.GetUser(user.Username)
	if err != nil {

		if strings.Contains(err.Error(), "missing") {
			fmt.Println("test passed", sessionId)
			return
     	} else {
     		t.Fatal(err)
		}
	}
	t.Fatal("error expected")

}

func Test_LogoutWithoutLogin_Success(t *testing.T) {

	helper := NewTestHelper()

	err := helper.Logout(TEST_USER)
	if err != nil {

		if strings.Contains(err.Error(), "missing") {
			fmt.Println("test passed")
			return
		} else {
			t.Fatal(err)
		}
	}
	t.Fatal("error expected")

}

func Test_ReconnectWithSameSession(t *testing.T) {

	helper := NewTestHelper()

	sessionId, closeCh, err := helper.Login(TEST_USER)
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second * 10)
	// close WS connection
	closeCh <- struct{}{}

	time.Sleep(time.Second)

	_, _, err = helper.Ws(sessionId)
	if err != nil {
		t.Fatal(err)
	}

	<-time.After(time.Second * 10)

}

func Test_MultipleConnectionsSameUser(t *testing.T) {

	helper := NewTestHelper()

	_, _, err := helper.Login(TEST_USER)
	if err != nil {
		t.Fatal(err)
	}

	user, err := helper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = helper.Login(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	m, err := helper.MonitorUserSessions(user.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.Sessions) != 1 {
		t.Fatal("expected one session")
	}

	_, _, err = helper.Login(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	m, err = helper.MonitorUserSessions(user.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.Sessions) != 2 {
		t.Fatal("expected two sessions")
	}

	<-time.After(time.Second * 10)

	err = helper.Logout(user.Username)
	if err != nil {
		t.Fatal(err)
	}

	m, err = helper.MonitorUserSessions(user.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.Sessions) != 0 {
		t.Fatal("expected no sessions")
	}

}