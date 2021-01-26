package create_user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nerzal/gocloak/v7"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
	bpm2 "gitlab.medzdrav.ru/prototype/bp/bpm"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/bp/repository/adapters/users"
	"gitlab.medzdrav.ru/prototype/kit/bpm"
	"gitlab.medzdrav.ru/prototype/kit/bpm/zeebe"
	"gitlab.medzdrav.ru/prototype/kit/queue/listener"
	"gitlab.medzdrav.ru/prototype/proto/mm"
	"log"
)

type bpImpl struct {
	userService users.Service
	mmService   mattermost.Service
	keycloak    gocloak.GoCloak
	bpm.BpBase
}

func NewBp(userService users.Service,
	mmService mattermost.Service,
	bpm bpm.Engine,
	keycloak gocloak.GoCloak) bpm2.BusinessProcess {

	bp := &bpImpl{
		userService: userService,
		mmService:   mmService,
		keycloak:    keycloak,
	}
	bp.Engine = bpm

	return bp

}

func (bp *bpImpl) Init() error {

	err := bp.registerBpmHandlers()
	if err != nil {
		return err
	}

	if err := bp.DeployBPMNs([]string{"../bp/bpm/create_user/bp.bpmn"}); err != nil {
		return err
	}

	return nil
}

func (bp *bpImpl) SetQueueListeners(ql listener.QueueListener) {
	ql.Add("users.draft-created", bp.userDraftCreatedHandler)
}

func (bp *bpImpl) GetId() string {
	return "p-create-user"
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
		"st-delete-user":    	bp.deleteUser,
	})
}

func (bp *bpImpl) userDraftCreatedHandler(payload []byte) error {

	var user map[string]interface{}
	if err := json.Unmarshal(payload, &user); err != nil {
		return err
	}

	var v = map[string]interface{}{"id": user["id"], "type": user["type"]}
	_, err := bp.StartProcess("p-create-user", v)

	return err

}

func (bp *bpImpl) createMMUser(client worker.JobClient, job entities.Job) {

	log.Println("createMMUser executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["id"].(string)
	user := bp.userService.Get(userId)

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
	mmRs, err := bp.mmService.CreateUser(&mm.CreateUserRequest{
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

	variables["mmId"] = mmRs.Id

	_, err = bp.userService.SetMMUserId(userId, mmRs.Id)
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

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["id"].(string)
	user := bp.userService.Get(userId)

	if user.Type == "client" {
		firstName := user.ClientDetails.FirstName
		lastName := user.ClientDetails.LastName

		chRs, err := bp.mmService.CreateClientChannel(&mm.CreateClientChannelRequest{
			ClientUserId: user.MMId,
			Name:         user.MMId,
			DisplayName:  fmt.Sprintf("Клиент %s %s - консультант", firstName, lastName),
		})
		if err != nil {
			err = bp.SendError(job.GetKey(), "err-create-mm-channel", err.Error())
			if err != nil {
				zeebe.FailJob(client, job, err)
			}
			return
		}

		variables["mmChannelId"] = chRs.ChannelId
		user.ClientDetails.MMChannelId = chRs.ChannelId

		_, err = bp.userService.SetClientDetails(userId, user.ClientDetails)
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

	variables, err := job.GetVariablesAsMap()
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

func (bp *bpImpl) createKKUser(client worker.JobClient, job entities.Job) {

	log.Println("createKKUser executed")

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["id"].(string)
	user := bp.userService.Get(userId)

	// TODO: move it from here
	token, err := bp.keycloak.LoginAdmin(context.Background(), "admin", "admin", "master")
	if err != nil {
		err = bp.SendError(job.GetKey(), "err-create-kk-user", err.Error())
		if err != nil {
			zeebe.FailJob(client, job, err)
		}
		return
	}

	kkUser := gocloak.User{
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(user.Username),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Type: gocloak.StringP("password"),
				Value: gocloak.StringP("12345"),
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

	kkId, err := bp.keycloak.CreateUser(context.Background(), token.AccessToken, "prototype", kkUser)
	if err != nil {
		err = bp.SendError(job.GetKey(), "err-create-kk-user", err.Error())
		if err != nil {
			zeebe.FailJob(client, job, err)
		}
		return
	}

	variables["kkId"] = kkId

	_, err = bp.userService.SetKKUserId(user.Id, kkId)
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

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	userId := variables["id"].(string)
	_, err = bp.userService.Activate(userId)
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

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	err = bp.mmService.DeleteUser(variables["mmId"].(string))
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

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	token, err := bp.keycloak.LoginAdmin(context.Background(), "admin", "admin", "master")
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	kkUserId := variables["kkId"].(string)
	err = bp.keycloak.DeleteUser(context.Background(), token.AccessToken, "prototype", kkUserId)
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

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		zeebe.FailJob(client, job, err)
		return
	}

	_, err = bp.userService.Delete(variables["id"].(string))
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