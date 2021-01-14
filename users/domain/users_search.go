package domain

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	pb "gitlab.medzdrav.ru/prototype/proto/mm"
	"gitlab.medzdrav.ru/prototype/users/repository/adapters/mattermost"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
)

type UserSearchService interface {
	Search(criteria *SearchCriteria) (*SearchResponse, error)
}

func NewUserSearchService(storage storage.UserStorage, mmService mattermost.Service) UserSearchService {
	return &searchImpl{
		storage: storage,
		mmService: mmService,
	}
}

type searchImpl struct {
	storage storage.UserStorage
	mmService mattermost.Service
}

func (s *searchImpl) Search(cr *SearchCriteria) (*SearchResponse, error) {
	if cr.PagingRequest == nil {
		cr.PagingRequest = &common.PagingRequest{}
	}

	if cr.Size == 0 {
		cr.Size = 100
	}

	if cr.Index == 0 {
		cr.Index = 1
	}

	r, err := s.storage.Search(criteriaToDto(cr))
	if err != nil {
		return nil, err
	}
	response := searchRsFromDto(r)

	if len(r.Users) > 0 && cr.OnlineStatuses != nil && len(cr.OnlineStatuses) > 0 {
		var mmUserIds []string
		modifiedResponse := &SearchResponse{Users: []*User{}, PagingResponse: &common.PagingResponse{
			Total: response.Total,
			Index: response.Index,
		}}

		for _, u := range r.Users {
			mmUserIds = append(mmUserIds, u.MMUserId)
		}

		if mmStatuses, err := s.mmService.GetUsersStatuses(&pb.GetUsersStatusesRequest{MMUserIds: mmUserIds}); err == nil {

			for _, user := range response.Users {

				for _, mmSt := range mmStatuses.Statuses {

					if mmSt.MMUserId == user.MMUserId {

						for _, criteriaStatus := range cr.OnlineStatuses {

							if criteriaStatus == mmSt.Status {
								modifiedResponse.Users = append(modifiedResponse.Users, user)
							}

						}

					}

				}

			}

			if response.Total < cr.Size {
				modifiedResponse.Total = len(modifiedResponse.Users)
			}

			return modifiedResponse, nil

		} else {
			return nil, err
		}
	}


	return response, nil
}