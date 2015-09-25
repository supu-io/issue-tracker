package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestIssuesService_ListByOrg(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"number":1,"title":"test","body":"body","comments":0,"html_url":"http://supu.io","url":"https://api.github.com/repos/octocat/Hello-World/pulls/1347"}]`)
	})

	g := Github{client: client, Org: "o"}
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
	if issue.URL != "http://supu.io" {
		t.Errorf("Issue url is not successfully mapped")
	}
	if issue.Repo != "octocat/Hello-World" {
		t.Errorf("Issue repo is not successfully mapped")
	}
}

func TestIssuesService_Update(t *testing.T) {
	setup()
	defer teardown()

	input := []string{"doing"}

	mux.HandleFunc("/repos/o/r/issues/1/labels/todo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels/doing", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels/review", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels/uat", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels/done", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})
	mux.HandleFunc("/repos/o/r/issues/1/labels", func(w http.ResponseWriter, r *http.Request) {
		v := new([]string)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(*v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `[{"name":"doing"}]`)
	})

	g := Github{client: client, Org: "o"}
	issue := Issue{Owner: "o", Repo: "r", Number: 1, Status: "doing"}
	labels := g.Update(&issue)

	if len(labels) != 1 {
		t.Errorf("Invalid number of labels %+v, want %+v", len(labels), 1)
	}
	label := labels[0]
	if label != "doing" {
		t.Errorf("Label is not successfully mapped")
	}

}
