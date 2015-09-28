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
func (t *Github) List(input *IssuesList) *Issues {
	options := github.IssueListOptions{}
	options.Filter = "all"
	options.Labels = []string{input.Status}
	githubIssues, _, err := t.client.Issues.ListByOrg(input.Org, &options)
	if err != nil {
		log.Println(err.Error())
	}

	issues := make(Issues, len(githubIssues))
	for i, issue := range githubIssues {
		issues[i] = t.mapIssue(&issue)
	}

	return &issues
}

// Gets issue details for the given issue id
func (t *Github) Details(i *Issue) *Issue {
	gIssue, _, _ := t.client.Issues.Get(i.Owner, i.Repo, i.Number)
	opt := github.IssueListCommentsOptions{}
	gComments, _, _ := t.client.Issues.ListComments(i.Owner, i.Repo, i.Number, &opt)

	issue := t.mapIssue(gIssue)
	issue.Comments = t.mapComments(&gComments)

	return issue
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

func (t *Github) mapIssue(gi *github.Issue) *Issue {
	i := Issue{}
	if gi.Number != nil {
		i.ID = strconv.Itoa(*gi.Number)
	}
	if gi.Number != nil {
		i.Number = *gi.Number
	}
	if gi.Title != nil {
		i.Title = *gi.Title
	}
	if gi.Body != nil {
		i.Body = *gi.Body
	}

	// i.Assignee = *gi.Assignee.Name
	if gi.URL != nil {
		i.Repo = t.getRepoFromURL(*gi.URL)
	}
	if gi.HTMLURL != nil {
		i.URL = *gi.HTMLURL
	}

	// TODO: Map extra fields

	for _, label := range gi.Labels {
		for _, status := range t.Labels {
			if *label.Name == status {
				i.Status = status
			}
		}
	}

	return &i
}

func (t *Github) mapComments(gComments *[]github.IssueComment) []Comment {
	comments := make([]Comment, len(*gComments))

	for index, gc := range *gComments {
		if gc.ID != nil {
			comments[index].ID = *gc.ID
		}
		if gc.Body != nil {
			comments[index].Body = *gc.Body
		}
		if gc.User != nil {
			comments[index].User = *gc.User.Name
		}
		if gc.URL != nil {
			comments[index].URL = *gc.URL
		}
		if gc.CreatedAt != nil {
			comments[index].CreatedAt = *gc.CreatedAt
		}
		if gc.UpdatedAt != nil {
			comments[index].UpdatedAt = *gc.UpdatedAt
		}
		if gc.HTMLURL != nil {
			comments[index].HTMLURL = *gc.HTMLURL
		}
		if gc.IssueURL != nil {
			comments[index].IssueURL = *gc.IssueURL
		}
	}
	return comments
}

func (t *Github) getRepoFromURL(url string) string {
	// "https://api.github.com/repos/octocat/Hello-World/pulls/1347"
	parts := strings.Split(url, "/")
	return parts[4] + "/" + parts[5]
}
