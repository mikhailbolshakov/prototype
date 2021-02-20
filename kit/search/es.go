package search

import (
	"bytes"
	"context"
	"github.com/olivere/elastic/v7"
	"gitlab.medzdrav.ru/prototype/kit"
	"gitlab.medzdrav.ru/prototype/kit/log"
	"text/template"
)

type FieldMapping struct {
	Type string
}

type Mapping struct {
	FieldsMapping map[string]*FieldMapping
}

type Search interface {
	CreateIndexIfNotExists(index string, mapping *Mapping) error
	Index(index string, id string, data interface{}) error
	IndexAsync(index string, id string, data interface{})
	GetClient() *elastic.Client
	Close()
}

type esImpl struct {
	client *elastic.Client
}

func NewEs(url string, trace bool) (Search, error) {

	l := log.L().Cmp("es").Mth("new").F(log.FF{"url": url})

	s := &esImpl{}

	opts := []elastic.ClientOptionFunc{elastic.SetURL(url)}
	if trace {
		opts = append(opts, elastic.SetTraceLog(log.GetLogger()))
	}

	cl, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	s.client = cl
	l.Inf("ok")
	return s, nil
}

const (
	T_KEYWORD = "keyword"
	T_TEXT    = "text"
	T_DATE    = "date"
)

func (s *esImpl) createIndex(index string, mapping *Mapping) error {

	mappingTmpl :=
		`
{{$l := isLast}}
{"mappings": {
	"properties": {
		{{range $key, $value := .FieldsMapping}}
			"{{$key}}": {"type":"{{$value.Type}}"}{{if not (call $l)}},{{else}}{{end}}
		{{end}}
	}
}}`

	l := log.L().Cmp("es").Mth("create-index")

	// here is a trick with closure to put commas correctly (avoid comma after the last item)
	isLast := func() func() bool {
		i := 0
		ln := len(mapping.FieldsMapping)
		return func() bool {
			i++
			return ln == i
		}
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"isLast": isLast,
	}).Parse(mappingTmpl)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, mapping)
	if err != nil {
		return err
	}

	_, err = s.client.CreateIndex(index).BodyString(body.String()).Do(context.Background())
	if err != nil {
		return err
	}

	l.Dbg("index %s created", index)

	return nil
}

func (s *esImpl) CreateIndexIfNotExists(index string, mapping *Mapping) error {

	l := log.L().Cmp("es").Mth("create-index")

	exists, err := s.client.IndexExists(index).Do(context.Background())
	if err != nil {
		return err
	}

	if exists {
		l.DbgF("index %s exists", index)
		return nil
	} else {
		return s.createIndex(index, mapping)
	}

}

func (s *esImpl) Index(index string, id string, data interface{}) error {

	log.L().Cmp("es").Mth("indexation").F(log.FF{"index": index, "id": id}).Dbg().Trc(kit.Json(data))

	_, err := s.client.Index().
		Index(index).
		Id(id).
		BodyJson(data).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *esImpl) IndexAsync(index string, id string, data interface{}) {

	go func() {

		l := log.L().Cmp("es").Mth("indexation").F(log.FF{"index": index, "id": id}).Dbg().Trc(kit.Json(data))

		_, err := s.client.Index().
			Index(index).
			Id(id).
			BodyJson(data).
			// don't refresh immediately, so index will be available for search in some point in the future
			Refresh("false").
			Do(context.Background())
		if err != nil {
			l.E(err).Err()
		}

	}()

}

func (s *esImpl) GetClient() *elastic.Client {
	return s.client
}


func (s *esImpl) Close() {
	s.client.Stop()
}
