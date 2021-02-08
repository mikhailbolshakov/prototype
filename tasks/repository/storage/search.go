package storage

import (
	"context"
	"github.com/olivere/elastic"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/search"
	"math"
)

const (
	IDX_TASKS = "tasks"
)

func (s *taskStorageImpl) EnsureIndex() error {

	tasksMapping := &search.Mapping{FieldsMapping: map[string]*search.FieldMapping{
		"id":               {Type: search.T_KEYWORD},
		"num":              {Type: search.T_KEYWORD},
		"type":             {Type: search.T_KEYWORD},
		"subtype":          {Type: search.T_KEYWORD},
		"status":           {Type: search.T_KEYWORD},
		"substatus":        {Type: search.T_KEYWORD},
		"assigneeUserId":   {Type: search.T_KEYWORD},
		"assigneeUsername": {Type: search.T_KEYWORD},
		"assigneeGroup":    {Type: search.T_KEYWORD},
		"assigneeType":     {Type: search.T_KEYWORD},
		"channelId":        {Type: search.T_KEYWORD},
		"title":            {Type: search.T_TEXT},
		"description":      {Type: search.T_TEXT},
		"deletedAt":        {Type: search.T_DATE},
	}}

	return s.infr.Search.CreateIndexIfNotExists(IDX_TASKS, tasksMapping)
}

func (t *Task) toIndex() *iTask {
	return &iTask{
		Id:               t.Id,
		Title:            t.Title,
		Description:      t.Description,
		Num:              t.Num,
		Type:             t.Type,
		SubType:          t.SubType,
		Status:           t.Status,
		SubStatus:        t.SubStatus,
		AssigneeType:     t.AssigneeType,
		AssigneeGroup:    t.AssigneeGroup,
		AssigneeUserId:   t.AssigneeUserId,
		AssigneeUsername: t.AssigneeUsername,
		ChannelId:        t.ChannelId,
	}
}

func (s *taskStorageImpl) Search(cr *SearchCriteria) (*SearchResponse, error) {

	response := &SearchResponse{
		PagingResponse: &common.PagingResponse{
			Total: 0,
			Index: 0,
		},
		Tasks: []*Task{},
	}

	cl := s.infr.Search.GetClient()

	bq := elastic.NewBoolQuery()
	bq = bq.Must(elastic.NewMatchAllQuery())

	var queries []elastic.Query

	if cr.Num != "" {
		queries = append(queries, elastic.NewTermQuery("num", cr.Num))
	}

	if cr.Type != "" {
		queries = append(queries, elastic.NewTermQuery("type", cr.Type))
	}

	if cr.SubType != "" {
		queries = append(queries, elastic.NewTermQuery("subtype", cr.SubType))
	}

	if cr.Status != "" {
		queries = append(queries, elastic.NewTermQuery("status", cr.Status))
	}

	if cr.SubStatus != "" {
		queries = append(queries, elastic.NewTermQuery("substatus", cr.SubStatus))
	}

	if cr.AssigneeUserId != "" {
		queries = append(queries, elastic.NewTermQuery("assigneeUserId", cr.AssigneeUserId))
	}

	if cr.AssigneeUsername != "" {
		queries = append(queries, elastic.NewTermQuery("assigneeUsername", cr.AssigneeUsername))
	}

	if cr.AssigneeType != "" {
		queries = append(queries, elastic.NewTermQuery("assigneeType", cr.AssigneeType))
	}

	if cr.AssigneeGroup != "" {
		queries = append(queries, elastic.NewTermQuery("assigneeGroup", cr.AssigneeGroup))
	}

	if cr.ChannelId != "" {
		queries = append(queries, elastic.NewTermQuery("channelId", cr.ChannelId))
	}

	// paging
	from := (cr.Index - 1) * cr.Size

	bq = bq.Filter(queries...)
	sr, err := cl.Search(IDX_TASKS).
		Query(bq).
		From(from).
		Size(cr.Size).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var ids []string

	if sr.TotalHits() > 0 {

		for _, sh := range sr.Hits.Hits {
			response.Tasks = append(response.Tasks, &Task{Id: sh.Id})
			ids = append(ids, sh.Id)
		}

	}

	response.PagingResponse.Total = int(math.Ceil(float64(sr.TotalHits()) / float64(cr.Size)))
	response.PagingResponse.Index = cr.Index

	if len(ids) > 0 {
		response.Tasks = s.GetByIds(ids)
	}

	return response, nil
}

func (s *taskStorageImpl) SearchDb(cr *SearchCriteria) (*SearchResponse, error) {

	response := &SearchResponse{
		PagingResponse: &common.PagingResponse{
			Total: 0,
			Index: 0,
		},
		Tasks: []*Task{},
	}

	selectClause := `*`

	query := s.infr.Db.Instance.
		Table(`tasks t`).
		Where(`t.deleted_at is null`)

	if cr.Num != "" {
		query = query.Where(`t.num = ?`, cr.Num)
	}

	if cr.Type != "" {
		query = query.Where(`t.type = ?`, cr.Type)
	}

	if cr.SubType != "" {
		query = query.Where(`t.subtype = ?`, cr.SubType)
	}

	if cr.Status != "" {
		query = query.Where(`t.status = ?`, cr.Status)
	}

	if cr.SubStatus != "" {
		query = query.Where(`t.substatus = ?`, cr.SubStatus)
	}

	if cr.AssigneeUserId != "" {
		query = query.Where(`t.assignee_user_id = ?`, cr.AssigneeUserId)
	}

	if cr.AssigneeUsername != "" {
		query = query.Where(`t.assignee_username = ?`, cr.AssigneeUsername)
	}

	if cr.AssigneeType != "" {
		query = query.Where(`t.assignee_type = ?`, cr.AssigneeType)
	}

	if cr.AssigneeGroup != "" {
		query = query.Where(`t.assignee_group = ?`, cr.AssigneeGroup)
	}

	if cr.ChannelId != "" {
		query = query.Where(`t.channel_id = ?`, cr.ChannelId)
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
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := &Task{}
		_ = s.infr.Db.Instance.ScanRows(rows, task)
		response.Tasks = append(response.Tasks, task)
	}

	return response, nil
}
