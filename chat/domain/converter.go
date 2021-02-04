package domain

import (
	"gitlab.medzdrav.ru/prototype/chat/repository/adapters/mattermost"
)

func (s *serviceImpl) convertAttachments(attachments []*PostAttachment) []*mattermost.PostAttachment {

	var repAttachments []*mattermost.PostAttachment

	if attachments == nil || len(attachments) == 0 {
		return repAttachments
	}

	for _, a := range attachments {

		sa := &mattermost.PostAttachment{
			Fallback:   a.Fallback,
			Color:      a.Color,
			Pretext:    a.Pretext,
			AuthorName: a.AuthorName,
			AuthorLink: a.AuthorLink,
			AuthorIcon: a.AuthorIcon,
			Title:      a.Title,
			TitleLink:  a.TitleLink,
			Text:       a.Text,
			ImageURL:   a.ImageURL,
			ThumbURL:   a.ThumbURL,
			Footer:     a.Footer,
			FooterIcon: a.FooterIcon,
		}

		if a.Actions != nil && len(a.Actions) > 0 {
			sa.Actions = []*mattermost.PostAction{}
			for _, act := range a.Actions {
				sAct := &mattermost.PostAction{
					Id:            act.Id,
					Type:          act.Type,
					Name:          act.Name,
					Disabled:      act.Disabled,
					Style:         act.Style,
					DataSource:    act.DataSource,
					Options:       []*mattermost.PostActionOptions{},
					DefaultOption: act.DefaultOption,
					Integration:   &mattermost.PostActionIntegration{},
					Cookie:        act.Cookie,
				}

				if act.Integration != nil {
					sAct.Integration.URL = act.Integration.URL
					sAct.Integration.Context = act.Integration.Context
				}

				if act.Options != nil && len(act.Options) > 0 {
					for _, o := range act.Options {
						sAct.Options = append(sAct.Options, &mattermost.PostActionOptions{
							Text:  o.Text,
							Value: o.Value,
						})
					}
				}

				sa.Actions = append(sa.Actions, sAct)
			}

		}

		repAttachments = append(repAttachments, sa)

	}

	return repAttachments
}
