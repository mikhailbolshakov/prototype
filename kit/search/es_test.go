package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"html/template"
	"reflect"
	"testing"
)

type Tweet struct {
	User     string
	Message  string
	Retweets int64
}

func Test_Simple(t *testing.T) {

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient()
	// Trace request and response details like this
	// client, err := elastic.NewClient(elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	if err != nil {
		// Handle error
		t.Fatal(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("twitter").Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}

	if exists {
		// Delete an index.
		deleteIndex, err := client.DeleteIndex("twitter").Do(context.Background())
		if err != nil {
			// Handle error
			t.Fatal(err)
		}
		if !deleteIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// Create a new index.
	mapping := `
{
	"mappings":{
		"properties":{
			"user":{
				"type":"keyword"
			},
			"message":{
				"type":"text"
			},
			"retweets":{
				"type":"long"
			}
		}
	}
}
`
	createIndex, err := client.CreateIndex("twitter").Body(mapping).Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
	}

	// Index a tweet (using JSON serialization)
	tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index("twitter").
		Id("1").
		BodyJson(tweet1).
		Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Index a second tweet (by string)
	tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
	put2, err := client.Index().
		Index("twitter").
		Id("2").
		BodyString(tweet2).
		Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put2.Id, put2.Index, put2.Type)

	// Get tweet with specified ID
	get1, err := client.Get().
		Index("twitter").
		Id("1").
		Do(context.Background())
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			t.Fatal(fmt.Sprintf("Document not found: %v", err))
		case elastic.IsTimeout(err):
			t.Fatal(fmt.Sprintf("Timeout retrieving document: %v", err))
		case elastic.IsConnErr(err):
			t.Fatal(fmt.Sprintf("Connection problem: %v", err))
		default:
			// Some other kind of error
			t.Fatal(err)
		}
	}
	fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)

	// Refresh to make sure the documents are searchable.
	_, err = client.Refresh().Index("twitter").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Search with a term query
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter"). // search in index "twitter"
		Query(termQuery). // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10). // take documents 0-9
		Pretty(true). // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		// Handle error
		t.Fatal(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(Tweet)
		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}

	// Update a tweet by the update API of Elasticsearch.
	// We just increment the number of retweets.
	script := elastic.NewScript("ctx._source.retweets += params.num").Param("num", 1)
	update, err := client.Update().Index("twitter").Id("1").
		Script(script).
		Upsert(map[string]interface{}{"retweets": 0}).
		Do(context.Background())
	if err != nil {
		// Handle error
		t.Fatal(err)
	}
	fmt.Printf("New version of tweet %q is now %d", update.Id, update.Version)

	// ...

}

func Test_MappingTemplate(t *testing.T) {
	mapping :=
		`
{{$l := isLast}}
{
"mappings": {
	"properties": {
		{{range $key, $value := .FieldsMapping}}
			"{{$key}}": {
							"type":"{{$value.Type}}"
						}{{if not (call $l)}},{{else}}{{end}}
		{{end}}
	}
}
}`

	m := &Mapping{FieldsMapping: map[string]*FieldMapping{
		"username":    {Type: T_KEYWORD},
		"description": {Type: T_TEXT},
	}}

	isLast := func() func() bool {
		i := 0
		ln := len(m.FieldsMapping)
		return func() bool {
			i++
			return ln == i
		}
	}

	tmpl := template.Must(template.New("M").Funcs(template.FuncMap{
		"isLast": isLast,
	}).Parse(mapping))

	var body bytes.Buffer
	err := tmpl.ExecuteTemplate(&body, "M", m)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(body.String())

}
