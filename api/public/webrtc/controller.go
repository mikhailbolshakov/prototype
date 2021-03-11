package webrtc

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gitlab.medzdrav.ru/prototype/api/public"
	kitHttp "gitlab.medzdrav.ru/prototype/kit/http"
	"net/http"
)

type Controller interface {
	CreateRoom(http.ResponseWriter, *http.Request)
	GetRoom(http.ResponseWriter, *http.Request)
}

type ctrlImpl struct {
	kitHttp.Controller
	roomService public.RoomService
}

func NewController(roomService public.RoomService) Controller {

	return &ctrlImpl{
		roomService:     roomService,
	}

}

func (c *ctrlImpl) CreateRoom(w http.ResponseWriter, r *http.Request) {

	rq := &CreateRoomRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(rq); err != nil {
		c.RespondError(w, http.StatusBadRequest, errors.New("invalid request"))
		return
	}

	room, err := c.roomService.Create(r.Context(), rq.ChannelId)
	if err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	c.RespondOK(w, c.toRoomApi(room))
}

func (c *ctrlImpl) GetRoom(w http.ResponseWriter, r *http.Request) {

	roomId := mux.Vars(r)["roomId"]

	if roomId == "" {
		c.RespondError(w, http.StatusBadRequest, errors.New("no roomId specified"))
		return
	}

	room, err := c.roomService.Get(r.Context(), roomId)
	if err != nil {
		c.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	if room != nil {
		c.RespondOK(w, c.toRoomApi(room))
	} else {
		c.RespondOK(w, struct{}{})
	}
}
