package context

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gitlab.medzdrav.ru/prototype/kit"
	"google.golang.org/grpc/metadata"
)

type requestContextKey struct{}

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

func NewRequestCtx() *RequestContext {
	return &RequestContext{}
}

func (r *RequestContext) GetRequestId() string {
	return r.Rid
}

func (r *RequestContext) GetSessionId() string {
	return r.Sid
}

func (r *RequestContext) GetUserId() string {
	return r.Uid
}

func (r *RequestContext) GetChatUserId() string {
	return r.Cid
}

func (r *RequestContext) GetClientType() string {
	return r.Cl
}

func (r *RequestContext) GetUsername() string {
	return r.Un
}

func (r *RequestContext) Empty() *RequestContext {

	return &RequestContext{
		Rid: kit.Nil(),
		Sid: kit.Nil(),
		Uid: kit.Nil(),
		Un:  "",
		Cid: "",
		Cl:  "none",
	}
}

func (r *RequestContext) WithRequestId(requestId string) *RequestContext {
	r.Rid = requestId
	return r
}

func (r *RequestContext) WithNewRequestId() *RequestContext {
	r.Rid = kit.NewId()
	return r
}

func (r *RequestContext) WithSessionId(sessionId string) *RequestContext {
	r.Sid = sessionId
	return r
}

func (r *RequestContext) WithChatUserId(chatUserId string) *RequestContext {
	r.Cid = chatUserId
	return r
}

func (r *RequestContext) Rest() *RequestContext {
	r.Cl = "rest"
	return r
}

func (r *RequestContext) Test() *RequestContext {
	r.Cl = "test"
	return r
}

func (r *RequestContext) Batch() *RequestContext {
	r.Cl = "batch"
	return r
}

func (r *RequestContext) Queue() *RequestContext {
	r.Cl = "queue"
	return r
}

func (r *RequestContext) Client(client string) *RequestContext {
	r.Cl = client
	return r
}

func (r *RequestContext) WithUser(userId, username string) *RequestContext {
	r.Uid = userId
	r.Un = username
	return r
}

func (r *RequestContext) ToContext(parent context.Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithValue(parent, requestContextKey{}, r)
}

func Request(context context.Context) (*RequestContext, bool) {
	if r, ok := context.Value(requestContextKey{}).(*RequestContext); ok {
		return r, true
	}
	return &RequestContext{}, false
}

func MustRequest(context context.Context) (*RequestContext, error) {
	if r, ok := context.Value(requestContextKey{}).(*RequestContext); ok {
		return r, nil
	}
	return &RequestContext{}, fmt.Errorf("context invalid")
}

func FromContextToGrpcMD(ctx context.Context) (metadata.MD, bool) {
	if r, ok := Request(ctx); ok {
		rm, _ := json.Marshal(*r)
		return metadata.Pairs("rq-bin", string(rm)), true
	}
	return metadata.Pairs(), false
}

func FromGrpcMD(ctx context.Context, md metadata.MD) context.Context {

	if rqb, ok := md["rq-bin"]; ok {
		if len(rqb) > 0 {
			rm := []byte(rqb[0])
			rq := &RequestContext{}
			_ = json.Unmarshal(rm, rq)
			return context.WithValue(ctx, requestContextKey{}, rq)
		}
	}
	return ctx
}

func FromMap(ctx context.Context, mp map[string]interface{}) (context.Context, error) {
	var r *RequestContext
	err := mapstructure.Decode(mp, &r)
	if err != nil {
		return nil, err
	}
	return r.ToContext(ctx), nil
}
