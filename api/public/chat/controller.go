package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.medzdrav.ru/prototype/api/public"
	kitCtx "gitlab.medzdrav.ru/prototype/kit/context"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"net/http"
)

type Controller interface {
	SetStatus(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
	EphemeralPost(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.Controller
	chatService public.ChatService
}

func NewController(chatService public.ChatService) Controller {

	return &ctrlImpl{
		chatService: chatService,
	}
}

// if userId parameter equals "me" value, try to take the current user from context
func (c *ctrlImpl) me(ctx context.Context, userId string) string {

	if userId == "me" {
		rq, _ := kitCtx.Request(ctx)
		return rq.Uid
	}
	return userId
}

func (c *ctrlImpl) SetStatus(writer http.ResponseWriter, r *http.Request) {

	userId := c.me(r.Context(), mux.Vars(r)["id"])
	status := mux.Vars(r)["status"]

	if err := c.chatService.SetStatus(r.Context(), userId, status); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, struct{}{})
	}

}

func (c *ctrlImpl) Post(writer http.ResponseWriter, r *http.Request) {

	userId := c.me(r.Context(), mux.Vars(r)["id"])

	rq := &PostRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	if err := c.chatService.Post(r.Context(), userId, rq.ChannelId, rq.Message); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, struct{}{})
	}

}

func (c *ctrlImpl) EphemeralPost(writer http.ResponseWriter, r *http.Request) {

	userId := c.me(r.Context(), mux.Vars(r)["id"])

	rq := &EphemeralPostRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(writer, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	if err := c.chatService.EphemeralPost(r.Context(), userId, rq.ToUserId, rq.ChannelId, rq.Message); err != nil {
		c.RespondError(writer, http.StatusInternalServerError, err)
	} else {
		c.RespondOK(writer, struct{}{})
	}

}