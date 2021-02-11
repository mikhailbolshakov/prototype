package storage

import (
	"context"
	"fmt"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/users/domain"
	"math"
)

func (s *storageImpl) ensureIndex() error {
	return nil
}

// TODO: ES
func (s *storageImpl) Search(ctx context.Context, cr *domain.SearchCriteria) (*domain.SearchResponse, error) {

	response := &domain.SearchResponse{
		PagingResponse: &common.PagingResponse{
			Total: 0,
			Index: 0,
		},
		Users: []*domain.User{},
	}

	selectClause := `*`

	query := s.c.Db.Instance.
		Table(`users u`).
		Where(`u.deleted_at is null`)

	if cr.Username != "" {
		query = query.Where(`u.username = ?`, cr.Username)
	}

	if cr.UserGroup != "" {
		query = query.Where(`exists(select 1 from users.user_groups ug where ug.group_code = ? and ug.user_id = u.id and ug.deleted_at is null)`, cr.UserGroup)
	}

	if cr.Status != "" {
		query = query.Where(`u.status = ?`, cr.Status)
	}

	if cr.CommonChannelId != "" {
		query = query.Where(`(u.details -> 'commonChannelId')::varchar = ?`, fmt.Sprintf(`"%s"`, cr.CommonChannelId))
	}

	if cr.MedChannelId != "" {
		query = query.Where(`(u.details -> 'medChannelId')::varchar = ?`, fmt.Sprintf(`"%s"`, cr.MedChannelId))
	}

	if cr.LawChannelId != "" {
		query = query.Where(`(u.details -> 'lawChannelId')::varchar = ?`, fmt.Sprintf(`"%s"`, cr.LawChannelId))
	}

	if cr.MMId != "" {
		query = query.Where(`u.mm_id = ?`, cr.MMId)
	}

	if cr.Email != "" {
		query = query.Where(`(u.details -> 'email')::varchar = ?`, fmt.Sprintf(`"%s"`, cr.Email))
	}

	if cr.Phone != "" {
		query = query.Where(`(u.details -> 'phone')::varchar = ?`, fmt.Sprintf(`"%s"`, cr.Phone))
	}

	if cr.UserType != "" {
		query = query.Where(`u.type = ?`, cr.UserType)
	}

	// paging
	var totalCount int64
	var offset int

	query.Count(&totalCount)

	if totalCount > int64(cr.Size) {
		offset = (cr.Index - 1) * cr.Size
	}

	response.PagingResponse.Total = int(math.Ceil(float64(totalCount) / float64(cr.Size)))
	response.PagingResponse.Index = cr.Index

	query = query.Select(selectClause).Offset(offset).Limit(cr.Size)

	rows, err := query.Rows()
	var users []*user
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := &user{Groups: []*userGroup{}}
		_ = s.c.Db.Instance.ScanRows(rows, user)
		users= append(users, user)
	}
	response.Users = s.toUsersDomain(users)

	return response, nil
}
