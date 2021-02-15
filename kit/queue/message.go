package queue

import (
	"context"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	kitCtx "gitlab.medzdrav.ru/prototype/kit/context"
)

type Message struct {
	Ctx     *kitCtx.RequestContext
	Payload interface{}
}

func Decode(parentCtx context.Context, msg []byte, payload interface{}) (context.Context, error) {

	var m Message

	err := json.Unmarshal(msg, &m)
	if err != nil {
		return nil, err
	}

	_, ok := payload.(map[string]interface{})
	// if target type isn't map[string]interface{} try to decode, otherwise it's already it
	if !ok {
		err = mapstructure.Decode(m.Payload, &payload)
		if err != nil {
			return nil, err
		}
		m.Payload = payload
	} else {
		payload = m.Payload
	}

	if parentCtx == nil {
		parentCtx = context.Background()
	}

	ctx := m.Ctx.ToContext(parentCtx)

	return ctx, nil

}