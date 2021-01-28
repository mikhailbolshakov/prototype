package mattermost

import (
	"errors"
	"fmt"
	"github.com/adacta-ru/mattermost-server/v6/model"
	"log"
	"net/http"
)

type Params struct {
	Url     string
	WsUrl   string
	Account string
	Pwd     string
	OpenWS  bool
}

type Client struct {
	RestApi *model.Client4
	WsApi   *model.WebSocketClient
	Token   string
	User    *model.User
	Quit    chan interface{}
	Params  *Params
}

type UserStatus struct {
	UserId string
	Status string
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

	if p.OpenWS {
		appErr := &model.AppError{}
		cl.WsApi, appErr = model.NewWebSocketClient(p.WsUrl, cl.Token)
		if appErr != nil {
			return nil, errors.New(appErr.Message)
		}
	}

	return cl, nil
}

func (c *Client) Logout() error {
	_, rs := c.RestApi.Logout()
	if err := handleResponse(rs); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateUser(teamName, username, password, email string) (string, error) {

	user, rs := c.RestApi.CreateUser(&model.User{
		Username: username,
		Password: password,
		Email:    email,
	})
	if err := handleResponse(rs); err != nil {
		return "", err
	}

	// add to team
	team, rs := c.RestApi.GetTeamByName(teamName, "")
	if err := handleResponse(rs); err != nil {
		return "", err
	}

	_, rs = c.RestApi.AddTeamMember(team.Id, user.Id)
	if err := handleResponse(rs); err != nil {
		return "", err
	}

	return user.Id, nil
}

func (c *Client) DeleteUser(userId string) error {

	_, rs := c.RestApi.DeleteUser(userId)
	if err := handleResponse(rs); err != nil {
		return err
	}
	log.Printf("user delete. Id: %s", userId)

	return nil

}

func (c *Client) CreateUserChannel(channelType, teamName, userId, displayName, name string) (string, error) {

	team, rs := c.RestApi.GetTeamByName(teamName, "")
	if err := handleResponse(rs); err != nil {
		return "", err
	}

	ch, rs := c.RestApi.CreateChannel(&model.Channel{
		TeamId:      team.Id,
		Type:        channelType,
		DisplayName: displayName,
		Name:        name,
	})
	if err := handleResponse(rs); err != nil {
		return "", err
	}

	_, rs = c.RestApi.AddChannelMember(ch.Id, userId)
	if err := handleResponse(rs); err != nil {
		return "", err
	}

	return ch.Id, nil
}

func (c *Client) SubscribeUser(channelId string, userId string) error {
	_, rs := c.RestApi.AddChannelMember(channelId, userId)
	if err := handleResponse(rs); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateEphemeralPost(channelId string, recipientUserId string, message string, attachments []*model.SlackAttachment) error {

	props := model.StringInterface{}
	if attachments != nil && len(attachments) > 0 {
		props["attachments"] = attachments
	}

	ep := &model.PostEphemeral{
		UserID: recipientUserId,
		Post: &model.Post{
			ChannelId: channelId,
			Message:   message,
			Props:     props,
		},
	}
	_, rs := c.RestApi.CreatePostEphemeral(ep)
	if err := handleResponse(rs); err != nil {
		return err
	}

	return nil

}

func (c *Client) CreatePost(channelId string, message string, attachments []*model.SlackAttachment) error {

	props := model.StringInterface{}
	if attachments != nil && len(attachments) > 0 {
		props["attachments"] = attachments
	}

	p := &model.Post{
		ChannelId: channelId,
		Message:   message,
		Props:     props,
	}

	_, rs := c.RestApi.CreatePost(p)
	if err := handleResponse(rs); err != nil {
		return err
	}

	return nil

}

func (c *Client) GetUsersStatuses(userIds []string) ([]*UserStatus, error) {

	statuses, rs := c.RestApi.GetUsersStatusesByIds(userIds)
	if err := handleResponse(rs); err != nil {
		return nil, err
	}

	var res []*UserStatus
	for _, s := range statuses {
		res = append(res, &UserStatus{
			UserId: s.UserId,
			Status: s.Status,
		})
	}

	return res, nil
}

func (c *Client) CreateDirectChannel(userId1, userId2 string) (string, error) {

	ch, rs := c.RestApi.CreateGroupChannel([]string{userId1, userId2})
	if err := handleResponse(rs); err != nil {
		return "", err
	}

	return ch.Id, nil

}

func (c *Client) GetChannelsForUserAndMembers(userId, teamName string, memberUserIds []string) ([]string, error) {

	team, rs := c.RestApi.GetTeamByName(teamName, "")
	if err := handleResponse(rs); err != nil {
		return nil, err
	}

	channels, rs := c.RestApi.GetChannelsForTeamForUser(team.Id, userId,false, "")
	if err := handleResponse(rs); err != nil {
		return nil, err
	}

	var res []string

	for _, ch := range channels {

		if ch.Type != model.CHANNEL_PRIVATE {
			continue
		}

		chMembers, rs := c.RestApi.GetChannelMembers(ch.Id, 0, 1000, "")
		if err := handleResponse(rs); err != nil {
			return nil, err
		}

		ok := true
		for _, srcM := range memberUserIds {

			found := false
			for _, m := range *chMembers {
				if srcM == m.UserId {
					found = true
					break
				}
			}

			if !found {
				ok = false
				break
			}

		}

		if ok {
			res = append(res, ch.Id)
		}

	}

	return res, nil

}
