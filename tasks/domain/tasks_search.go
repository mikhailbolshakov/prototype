package domain

import (
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/tasks/repository/storage"
)

type TaskSearchService interface {
	// search tasks by the given criteria
	Search(cr *SearchCriteria) (*SearchResponse, error)
}

type searchServiceImpl struct {
	storage storage.TaskStorage
}

func NewTaskSearchService(storage storage.TaskStorage) TaskSearchService {
	return &searchServiceImpl{
		storage: storage,
	}
}

func (s *searchServiceImpl) Search(cr *SearchCriteria) (*SearchResponse, error) {

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
