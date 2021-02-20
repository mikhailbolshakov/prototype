package chat

const (
 	QUEUE_TOPIC_MATTERMOST_POST_MESSAGE =  "mm.posts"
)

type MattermostPostMessagePayload struct {
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
