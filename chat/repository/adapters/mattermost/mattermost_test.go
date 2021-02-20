package mattermost

import (
	"log"
	"testing"
)

func Test_GetAdminAccessToken(t *testing.T) {

	admin, err := Login(&Params{
		Url:     "http://localhost:8065",
		WsUrl:   "ws://localhost:8065",
		Account: "admin",
		Pwd:     "admin",
		OpenWS:  false,
	})
	if err != nil {
		t.Fatal(err)
	}

	uat, rs := admin.RestApi.CreateUserAccessToken(admin.User.Id, "admin access token")
	if err := HandleResponse(rs); err != nil {
		t.Fatal(err)
	}
	log.Printf("%v", uat)

}

func Test_Ping(t *testing.T) {
	admin, err := Login(&Params{
		Url:     "http://localhost:8065",
		WsUrl:   "ws://localhost:8065",
		Account: "admin",
		Pwd:     "admin",
		OpenWS:  false,
	})
	if err != nil {
		t.Fatal(err)
	}

	r, rs := admin.RestApi.GetPing()
	if err := HandleResponse(rs); err != nil {
		t.Fatal(err)
	}

	log.Printf("%v", r)

}
