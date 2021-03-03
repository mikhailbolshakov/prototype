package api

import (
	"log"
	"testing"
	"time"
)

func Test_Webrtc_gRPC_Success(t *testing.T) {

	helper := NewTestHelper()

	_, _, err := helper.Login(TEST_USER)
	if err != nil {
		t.Fatal(err)
	}

	user1, err := helper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	user2, err := helper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	_ = helper.Logout(TEST_USER)

	_, _, err = helper.Login(user1.Username)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = helper.Login(user2.Username)
	if err != nil {
		t.Fatal(err)
	}
	_ = helper.NewPeer("123", user2.Username)
	time.Sleep(time.Second)
	_ = helper.NewPeer("123", user1.Username)


	select{}
}

func Test_Webrtc_jsonrpc_Success(t *testing.T) {

	helper := NewTestHelper()

	_, _, err := helper.Login(TEST_USER)
	if err != nil {
		t.Fatal(err)
	}

	user1, err := helper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	user2, err := helper.CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	_ = helper.Logout(TEST_USER)

	s1, _, err := helper.Login(user1.Username)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("sid-1: %s", s1)

	s2, _, err := helper.Login(user2.Username)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("sid-2: %s", s2)

	select{}
}
