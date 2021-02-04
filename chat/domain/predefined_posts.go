package domain

import (
	"errors"
	"fmt"
)

const (
	TP_CLIENT_NEW_REQUEST             = "client.new-request"
	TP_CLIENT_NEW_MED_REQUEST         = "client.new-med-request"
	TP_CLIENT_NEW_LAW_REQUEST         = "client.new-law-request"
	TP_CLIENT_REQUEST_ASSIGNED        = "client.request-assigned"
	TP_CONSULTANT_REQUEST_ASSIGNED    = "consultant.request-assigned"
	TP_CLIENT_NEW_EXPERT_CONSULTATION = "client.new-expert-consultation"
	TP_EXPERT_NEW_EXPERT_CONSULTATION = "expert.new-expert-consultation"
	TP_CLIENT_NO_CONSULTANT           = "client.no-consultant-available"
	TP_TASK_SOLVED                    = "client.task-solved"
	TP_CLIENT_FEEDBACK                = "client.feedback"
)

type handler func(params params) (*Post, error)
type params map[string]interface{}

var postMap = map[string]handler{
	TP_CLIENT_NEW_REQUEST:             clientNewRequest,
	TP_CLIENT_NEW_MED_REQUEST:         clientNewMedRequest,
	TP_CLIENT_NEW_LAW_REQUEST:         clientNewLawRequest,
	TP_CLIENT_REQUEST_ASSIGNED:        clientRequestAssigned,
	TP_CONSULTANT_REQUEST_ASSIGNED:    consultantRequestAssigned,
	TP_CLIENT_NEW_EXPERT_CONSULTATION: clientNewExpertConsultation,
	TP_EXPERT_NEW_EXPERT_CONSULTATION: expertNewExpertConsultation,
	TP_CLIENT_NO_CONSULTANT:           clientNoConsultantAvailable,
	TP_TASK_SOLVED:                    taskSolved,
	TP_CLIENT_FEEDBACK:                clientFeedback,
}

func (s *serviceImpl) predefinedPost(p *Post) (*Post, error) {

	if postFunc, ok := postMap[p.PredefinedPost.Code]; ok {
		predefinedPost, err := postFunc(p.PredefinedPost.Params)
		if err != nil {
			return nil, err
		}
		p.Attachments = predefinedPost.Attachments
		p.Message = predefinedPost.Message
		return p, nil
	} else {
		return nil, errors.New(fmt.Sprintf("trigger post with code %s not supported", p.PredefinedPost.Code))
	}

}

func clientNewRequest(params params) (*Post, error) {

	//	clientName, ok := params["client.name"]

	attach := []*PostAttachment{
		{
			Text:     "Идет поиск подходящего консультанта. Ожидайте....",
			ImageURL: "https://i.gifer.com/9XcW.gif",
		},
	}

	post := &Post{
		Message:     "К сожалению, я не могу ответить на Ваш вопрос,\nя найду для Вас консультанта...",
		Attachments: attach,
	}

	return post, nil
}

func clientNewMedRequest(params params) (*Post, error) {

	//	clientName, ok := params["client.name"]

	attach := []*PostAttachment{
		{
			Text:     "Идет поиск подходящего консультанта. Ожидайте....",
			ImageURL: "https://i.gifer.com/9XcW.gif",
		},
	}

	post := &Post{
		Message:     "Спасибо за обращение,\nя найду для Вас медицинского консультанта...",
		Attachments: attach,
	}

	return post, nil
}

func clientNewLawRequest(params params) (*Post, error) {

	//	clientName, ok := params["client.name"]

	attach := []*PostAttachment{
		{
			Text:     "Идет поиск подходящего консультанта. Ожидайте....",
			ImageURL: "https://i.gifer.com/9XcW.gif",
		},
	}

	post := &Post{
		Message:     "Спасибо за обращение,\nя найду для Вас консультанта-юриста...",
		Attachments: attach,
	}

	return post, nil
}

