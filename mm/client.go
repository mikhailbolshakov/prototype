package mm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adacta-ru/mattermost-server/v6/model"
	"log"
	"net/http"
)

type Params struct {
	Url string
	WsUrl string
	Account string
	Pwd string
}

type Client struct {
	RestApi *model.Client4
	WsApi *model.WebSocketClient
	Token string
	User *model.User
	Quit chan interface{}
	Params *Params
}

func handleResponse(rs *model.Response) error {

	success := map[int]bool{
		http.StatusOK:       true,
		http.StatusCreated:  true,
		http.StatusAccepted: true,
	}

	if ok := success[rs.StatusCode]; !ok {
		return errors.New(fmt.Sprintf("Status: %d, %v", rs.StatusCode, rs.Error))
	}
	return nil
}

func Login(p *Params) (*Client, error) {

	cl := &Client{
		Quit: make(chan interface{}),
	}

	cl.RestApi = model.NewAPIv4Client(p.Url)
	user, resp := cl.RestApi.Login(p.Account, p.Pwd)
	if err := handleResponse(resp); err != nil {
		return nil, err
	}
	cl.User = user
	cl.Token = resp.Header.Get("Token")

	log.Printf("muttermost connected. user: %s", cl.User.Email)

	appErr := &model.AppError{}
	cl.WsApi, appErr = model.NewWebSocketClient(p.WsUrl, cl.Token)
	if appErr != nil {
		return nil, errors.New(appErr.Message)
	}

	go cl.WsApi.Listen()
	go func() {
		for {
			select {

			case event := <-cl.WsApi.EventChannel:
				s, _ := json.MarshalIndent(event, "", "\t")
				log.Printf("[WS event]. %s", s)
			case response := <-cl.WsApi.ResponseChannel:
				s, _ := json.MarshalIndent(response, "", "\t")
				log.Printf("[WS response]. %s", s)
			case <-cl.Quit:
				log.Printf("[WS closing]. Closing request for user %s", cl.User.Email)
				cl.WsApi.Close()
				return
			}
		}
	}()

	return cl, nil
}

func (c *Client) Logout() error {
	_, rs := c.RestApi.Logout()
	if err := handleResponse(rs); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateUser(rq *CreateUserRequest) (*CreateUserResponse, error) {

	user, rs := c.RestApi.CreateUser(&model.User{
		Username: rq.Username,
		Password: rq.Password,
		Email:    rq.Email,
	})
	if err := handleResponse(rs); err != nil {
		return nil, err
	}
	log.Printf("User created. Id: %s, email: %s", user.Id, user.Email)

	// add to team
	team, rs := c.RestApi.GetTeamByName(rq.TeamName, "")
	if err := handleResponse(rs); err != nil {
		return nil, err
	}

	_, rs = c.RestApi.AddTeamMember(team.Id, user.Id)
	if err := handleResponse(rs); err != nil {
		return nil, err
	}

	response := &CreateUserResponse{
		Id: user.Id,
	}

	return response, nil
}

