package http

import (
	"encoding/json"
	"net/http"
)

type Controller struct {}

func (c *Controller) RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {

	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}

func (c *Controller) RespondError(w http.ResponseWriter, code int, err error) {
	c.RespondWithJson(w, code, map[string]string{"error": err.Error()})
}

func (c *Controller) RespondOK(w http.ResponseWriter, payload interface{}) {
	c.RespondWithJson(w, http.StatusOK, payload)
}

