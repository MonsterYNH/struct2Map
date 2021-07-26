package client

import (
	"fmt"
	"testing"
)

type QueryRequest struct {
	Query  string   `name:"query"`
	Page   int      `name:"page"`
	Page8  int8     `name:"page8"`
	Page16 int16    `name:"page16"`
	Page32 int32    `name:"page32"`
	Page64 int64    `name:"page64"`
	IsTrue bool     `name:"is_true"`
	Match  []string `name:"match"`
}

func TestRangeQuery(t *testing.T) {
	values, err := structToMap(&QueryRequest{
		Query:  "asdasd",
		IsTrue: false,
		Match:  []string{"1", "2"},
	})

	fmt.Println(values.Encode(), err)
}

type LabelsRequest struct {
	Query string   `name:"query"`
	Start string   `name:"start"`
	End   string   `name:"end"`
	Match []string `name:"match[]"`
}

func TestMakeRequest(t *testing.T) {
	data, err := MakeRequest("http://10.4.0.35:19192", "/api/v1/labels", "POST", "x-www-form-urlencoded", &LabelsRequest{
		// Match: []string{"go_gc_duration_seconds", "go_gc_duration_seconds_count"},
		Query: "go_gc_duration_seconds",
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(data))
}
