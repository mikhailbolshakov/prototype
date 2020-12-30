package domain

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/users/repository/storage"
)

type UserSearchService interface {
	Search(criteria *SearchCriteria) (*SearchResponse, error)
}

func NewUserSearchService(storage storage.UserStorage) UserSearchService {
	return &searchImpl{storage: storage}
}

type searchImpl struct {
	storage storage.UserStorage
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
	return searchRsFromDto(r), nil
}