package main

import (
	"encoding/json"
	"fmt"
	"github.com/adacta-ru/mattermost-server/v6/model"
	"github.com/nats-io/stan.go"
	"time"
)

const (
	CLIENT = "mattermost"
)

// Message corresponds app message format
type Message struct {
	Ctx     *RequestContext `json:"ctx"`
	Payload interface{}     `json:"pl"`
}

type RequestContext struct {
	// request ID
	Rid string `json:"rid"`
	// session ID
	Sid string `json:"sid"`
	// user ID
	Uid string `json:"uid"`
	// username
	Un string `json:"un"`
	// chat user id
	Cid string `json:"cid"`
	// client type
	Cl string `json:"cl"`
}

type postStanMsg struct {
	Id        string `json:"id"`
	CreateAt  int64  `json:"createAt"`
	UpdateAt  int64  `json:"updateAt"`
	EditAt    int64  `json:"editAt"`
	DeleteAt  int64  `json:"deleteAt"`
	UserId    string `json:"userId"`
	ChannelId string `json:"channelId"`
	Message   string `json:"message"`
	Type      string `json:"type"`
}

func (p *Plugin) reconnectHandler(c stan.Conn, err error) {

	for {
		select {
			case <-time.NewTicker(time.Second).C:
				err := p.ConnectStan()
				if err != nil {
					p.API.LogDebug("[STAN] reconnection failed")
				} else {
					p.API.LogDebug("[STAN] reconnected")
					return
				}
			case <-p.close:
				return
		}

	}

}

func (p *Plugin) ConnectStan() error {

	if p.stanConn != nil {
		p.stanConn.Close()
	}

	sc, err := stan.Connect(p.cfg.NatsClusterId, CLIENT, stan.NatsURL(p.cfg.NatsUrl), stan.SetConnectionLostHandler(p.reconnectHandler))
	if err != nil {
		return err
	}
	p.stanConn = sc
	p.API.LogDebug("[STAN] connected")
	return nil
}

func (p *Plugin) StanPublishPost(post *model.Post) {

	var err error
	pl := &postStanMsg{
		Id:        post.Id,
		CreateAt:  post.CreateAt,
		UpdateAt:  post.UpdateAt,
		EditAt:    post.EditAt,
		DeleteAt:  post.DeleteAt,
		UserId:    post.UserId,
		ChannelId: post.ChannelId,
		Message:   post.Message,
		Type:      post.Type,
	}

	msg := &Message{
		Ctx:     &RequestContext{
			Rid: model.NewId(),
			Cid: post.UserId,
			Cl:  "mm",
		},
		Payload: pl,
	}

	j, err := json.Marshal(msg)
	if err != nil {
		p.API.LogError(fmt.Sprintf("[STAN] marshal error: %v", err))
		return
	}

	err = p.stanConn.Publish(p.cfg.TopicNewPost, j)
	if err != nil {
		p.API.LogError(fmt.Sprintf("[STAN] message publish error: %v", err))
		return
	}
	p.API.LogDebug(fmt.Sprintf("[STAN] post published %v", string(j)))

}

func (p *Plugin) StanClose() error {
	if p.stanConn != nil {
		p.stanConn.Close()
		p.API.LogDebug("[STAN] disconnected")
	}
	return nil
}
