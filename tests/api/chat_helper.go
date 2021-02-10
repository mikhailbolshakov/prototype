package api

import "gitlab.medzdrav.ru/prototype/kit/chat/mattermost"

func (h *TestHelper) ChatLogin(account string, goOnline bool) (*mattermost.Client, error) {

	client, err := mattermost.Login(&mattermost.Params{
		Url:     MM_URL,
		Account: account,
		Pwd:     DEFAULT_PWD,
		OpenWS:  false,
	})
	if err != nil {
		return nil, err
	}

	if goOnline {
		if err := client.SetStatus(client.User.Id, "online"); err != nil {
			return nil, err
		}
	}

	return client, nil

}
