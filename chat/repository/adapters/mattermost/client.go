package mattermost

import (
	"errors"
	"fmt"
	"github.com/adacta-ru/mattermost-server/v6/model"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"net/http"
)

type Params struct {
	Url         string
	WsUrl       string
	Account     string
	Pwd         string
	AccessToken string
	OpenWS      bool
}

type Client struct {
	RestApi *model.Client4
	WsApi   *model.WebSocketClient
	Token   string
	User    *model.User
	Params  *Params
}

type UserStatus struct {
	UserId string
	Status string
}

func HandleResponse(rs *model.Response) error {

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

	l := log.L().Cmp("mm-client").Mth("login").F(log.FF{"user": p.Account})

	cl := &Client{Params: p}

	cl.RestApi = model.NewAPIv4Client(p.Url)

	var user *model.User
	var rs *model.Response
	if p.AccessToken != "" {
		cl.RestApi.SetOAuthToken(p.AccessToken)
		user, rs = cl.RestApi.GetUserByUsername(p.Account, "")
		if err := HandleResponse(rs); err != nil {
			return nil, err
		}
		cl.Token = p.AccessToken
	} else {
		user, rs = cl.RestApi.Login(p.Account, p.Pwd)
		if err := HandleResponse(rs); err != nil {
			return nil, err
		}
		cl.Token = rs.Header.Get("Token")
	}

	cl.User = user
	l.Trc("logged")

	if p.OpenWS {
		appErr := &model.AppError{}
		cl.WsApi, appErr = model.NewWebSocketClient(p.WsUrl, cl.Token)
		if appErr != nil {
			return nil, errors.New(appErr.Message)
		}
	}

	return cl, nil
}

// Ping
func (c *Client) Ping() bool {
	r, rs := c.RestApi.GetPing()
	if err := HandleResponse(rs); err != nil || r != "OK" {
		return false
	}
	return true
}

// This version of Ping can be called before Login
func ping(url, accessToken string) bool {
	restApi := model.NewAPIv4Client(url)
	restApi.SetOAuthToken(accessToken)
	r, rs := restApi.GetPing()
	if err := HandleResponse(rs); err != nil || r != "OK" {
		return false
	}
	return true
}

func (c *Client) SetStatus(userId, status string) error {

	s, rs := c.RestApi.UpdateUserStatus(userId, &model.Status{
		UserId:         userId,
		Status:         status,
	})
	if err := HandleResponse(rs); err != nil {
		return err
	}

	if s.Status != status {
		return fmt.Errorf("status hasn't been changed")
	}

	return nil
}

func (c *Client) Logout() error {

	l := log.L().Cmp("mm-client").Mth("del-user").F(log.FF{"user": c.User.Username})

	_, rs := c.RestApi.Logout()
	if err := HandleResponse(rs); err != nil {
		return err
	}
	l.Dbg("logged out")

	return nil
}

func (c *Client) CreateUser(teamName, username, password, email string) (string, error) {

	l := log.L().Cmp("mm-client").Mth("create-user").F(log.FF{"user": username})

	user, rs := c.RestApi.CreateUser(&model.User{
		Username: username,
		Password: password,
		Email:    email,
	})
	if err := HandleResponse(rs); err != nil {
		return "", err
	}

	// add to team
	team, rs := c.RestApi.GetTeamByName(teamName, "")
	if err := HandleResponse(rs); err != nil {
		return "", err
	}

	_, rs = c.RestApi.AddTeamMember(team.Id, user.Id)
	if err := HandleResponse(rs); err != nil {
		return "", err
	}

	l.Dbg("created")

	return user.Id, nil
}

func (c *Client) DeleteUser(userId string) error {

	l := log.L().Cmp("mm-client").Mth("del-user").F(log.FF{"user": userId})

	_, rs := c.RestApi.DeleteUser(userId)
	if err := HandleResponse(rs); err != nil {
		return err
	}
	l.Dbg("deleted")

	return nil

}

func (c *Client) CreateUserChannel(channelType, teamName, userId, displayName, name string) (string, error) {

	l := log.L().Cmp("mm-client").Mth("create-channel").F(log.FF{"user": c.User.Username})

	team, rs := c.RestApi.GetTeamByName(teamName, "")
	if err := HandleResponse(rs); err != nil {
		return "", err
	}

	ch, rs := c.RestApi.CreateChannel(&model.Channel{
		TeamId:      team.Id,
		Type:        channelType,
		DisplayName: displayName,
		Name:        name,
	})
	if err := HandleResponse(rs); err != nil {
		return "", err
	}

	_, rs = c.RestApi.AddChannelMember(ch.Id, userId)
	if err := HandleResponse(rs); err != nil {
		return "", err
	}

	l.Dbg("created")

	return ch.Id, nil
}

func (c *Client) SubscribeUser(channelId string, userId string) error {
	_, rs := c.RestApi.AddChannelMember(channelId, userId)
	if err := HandleResponse(rs); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateEphemeralPost(channelId string, recipientUserId string, message string, attachments []*model.SlackAttachment) error {

	l := log.L().Cmp("mm-client").Mth("create-eph-post").F(log.FF{"user": c.User.Username})

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
	if err := HandleResponse(rs); err != nil {
		return err
	}

	l.Dbg("posted")

	return nil

}

func (c *Client) CreatePost(channelId string, message string, attachments []*model.SlackAttachment) error {

	l := log.L().Cmp("mm-client").Mth("create-post").F(log.FF{"user": c.User.Username})

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
	if err := HandleResponse(rs); err != nil {
		return err
	}

	l.Dbg("posted")

	return nil

}

func (c *Client) GetUsersStatuses(userIds []string) ([]*UserStatus, error) {

	statuses, rs := c.RestApi.GetUsersStatusesByIds(userIds)
	if err := HandleResponse(rs); err != nil {
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
	if err := HandleResponse(rs); err != nil {
		return "", err
	}

	return ch.Id, nil

}

func (c *Client) GetChannelsForUserAndMembers(userId, teamName string, memberUserIds []string) ([]string, error) {

	team, rs := c.RestApi.GetTeamByName(teamName, "")
	if err := HandleResponse(rs); err != nil {
		return nil, err
	}

	channels, rs := c.RestApi.GetChannelsForTeamForUser(team.Id, userId, false, "")
	if err := HandleResponse(rs); err != nil {
		return nil, err
	}

	var res []string

	for _, ch := range channels {

		if ch.Type != model.CHANNEL_PRIVATE {
			continue
		}

		chMembers, rs := c.RestApi.GetChannelMembers(ch.Id, 0, 1000, "")
		if err := HandleResponse(rs); err != nil {
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

func (c *Client) CreateBotIfNotExists(username, displayName, description, ownerId string) (string, error) {

	bots, rs := c.RestApi.GetBots(0, 1000, "")
	if err := HandleResponse(rs); err != nil {
		return "", err
	}

	for _, b := range bots {
		if b.Username == username {
			return b.UserId, nil
		}
	}

	b, rs := c.RestApi.CreateBot(&model.Bot{
		Username:    username,
		DisplayName: displayName,
		Description: description,
		OwnerId:     ownerId,
	})
	if err := HandleResponse(rs); err != nil {
		return "", err
	}

	return b.UserId, nil
}

func (c *Client) UpdateStatus(userId, status string) error {
	_, rs := c.RestApi.UpdateUserStatus(userId, &model.Status{UserId: userId, Status: status})
	return HandleResponse(rs)
}