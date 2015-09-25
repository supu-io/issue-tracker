package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"strconv"
	"strings"
)

// Issue tracker for github
type Github struct {
	Token  string `json:"token"`
	Org    string `json:"org"`
	Labels []string
	client *github.Client
}

// Setup the client for github
func (t *Github) setup() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: t.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	t.client = github.NewClient(tc)
	t.Labels = []string{"todo", "doing", "review", "uat", "done"}
}

// Get a list of issues for a given status
func (t *Github) List(status string) *Issues {
	options := github.IssueListOptions{}
	options.Filter = "all"
	options.Labels = []string{status}
	githubIssues, _, err := t.client.Issues.ListByOrg(t.Org, &options)
	if err != nil {
		log.Println(err.Error())
	}

	issues := make(Issues, len(githubIssues))
	for i, issue := range githubIssues {
		issues[i] = t.exportIssue(&issue)
	}

	return &issues
}

// Gets issue details for the given issue id
func (t *Github) Details(id string) *Issue {
	// func (s *IssuesService) Get(owner string, repo string, number int) (*Issue, *Response, error)
	return &Issue{}
}

// Updates an issue by id
func (t *Github) Update(i *Issue) []string {
	for _, status := range t.Labels {
		t.client.Issues.RemoveLabelForIssue(i.Owner, i.Repo, i.Number, status)
	}
	labels, _, _ := t.client.Issues.AddLabelsToIssue(i.Owner, i.Repo, i.Number, []string{i.Status})
	ls := make([]string, len(labels))
	for index, label := range labels {
		ls[index] = *label.Name
	}
	return ls
}

// Adds a comment on an issue
func (t *Github) Comment(id string, body string) {
}

func (t *Github) exportIssue(gi *github.Issue) *Issue {
	i := Issue{}
	i.ID = strconv.Itoa(*gi.Number)
	i.Number = *gi.Number
	i.Title = *gi.Title
	i.Body = *gi.Body
	// i.Assignee = *gi.Assignee.Name
	i.Repo = t.getRepoFromURL(*gi.URL)
	i.Comments = *gi.Comments
	i.URL = *gi.HTMLURL

	for _, label := range gi.Labels {
		for _, status := range t.Labels {
			if *label.Name == status {
				i.Status = status
			}
		}
	}

	return &i
}

func (t *Github) getRepoFromURL(url string) string {
	// "https://api.github.com/repos/octocat/Hello-World/pulls/1347"
	parts := strings.Split(url, "/")
	return parts[4] + "/" + parts[5]
}
