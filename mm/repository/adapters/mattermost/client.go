package mattermost

import (
	"encoding/json"
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
}

type Client struct {
	RestApi *model.Client4
	WsApi   *model.WebSocketClient
	Token   string
	User    *model.User
	Quit    chan interface{}
	Params  *Params
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

func login(p *Params) (*Client, error) {

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

func (c *Client) logout() error {
	_, rs := c.RestApi.Logout()
	if err := handleResponse(rs); err != nil {
		return err
	}
	return nil
}

func (c *Client) createUser(rq *CreateUserRequest) (*CreateUserResponse, error) {

	user, rs := c.RestApi.CreateUser(&model.User{
		Username: rq.Username,
		Password: rq.Password,
		Email:    rq.Email,
	})
	if err := handleResponse(rs); err != nil {
		return nil, err
	}
	log.Printf("By created. Id: %s, email: %s", user.Id, user.Email)

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

func (c *Client) deleteUser(userId string) error {

	_, rs := c.RestApi.DeleteUser(userId)
	if err := handleResponse(rs); err != nil {
		return err
	}
	log.Printf("user delete. Id: %s", userId)

	return nil

}

func (c *Client) createClientChannel(rq *CreateClientChannelRequest) (*CreateChannelResponse, error) {

	team, rs := c.RestApi.GetTeamByName(rq.TeamName, "")
	if err := handleResponse(rs); err != nil {
		return nil, err
	}

	ch, rs := c.RestApi.CreateChannel(&model.Channel{
		TeamId:      team.Id,
		Type:        "P",
		DisplayName: rq.DisplayName,
		Name:        rq.Name,
	})
	if err := handleResponse(rs); err != nil {
		return nil, err
	}

	_, rs = c.RestApi.AddChannelMember(ch.Id, rq.ClientUserId)
	if err := handleResponse(rs); err != nil {
		return nil, err
	}

	return &CreateChannelResponse{ChannelId: ch.Id}, nil
}

func (c *Client) subscribeUser(channelId string, userId string) error {
	_, rs := c.RestApi.AddChannelMember(channelId, userId)
	if err := handleResponse(rs); err != nil {
		return err
	}
	return nil
}

func (c *Client) convertAttachments(attachments []*PostAttachment) []*model.SlackAttachment {

	var slackAttachments []*model.SlackAttachment

	for _, a := range attachments {

		sa := &model.SlackAttachment{
			Fallback:   a.Fallback,
			Color:      a.Color,
			Pretext:    a.Pretext,
			AuthorName: a.AuthorName,
			AuthorLink: a.AuthorLink,
			AuthorIcon: a.AuthorIcon,
			Title:      a.Title,
			TitleLink:  a.TitleLink,
			Text:       a.Text,
			ImageURL:   a.ImageURL,
			ThumbURL:   a.ThumbURL,
			Footer:     a.Footer,
			FooterIcon: a.FooterIcon,
		}

		if a.Fields != nil && len(a.Fields) > 0 {
			sa.Fields = []*model.SlackAttachmentField{}

			for _, f := range a.Fields {
				sa.Fields = append(sa.Fields, &model.SlackAttachmentField{
					Title: f.Title,
					Value: f.Value,
					Short: model.SlackCompatibleBool(f.Short),
				})

			}
		}

		slackAttachments = append(slackAttachments, sa)

	}

	return slackAttachments
}

func (c *Client) createEphemeralPost(channelId string, recipientUserId string, message string, attachments []*PostAttachment) error {

	props := model.StringInterface{}
	if attachments != nil && len(attachments) > 0 {
		props["attachments"] = c.convertAttachments(attachments)
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

func (c *Client) createPost(channelId string, message string, attachments []*PostAttachment) error {

	props := model.StringInterface{}
	if attachments != nil && len(attachments) > 0 {
		props["attachments"] = c.convertAttachments(attachments)
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

func (c *Client) getUsersStatuses(userIds []string) ([]*UserStatus, error) {

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

func (c *Client) createDirectChannel(userId1, userId2 string) (*CreateChannelResponse, error) {

	ch, rs := c.RestApi.CreateGroupChannel([]string{userId1, userId2})
	if err := handleResponse(rs); err != nil {
		return nil, err
	}

	return &CreateChannelResponse{ChannelId: ch.Id}, nil

}

func (c *Client) getChannelsForUserAndMembers(userId, teamName string, memberUserIds []string) ([]string, error) {

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
