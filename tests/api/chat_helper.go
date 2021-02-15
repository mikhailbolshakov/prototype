package api

import (
	"encoding/json"
	"fmt"
	chatApi "gitlab.medzdrav.ru/prototype/api/public/chat"
)

func (h *TestHelper) SetStatus(userId, status string) error {
	_, err := h.PUT(fmt.Sprintf("%s/api/chat/users/%s/status/%s", BASE_URL, userId, status), []byte{})
	return err
}

func (h *TestHelper) MyPost(channelId, message string) error {

	pr := &chatApi.PostRequest{
		ChannelId: channelId,
		Message:   message,
	}
	prJ, _ := json.Marshal(pr)

	_, err := h.POST(fmt.Sprintf("%s/api/chat/users/me/posts", BASE_URL), prJ)
	return err
}