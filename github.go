package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"strconv"
)

// Issue tracker for github
type Github struct {
	client *github.Client
}

// Setup the client for github
func (t *Github) setup() {
	t.client = github.NewClient(nil)
}

// Setup the client for github with an oauth2 token
func (t *Github) auth(token string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	t.client = github.NewClient(tc)
}

// Get a list of issues for a given status
func (t *Github) List(status string) *[]*Issue {
	s := []string{status}
	options := github.IssueListOptions{Labels: s}
	githubIssues, _, err := t.client.Issues.ListByOrg("supu-io", &options)
	if err != nil {
		log.Println(err.Error())
	}

	issues := []*Issue{}
	for i, issue := range githubIssues {
		issues[i] = t.exportIssue(&issue)
	}

	return &issues
}

// Gets issue details for the fiven issue id
func (t *Github) Details(id string) *Issue {
	return &Issue{}
}

// Updates an issue by id
func (t *Github) Update(id string, issue *Issue) {
}

// Adds a comment on an issue
func (t *Github) Comment(id string, body string) {
}

func (t *Github) exportIssue(gi *github.Issue) *Issue {
	i := Issue{}

	i.ID = strconv.Itoa(*gi.Number)
	i.Title = *gi.Title
	i.Body = *gi.Body
	i.Assignee = *gi.Assignee.Name
	i.Comments = *gi.Comments
	i.URL = *gi.URL

	statuses := []string{"todo", "doing", "review", "uat", "done"}

	for _, label := range gi.Labels {
		for _, status := range statuses {
			if *label.Name == status {
				i.Status = status
			}
		}
	}

	return &i
}
