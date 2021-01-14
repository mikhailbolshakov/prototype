package domain

import (
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/mm/repository/adapters/mattermost"
)

const (
	TP_CLIENT_NEW_REQUEST = "client.new-request"
	TP_CLIENT_REQUEST_ASSIGNED = "client.request-assigned"
	TP_CONSULTANT_REQUEST_ASSIGNED = "consultant.request-assigned"
	TP_CLIENT_NEW_EXPERT_CONSULTATION = "client.new-expert-consultation"
	TP_EXPERT_NEW_EXPERT_CONSULTATION = "expert.new-expert-consultation"
	TP_TASK_REMINDER = "task-reminder"
)

type triggerPostHandler func(params triggerPostParams) (*mattermost.Post, error)
type triggerPostParams map[string]interface{}

var postMap = map[string]triggerPostHandler{
	TP_CLIENT_NEW_REQUEST: clientNewRequest,
	TP_CLIENT_REQUEST_ASSIGNED: clientRequestAssigned,
	TP_CONSULTANT_REQUEST_ASSIGNED: consultantRequestAssigned,
	TP_CLIENT_NEW_EXPERT_CONSULTATION: clientNewExpertConsultation,
	TP_EXPERT_NEW_EXPERT_CONSULTATION: expertNewExpertConsultation,
	TP_TASK_REMINDER: taskReminder,
}

func (s *serviceImpl) sendTriggerPost(postCode string, userId string, channelId string, params triggerPostParams) error {

	if postFunc, ok := postMap[postCode]; ok {
		post, err := postFunc(params)
		if err != nil {
			return err
		}

		post.ChannelId = channelId

		if err := s.mmService.CreateEphemeralPost(&mattermost.EphemeralPost{
			Post:   post,
			UserId: userId,
		}); err != nil {
			return err
		}

		//if err := s.mmService.CreatePost(post); err != nil {
		//	return err
		//}

	} else {
		return errors.New(fmt.Sprintf("trigger post with code %s not supported", postCode))
	}

	return nil
}

func clientNewRequest(params triggerPostParams) (*mattermost.Post, error) {

	clientName, ok := params["client.name"]
	if !ok {
		return nil, errors.New("parameter 'client.name' is empty")
	}
	
	attach := []*mattermost.PostAttachment{
		{
			Text:    fmt.Sprintf("## добрый день, **%s** \n ### Мы подбираем для Вас консультанта...", clientName),
			ImageURL:   "https://i.gifer.com/VAyR.gif",
		},
	}
	
	post := &mattermost.Post{
		Attachments: attach,
	}

	return post, nil
}

func clientRequestAssigned(params triggerPostParams) (*mattermost.Post, error) {

	consFirstName, ok := params["consultant.first-name"]
	if !ok {
		return nil, errors.New("parameter 'consultant.first-name' is empty")
	}

	consLastName, ok := params["consultant.last-name"]
	if !ok {
		return nil, errors.New("parameter 'consultant.last-name' is empty")
	}

	consUrl, ok := params["consultant.url"]
	if !ok {
		return nil, errors.New("parameter 'consultant.url' is empty")
	}

	attach := []*mattermost.PostAttachment{
		{
			Text:    fmt.Sprintf("## Ваш консультант - **%s %s**", consFirstName.(string), consLastName.(string)),
			ImageURL: consUrl.(string),
		},
	}

	post := &mattermost.Post{
		Attachments: attach,
	}

	return post, nil
}

func consultantRequestAssigned(params triggerPostParams) (*mattermost.Post, error) {

	clientFirstName, ok := params["client.first-name"]
	if !ok {
		return nil, errors.New("parameter 'client.first-name' is empty")
	}

	clientLastName, ok := params["client.last-name"]
	if !ok {
		return nil, errors.New("parameter 'client.last-name' is empty")
	}

	clientPhone, ok := params["client.phone"]
	if !ok {
		return nil, errors.New("parameter 'client.phone' is empty")
	}

	clientUrl, ok := params["client.url"]
	if !ok {
		return nil, errors.New("parameter 'client.url' is empty")
	}

	clientMedcardUrl, ok := params["client.med-card"]
	if !ok {
		return nil, errors.New("parameter 'client.med-card' is empty")
	}

	attach := []*mattermost.PostAttachment{
		{
			Text:    fmt.Sprintf("## Вам назначена консультация \n ### Клиент [**%s %s**](%s) \n #### Телефон: %s \n #### [МедКарта](%s)", clientFirstName, clientLastName, clientMedcardUrl, clientPhone, clientMedcardUrl),
			ImageURL:   clientUrl.(string),
		},
	}

	post := &mattermost.Post{
		Attachments: attach,
	}

	return post, nil
}

func clientNewExpertConsultation(params triggerPostParams) (*mattermost.Post, error) {

	expertFirstName := params["expert.first-name"]
	expertLastName := params["expert.last-name"]
	dueDate := params["due-date"]
	expertUrl := params["expert.url"]
	expertPhotoUrl := params["expert.photo-url"]

	// https://pmed.moi-service.ru/doctor/7712?cityIdDF=1
	attach := []*mattermost.PostAttachment{
		{
			Text:    fmt.Sprintf("### Для Вас назначена консультация с экспертом %s \n ##  Эксперт [**%s %s**](%s) ", dueDate, expertFirstName, expertLastName, expertUrl),
			ImageURL:   expertPhotoUrl.(string),
		},
	}

	post := &mattermost.Post{
		Attachments: attach,
	}

	return post, nil
}

func expertNewExpertConsultation(params triggerPostParams) (*mattermost.Post, error) {

	clientFirstName := params["client.first-name"]
	clientLastName := params["client.last-name"]
	clientPhone := params["client.phone"]
	clientUrl := params["client.url"]
	clientMedcardUrl := params["client.med-card"]
	dueDate := params["due-date"]

	attach := []*mattermost.PostAttachment{
		{
			Text:    fmt.Sprintf("## Вам назначена консультация %s\n ### Клиент [**%s %s**](%s) \n #### Телефон: %s \n #### [МедКарта](%s)", dueDate, clientFirstName, clientLastName, clientMedcardUrl, clientPhone, clientMedcardUrl),
			ImageURL:   clientUrl.(string),
		},
	}

	post := &mattermost.Post{
		Attachments: attach,
	}

	return post, nil
}

func taskReminder(params triggerPostParams) (*mattermost.Post, error) {

	attach := []*mattermost.PostAttachment{}

	post := &mattermost.Post{
		Message: "Напоминание о задаче",
		Attachments: attach,
	}

	return post, nil
}