func clientNoConsultantAvailable(params params) (*Post, error) {

	attach := []*PostAttachment{
		{
			Text:     "Идет поиск подходящего консультанта. Ожидайте....",
			ImageURL: "https://i.gifer.com/9XcW.gif",
		},
	}

	post := &Post{
		Message:     "К сожалению, сегодня очень сложный день\nвсе консультанты в данный момент заняты\nПодождите, пожалуйста, я назначу первого освободившегося консультанта",
		Attachments: attach,
	}

	return post, nil
}

func clientRequestAssigned(params params) (*Post, error) {

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

	attach := []*PostAttachment{
		{
			Text:     fmt.Sprintf("#### Ваш консультант - **%s %s**", consFirstName.(string), consLastName.(string)),
			ImageURL: consUrl.(string),
		},
	}

	post := &Post{
		Attachments: attach,
	}

	return post, nil
}

func consultantRequestAssigned(params params) (*Post, error) {

	clientFirstName := params["client.first-name"]
	clientLastName := params["client.last-name"]
	clientPhone := params["client.phone"]
	clientUrl := params["client.url"]
	//clientMedcardUrl := params["client.med-card"]

	attach := []*PostAttachment{
		{
			Text:     fmt.Sprintf("###### Клиент **%s %s**\n ###### Телефон: %s", clientFirstName, clientLastName, clientPhone),
			ImageURL: clientUrl.(string),
		},
	}

	post := &Post{
		Message:     "Вам назначена консультация",
		Attachments: attach,
	}

	return post, nil
}

func clientNewExpertConsultation(params params) (*Post, error) {

	expertFirstName := params["expert.first-name"]
	expertLastName := params["expert.last-name"]
	dueDate := params["due-date"]
	expertUrl := params["expert.url"]
	expertPhotoUrl := params["expert.photo-url"]

	attach := []*PostAttachment{
		{
			Text:     fmt.Sprintf("#### Для Вас назначена консультация с экспертом %s \n #### Эксперт [**%s %s**](%s) ", dueDate, expertFirstName, expertLastName, expertUrl),
			ImageURL: expertPhotoUrl.(string),
		},
	}

	post := &Post{
		Attachments: attach,
	}

	return post, nil
}

func expertNewExpertConsultation(params params) (*Post, error) {

	clientFirstName := params["client.first-name"]
	clientLastName := params["client.last-name"]
	clientPhone := params["client.phone"]
	clientUrl := params["client.url"]
	clientMedcardUrl := params["client.med-card"]
	dueDate := params["due-date"]

	attach := []*PostAttachment{
		{
			Text:     fmt.Sprintf("#### Вам назначена консультация %s\n #### Клиент [**%s %s**](%s) \n ##### Телефон: %s \n ##### [МедКарта](%s)", dueDate, clientFirstName, clientLastName, clientMedcardUrl, clientPhone, clientMedcardUrl),
			ImageURL: clientUrl.(string),
		},
	}

	post := &Post{
		Attachments: attach,
	}

	return post, nil
}

func taskSolved(params params) (*Post, error) {

	attach := []*PostAttachment{}
	taskNum := params["task-num"]

	post := &Post{
		Message:     fmt.Sprintf("### Задача %s завершена", taskNum),
		Attachments: attach,
	}

	return post, nil
}

func clientFeedback(params params) (*Post, error) {

	attach := []*PostAttachment{
		{
			Text: "Нашим специалистам удалось решить Вашу проблему?",
			Actions: []*PostAction{
				{
					Id: "yes",
					Name: "Да, полностью",
					Integration: &PostActionIntegration{
						URL:     "http://localhost:8065/yes",
						Context: map[string]interface{}{},
					},
				},
				{
					Id: "yes-partially",
					Name: "Да, но остались вопросы",
					Integration: &PostActionIntegration{
						URL:     "http://localhost:8065/partially",
						Context: map[string]interface{}{},
					},
				},
				{
					Id: "no",
					Name: "Нет",
					Integration: &PostActionIntegration{
						URL:     "http://localhost:8065/no",
						Context: map[string]interface{}{},
					},
				},
			},
		},
	}
	post := &Post{
		Message: "",
		Attachments: attach,
	}

	return post, nil
}
