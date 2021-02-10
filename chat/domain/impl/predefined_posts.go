package impl

import (
	"errors"
	"fmt"
	"gitlab.medzdrav.ru/prototype/chat/domain"
)

type handler func(params params) (*domain.Post, error)
type params map[string]interface{}

var postMap = map[string]handler{
	domain.TP_CLIENT_NEW_REQUEST:             clientNewRequest,
	domain.TP_CLIENT_NEW_MED_REQUEST:         clientNewMedRequest,
	domain.TP_CLIENT_NEW_LAW_REQUEST:         clientNewLawRequest,
	domain.TP_CLIENT_REQUEST_ASSIGNED:        clientRequestAssigned,
	domain.TP_CONSULTANT_REQUEST_ASSIGNED:    consultantRequestAssigned,
	domain.TP_CLIENT_NEW_EXPERT_CONSULTATION: clientNewExpertConsultation,
	domain.TP_EXPERT_NEW_EXPERT_CONSULTATION: expertNewExpertConsultation,
	domain.TP_CLIENT_NO_CONSULTANT:           clientNoConsultantAvailable,
	domain.TP_TASK_SOLVED:                    taskSolved,
	domain.TP_CLIENT_FEEDBACK:                clientFeedback,
}

func (s *serviceImpl) predefinedPost(p *domain.Post) (*domain.Post, error) {

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

func clientNewRequest(params params) (*domain.Post, error) {

	//	clientName, ok := params["client.name"]

	attach := []*domain.PostAttachment{
		{
			Text:     "Идет поиск подходящего консультанта. Ожидайте....",
			ImageURL: "https://i.gifer.com/9XcW.gif",
		},
	}

	post := &domain.Post{
		Message:     "К сожалению, я не могу ответить на Ваш вопрос,\nя найду для Вас консультанта...",
		Attachments: attach,
	}

	return post, nil
}

func clientNewMedRequest(params params) (*domain.Post, error) {

	//	clientName, ok := params["client.name"]

	attach := []*domain.PostAttachment{
		{
			Text:     "Идет поиск подходящего консультанта. Ожидайте....",
			ImageURL: "https://i.gifer.com/9XcW.gif",
		},
	}

	post := &domain.Post{
		Message:     "Спасибо за обращение,\nя найду для Вас медицинского консультанта...",
		Attachments: attach,
	}

	return post, nil
}

func clientNewLawRequest(params params) (*domain.Post, error) {

	//	clientName, ok := params["client.name"]

	attach := []*domain.PostAttachment{
		{
			Text:     "Идет поиск подходящего консультанта. Ожидайте....",
			ImageURL: "https://i.gifer.com/9XcW.gif",
		},
	}

	post := &domain.Post{
		Message:     "Спасибо за обращение,\nя найду для Вас консультанта-юриста...",
		Attachments: attach,
	}

	return post, nil
}

func clientNoConsultantAvailable(params params) (*domain.Post, error) {

	attach := []*domain.PostAttachment{
		{
			Text:     "Идет поиск подходящего консультанта. Ожидайте....",
			ImageURL: "https://i.gifer.com/9XcW.gif",
		},
	}

	post := &domain.Post{
		Message:     "К сожалению, сегодня очень сложный день\nвсе консультанты в данный момент заняты\nПодождите, пожалуйста, я назначу первого освободившегося консультанта",
		Attachments: attach,
	}

	return post, nil
}

func clientRequestAssigned(params params) (*domain.Post, error) {

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

	attach := []*domain.PostAttachment{
		{
			Text:     fmt.Sprintf("#### Ваш консультант - **%s %s**", consFirstName.(string), consLastName.(string)),
			ImageURL: consUrl.(string),
		},
	}

	post := &domain.Post{
		Attachments: attach,
	}

	return post, nil
}

func consultantRequestAssigned(params params) (*domain.Post, error) {

	clientFirstName := params["client.first-name"]
	clientLastName := params["client.last-name"]
	clientPhone := params["client.phone"]
	clientUrl := params["client.url"]
	//clientMedcardUrl := params["client.med-card"]

	attach := []*domain.PostAttachment{
		{
			Text:     fmt.Sprintf("###### Клиент **%s %s**\n ###### Телефон: %s", clientFirstName, clientLastName, clientPhone),
			ImageURL: clientUrl.(string),
		},
	}

	post := &domain.Post{
		Message:     "Вам назначена консультация",
		Attachments: attach,
	}

	return post, nil
}

func clientNewExpertConsultation(params params) (*domain.Post, error) {

	expertFirstName := params["expert.first-name"]
	expertLastName := params["expert.last-name"]
	dueDate := params["due-date"]
	expertUrl := params["expert.url"]
	expertPhotoUrl := params["expert.photo-url"]

	attach := []*domain.PostAttachment{
		{
			Text:     fmt.Sprintf("#### Для Вас назначена консультация с экспертом %s \n #### Эксперт [**%s %s**](%s) ", dueDate, expertFirstName, expertLastName, expertUrl),
			ImageURL: expertPhotoUrl.(string),
		},
	}

	post := &domain.Post{
		Attachments: attach,
	}

	return post, nil
}

func expertNewExpertConsultation(params params) (*domain.Post, error) {

	clientFirstName := params["client.first-name"]
	clientLastName := params["client.last-name"]
	clientPhone := params["client.phone"]
	clientUrl := params["client.url"]
	clientMedcardUrl := params["client.med-card"]
	dueDate := params["due-date"]

	attach := []*domain.PostAttachment{
		{
			Text:     fmt.Sprintf("#### Вам назначена консультация %s\n #### Клиент [**%s %s**](%s) \n ##### Телефон: %s \n ##### [МедКарта](%s)", dueDate, clientFirstName, clientLastName, clientMedcardUrl, clientPhone, clientMedcardUrl),
			ImageURL: clientUrl.(string),
		},
	}

	post := &domain.Post{
		Attachments: attach,
	}

	return post, nil
}

func taskSolved(params params) (*domain.Post, error) {

	attach := []*domain.PostAttachment{}
	taskNum := params["task-num"]

	post := &domain.Post{
		Message:     fmt.Sprintf("### Задача %s завершена", taskNum),
		Attachments: attach,
	}

	return post, nil
}

func clientFeedback(params params) (*domain.Post, error) {

	attach := []*domain.PostAttachment{
		{
			Text: "Нашим специалистам удалось решить Вашу проблему?",
			Actions: []*domain.PostAction{
				{
					Id: "yes",
					Name: "Да, полностью",
					Integration: &domain.PostActionIntegration{
						URL:     "http://localhost:8065/yes",
						Context: map[string]interface{}{},
					},
				},
				{
					Id: "yes-partially",
					Name: "Да, но остались вопросы",
					Integration: &domain.PostActionIntegration{
						URL:     "http://localhost:8065/partially",
						Context: map[string]interface{}{},
					},
				},
				{
					Id: "no",
					Name: "Нет",
					Integration: &domain.PostActionIntegration{
						URL:     "http://localhost:8065/no",
						Context: map[string]interface{}{},
					},
				},
			},
		},
	}
	post := &domain.Post{
		Message: "",
		Attachments: attach,
	}

	return post, nil
}
