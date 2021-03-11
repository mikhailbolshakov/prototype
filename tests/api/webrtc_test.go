package api

import (
	"gitlab.medzdrav.ru/prototype/kit"
	"log"
	"testing"
	"time"
)

func Test_Webrtc_gRPC_StreamToFile_Success(t *testing.T) {

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
	_ = helper.NewPeer("123", user2.Username, true)
	time.Sleep(time.Second)
	_ = helper.NewPeer("123", user1.Username, true)

	select{}
}

func Test_Webrtc_gRPC_WithFakeClients_Success(t *testing.T) {

	helper := NewTestHelper()
	roomId := kit.NewId()
	_ = helper.NewPeer(roomId, "user1", false)
	time.Sleep(time.Second)
	_ = helper.NewPeer(roomId, "user2", false)

	select{}
}

func Test_Webrtc_Login(t *testing.T) {

	helper := NewTestHelper()

	s1, _, err := helper.Login("1614935271891000474")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("sid-1: %s", s1)

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

	//user3, err := helper.CreateClient()
	//if err != nil {
	//	t.Fatal(err)
	//}

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

	//s3, _, err := helper.Login(user3.Username)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//log.Printf("sid-3: %s", s3)

	room, err := helper.CreateRoom("1234")
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%s", kit.MustJson(room))

	select{}
}
