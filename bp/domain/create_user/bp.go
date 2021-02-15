package create_user

import (
	"fmt"
	"github.com/Nerzal/gocloak/v7"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	b "gitlab.medzdrav.ru/prototype/bp/domain"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	"gitlab.medzdrav.ru/prototype/kit/queue"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	pbChat "gitlab.medzdrav.ru/prototype/proto/chat"
	"log"
)

type bpImpl struct {
	userService      b.UserService
	chatService      b.ChatService
	keycloakProvider b.KeycloakProvider
	bpm.BpBase
}

func NewBp(userService b.UserService,
	chatService b.ChatService,
	bpm bpm.Engine,
	keycloak b.KeycloakProvider) b.BusinessProcess {

	bp := &bpImpl{
		userService:      userService,
		chatService:      chatService,
		keycloakProvider: keycloak,
	}
	bp.Engine = bpm

	return bp

}

func (bp *bpImpl) Init() error {

	err := bp.registerBpmHandlers()
	if err != nil {
		return err
	}
	return nil
}

func (bp *bpImpl) SetQueueListeners(ql listener.QueueListener) {
	ql.Add(queue.QUEUE_TYPE_AT_LEAST_ONCE, "users.draft-created", bp.userDraftCreatedHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-create-user"
}

func (bp *bpImpl) GetBPMNPath() string {
	return "../bp/domain/create_user/bp.bpmn"
}

func (bp *bpImpl) registerBpmHandlers() error {
	return bp.RegisterTaskHandlers(map[string]interface{}{
		"st-create-mm-user":    bp.createMMUser,
		"st-create-mm-channel": bp.createMMChannel,
		"st-create-send-hello": bp.sendHello,
		"st-create-kk-user":    bp.createKKUser,
		"st-activate-user":     bp.activateUser,
		"st-delete-mm-user":    bp.deleteMMUser,
		"st-delete-kk-user":    bp.deleteKKUser,
		"st-delete-user":       bp.deleteUser,
	})
}

func (bp *bpImpl) userDraftCreatedHandler(msg []byte) error {

	var user map[string]interface{}
	ctx, err := queue.Decode(nil, msg, &user)
	if err != nil {
		return err
	}

	var vars = map[string]interface{}{"id": user["id"], "type": user["type"]}
	if err := zeebe.CtxToVars(ctx, vars); err != nil {
		return err
	}

	_, err = bp.StartProcess("p-create-user", vars)

	return err

}

func (bp *bpImpl) createMMUser(client worker.JobClient, job entities.Job) {

	log.Println("createMMUser executed")

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["id"].(string)
	user := bp.userService.Get(ctx, userId)

	var email string
	switch user.Type {
	case "client":
		email = user.ClientDetails.Email
	case "consultant":
		email = user.ConsultantDetails.Email
	case "expert":
		email = user.ExpertDetails.Email
	default:
		zeebe.FailJob(client, job, fmt.Errorf("unknow user type %s", user.Type))
		return
	}

	//create a new user in MM
	chatUserId, err := bp.chatService.CreateUser(ctx, &pbChat.CreateUserRequest{
		Username: user.Username,
		Email:    email,
	})
	if err != nil {
		err = bp.SendError(job.GetKey(), "err-create-mm-user", err.Error())
		if err != nil {
			zeebe.FailJob(client, job, err)
		}
		return
	}

	variables["mmId"] = chatUserId

	_, err = bp.userService.SetMMUserId(ctx, userId, chatUserId)
	if err != nil {
		err = bp.SendError(job.GetKey(), "err-create-mm-user", err.Error())
		if err != nil {
			zeebe.FailJob(client, job, err)
		}
		return
	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) createMMChannel(client worker.JobClient, job entities.Job) {

	log.Println("createMMChannel executed")

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["id"].(string)
	user := bp.userService.Get(ctx, userId)

	if user.Type == "client" {
		firstName := user.ClientDetails.FirstName
		lastName := user.ClientDetails.LastName

		channelId, err := bp.chatService.CreateClientChannel(ctx, &pbChat.CreateClientChannelRequest{
			ChatUserId:  user.MMId,
			Name:        user.MMId,
			DisplayName: fmt.Sprintf("Общие вопросы (клиент %s %s)", firstName, lastName),
		})
		if err != nil {
			err = bp.SendError(job.GetKey(), "err-create-mm-channel", err.Error())
			if err != nil {
				zeebe.FailJob(client, job, err)
			}
			return
		}

		variables["mmChannelId"] = channelId
		user.ClientDetails.CommonChannelId = channelId

		_, err = bp.userService.SetClientDetails(ctx, userId, user.ClientDetails)
		if err != nil {
			err = bp.SendError(job.GetKey(), "err-create-mm-channel", err.Error())
			if err != nil {
				zeebe.FailJob(client, job, err)
			}
			return
		}

	} else {
		err = bp.SendError(job.GetKey(), "err-create-mm-channel", "this operation is valid for clients only")
		if err != nil {
			zeebe.FailJob(client, job, err)
		}
		return
	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

}

func (bp *bpImpl) sendHello(client worker.JobClient, job entities.Job) {

	log.Println("sendHello executed")

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	id := variables["id"].(string)
	channelId := variables["mmChannelId"].(string)

	if channelId != "" {
		user := bp.userService.Get(ctx, id)

		if user.ClientDetails != nil {
			msg := fmt.Sprintf("Добрый день, **%s %s**\nДобро пожаловать!!!", user.ClientDetails.FirstName, user.ClientDetails.MiddleName)
			if err := bp.chatService.Post(ctx, msg, channelId, "", false); err != nil {
				zeebe.FailJob(client, job, err)
				return
			}
		}

	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) createKKUser(client worker.JobClient, job entities.Job) {

	log.Println("createKKUser executed")

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["id"].(string)
	user := bp.userService.Get(ctx, userId)

	// TODO: config
	token, err := bp.keycloakProvider().LoginAdmin(ctx, "admin", "admin", "master")
	if err != nil {
		err = bp.SendError(job.GetKey(), "err-create-kk-user", err.Error())
		if err != nil {
			zeebe.FailJob(client, job, err)
		}
		return
	}

	kkUser := gocloak.User{
		Enabled:  gocloak.BoolP(true),
		Username: gocloak.StringP(user.Username),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP("12345"),
				Temporary: gocloak.BoolP(false),
			},
		},
	}

	switch user.Type {
	case "client":
		kkUser.FirstName = gocloak.StringP(user.ClientDetails.FirstName)
		kkUser.LastName = gocloak.StringP(user.ClientDetails.LastName)
		kkUser.Email = gocloak.StringP(user.ClientDetails.Email)
	case "consultant":
		kkUser.FirstName = gocloak.StringP(user.ConsultantDetails.FirstName)
		kkUser.LastName = gocloak.StringP(user.ConsultantDetails.LastName)
		kkUser.Email = gocloak.StringP(user.ConsultantDetails.Email)
	case "expert":
		kkUser.FirstName = gocloak.StringP(user.ExpertDetails.FirstName)
		kkUser.LastName = gocloak.StringP(user.ExpertDetails.LastName)
		kkUser.Email = gocloak.StringP(user.ExpertDetails.Email)
	default:
		zeebe.FailJob(client, job, fmt.Errorf("unknow user type %s", user.Type))
		return
	}

	kkId, err := bp.keycloakProvider().CreateUser(ctx, token.AccessToken, "prototype", kkUser)
	if err != nil {
		err = bp.SendError(job.GetKey(), "err-create-kk-user", err.Error())
		if err != nil {
			zeebe.FailJob(client, job, err)
		}
		return
	}

	variables["kkId"] = kkId

	_, err = bp.userService.SetKKUserId(ctx, user.Id, kkId)
	if err != nil {
		err = bp.SendError(job.GetKey(), "err-create-kk-user", err.Error())
		if err != nil {
			zeebe.FailJob(client, job, err)
		}
		return
	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) activateUser(client worker.JobClient, job entities.Job) {

	log.Println("activateUser executed")

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["id"].(string)
	_, err = bp.userService.Activate(ctx, userId)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) deleteMMUser(client worker.JobClient, job entities.Job) {

	log.Println("deleteMMUser executed")

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = bp.chatService.DeleteUser(ctx, variables["mmId"].(string))
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) deleteKKUser(client worker.JobClient, job entities.Job) {

	log.Println("deleteKKUser executed")

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	// TODO: config

	token, err := bp.keycloakProvider().LoginAdmin(ctx, "admin", "admin", "master")
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	kkUserId := variables["kkId"].(string)
	err = bp.keycloakProvider().DeleteUser(ctx, token.AccessToken, "prototype", kkUserId)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}

func (bp *bpImpl) deleteUser(client worker.JobClient, job entities.Job) {

	log.Println("deleteUser executed")

	variables, ctx, err := zeebe.GetVarsAndCtx(job)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	_, err = bp.userService.Delete(ctx, variables["id"].(string))
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = zeebe.CompleteJob(client, job, variables)
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}
}
