package storage

import (
	"context"
	"github.com/olivere/elastic"
	"gitlab.medzdrav.ru/prototype/kit/common"
	"gitlab.medzdrav.ru/prototype/kit/search"
	"gitlab.medzdrav.ru/prototype/tasks/domain"
	"math"
)

const (
	IDX_TASKS = "tasks"
)

type iTask struct {
	Id               string     `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Num              string     `json:"num"`
	Type             string     `json:"type"`
	SubType          string     `json:"subtype"`
	Status           string     `json:"status"`
	SubStatus        string     `json:"substatus"`
	AssigneeType     string     `json:"assigneeType"`
	AssigneeGroup    string     `json:"assigneeGroup"`
	AssigneeUserId   string     `json:"assigneeUserId"`
	AssigneeUsername string     `json:"assigneeUsername"`
	ChannelId        string     `json:"channelId"`
}

func (s *taskStorageImpl) ensureIndex() error {

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

	return s.c.Search.CreateIndexIfNotExists(IDX_TASKS, tasksMapping)
}

func (s *taskStorageImpl) Search(cr *domain.SearchCriteria) (*domain.SearchResponse, error) {

	response := &domain.SearchResponse{
		PagingResponse: &common.PagingResponse{
			Total: 0,
			Index: 0,
		},
		Tasks: []*domain.Task{},
	}

	cl := s.c.Search.GetClient()

	bq := elastic.NewBoolQuery()
	bq = bq.Must(elastic.NewMatchAllQuery())

	var queries []elastic.Query

	if cr.Num != "" {
		queries = append(queries, elastic.NewTermQuery("num", cr.Num))
	}

	if cr.Type != nil && cr.Type.Type != "" {
		queries = append(queries, elastic.NewTermQuery("type", cr.Type.Type))
	}

	if cr.Type != nil && cr.Type.SubType != "" {
		queries = append(queries, elastic.NewTermQuery("subtype", cr.Type.SubType))
	}

	if cr.Status != nil && cr.Status.Status != "" {
		queries = append(queries, elastic.NewTermQuery("status", cr.Status.Status))
	}

	if cr.Status != nil && cr.Status.SubStatus != "" {
		queries = append(queries, elastic.NewTermQuery("substatus", cr.Status.SubStatus))
	}

	if cr.Assignee != nil && cr.Assignee.UserId != "" {
		queries = append(queries, elastic.NewTermQuery("assigneeUserId", cr.Assignee.UserId))
	}

	if cr.Assignee != nil && cr.Assignee.Username != "" {
		queries = append(queries, elastic.NewTermQuery("assigneeUsername", cr.Assignee.Username))
	}

	if cr.Assignee != nil && cr.Assignee.Type != "" {
		queries = append(queries, elastic.NewTermQuery("assigneeType", cr.Assignee.Type))
	}

	if cr.Assignee != nil && cr.Assignee.Group != "" {
		queries = append(queries, elastic.NewTermQuery("assigneeGroup", cr.Assignee.Group))
	}

	if cr.ChannelId != "" {
		queries = append(queries, elastic.NewTermQuery("channelId", cr.ChannelId))
	}

	// paging
	from := (cr.Index - 1) * cr.Size
	if from < 0 {
		from = 0
	}

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

//func (s *taskStorageImpl) SearchDb(cr *searchCriteria) (*searchResponse, error) {
//
//	response := &searchResponse{
//		PagingResponse: &common.PagingResponse{
//			Total: 0,
//			Index: 0,
//		},
//		Tasks: []*task{},
//	}
//
//	selectClause := `*`
//
//	query := s.c.Db.Instance.
//		Table(`tasks t`).
//		Where(`t.deleted_at is null`)
//
//	if cr.Num != "" {
//		query = query.Where(`t.num = ?`, cr.Num)
//	}
//
//	if cr.Type != "" {
//		query = query.Where(`t.type = ?`, cr.Type)
//	}
//
//	if cr.SubType != "" {
//		query = query.Where(`t.subtype = ?`, cr.SubType)
//	}
//
//	if cr.Status != "" {
//		query = query.Where(`t.status = ?`, cr.Status)
//	}
//
//	if cr.SubStatus != "" {
//		query = query.Where(`t.substatus = ?`, cr.SubStatus)
//	}
//
//	if cr.AssigneeUserId != "" {
//		query = query.Where(`t.assignee_user_id = ?`, cr.AssigneeUserId)
//	}
//
//	if cr.AssigneeUsername != "" {
//		query = query.Where(`t.assignee_username = ?`, cr.AssigneeUsername)
//	}
//
//	if cr.AssigneeType != "" {
//		query = query.Where(`t.assignee_type = ?`, cr.AssigneeType)
//	}
//
//	if cr.AssigneeGroup != "" {
//		query = query.Where(`t.assignee_group = ?`, cr.AssigneeGroup)
//	}
//
//	if cr.ChannelId != "" {
//		query = query.Where(`t.channel_id = ?`, cr.ChannelId)
//	}
//
//	// paging
//	var totalCount int64
//	var offset int
//
//	query.Count(&totalCount)
//
//	if totalCount > int64(cr.Size) {
//		offset = (cr.Index - 1) * cr.Size
//	}
//
//	response.PagingResponse.Total = int(math.Ceil(float64(totalCount) / float64(cr.Size)))
//	response.PagingResponse.Index = cr.Index
//
//	query = query.Select(selectClause).Offset(offset).Limit(cr.Size)
//
//	rows, err := query.Rows()
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		task := &task{}
//		_ = s.c.Db.Instance.ScanRows(rows, task)
//		response.Tasks = append(response.Tasks, task)
//	}
//
//	return response, nil
//}
