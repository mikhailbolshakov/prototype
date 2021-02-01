package mattermost

import (
	"github.com/adacta-ru/mattermost-server/v6/model"
	"log"
	"testing"
)

func Test_CreateBot(t *testing.T) {

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

	bot, rs := admin.RestApi.CreateBot(&model.Bot{Username: "test1.bot", DisplayName: "test", Description: "test"})
	if err := handleResponse(rs); err != nil {
		t.Fatal(err)
	}
	log.Printf("bot created %s\n", bot.UserId)

	uat, rs := admin.RestApi.CreateUserAccessToken(bot.UserId, "bot token")
	if err := handleResponse(rs); err != nil {
		t.Fatal(err)
	}
	log.Printf("token: %s\n", uat.Token)


}

func Test_GetBotToken(t *testing.T) {
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

	bots, rs := admin.RestApi.GetBots(0, 100, "")
	if err := handleResponse(rs); err != nil {
		t.Fatal(err)
	}

	var testBot *model.Bot
	for _, b := range bots {
		if b.Username == "test.bot" {
			testBot = b
			break
		}
	}
	if testBot == nil {
		t.Fatal("bot not found")
	}

	tokens, rs := admin.RestApi.GetUserAccessTokensForUser(testBot.UserId, 0, 100)
	if err := handleResponse(rs); err != nil {
		t.Fatal(err)
	}

	var botToken string
	for _, t := range tokens {
		if t.IsActive {
			botToken = t.Token
			log.Printf("bot token found: %s\n", botToken)
			break
		}
	}

}

func Test_DeleteBot(t *testing.T) {
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

	bots, rs := admin.RestApi.GetBots(0, 100, "")
	if err := handleResponse(rs); err != nil {
		t.Fatal(err)
	}

	var testBot *model.Bot
	for _, b := range bots {
		if b.Username == "test.bot" {
			testBot = b
			break
		}
	}
	if testBot == nil {
		t.Fatal("bot not found")
	}

	_, rs = admin.RestApi.DisableBot(testBot.UserId)
	if err := handleResponse(rs); err != nil {
		t.Fatal(err)
	}

	_, rs = admin.RestApi.DeleteUser(testBot.UserId)
	if err := handleResponse(rs); err != nil {
		t.Fatal(err)
	}

}

func Test_BotLoginWithAccessToken(t *testing.T) {

	token := "zgpbxezgw3gktn54aw7o3enfgc"
	a := model.NewAPIv4Client("http://localhost:8065")
	a.SetOAuthToken(token)
	bot, rs := a.GetMe("")
	if err := handleResponse(rs); err != nil {
		t.Fatal(err)
	}
	log.Printf("bot username %s\n", bot.Username)

}


