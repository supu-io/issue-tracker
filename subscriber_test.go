package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestIssuesService_ListByOrg(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"number":1,"title":"test","body":"body","comments":0,"html_url":"http://supu.io"}]`)
	})

	g := Github{client: client, Org: "o"}
	// issues := s.issuesList(m)
	issues := *g.List("todo")
	if len(issues) != 1 {
		t.Errorf("Issues.List returned %+v, want %+v", len(issues), 1)
	}
	issue := issues[0]
	if issue.Title != "test" {
		t.Errorf("Issue title is not successfully mapped")
	}
	if issue.ID != "1" {
		t.Errorf("Issue number is not successfully mapped")
	}
	if issue.Body != "body" {
		t.Errorf("Issue body is not successfully mapped")
	}
	if issue.Comments != 0 {
		t.Errorf("Issue comments is not successfully mapped")
	}
	if issue.URL != "http://supu.io" {
		t.Errorf("Issue url is not successfully mapped")
	}
}
