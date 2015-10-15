package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/supu-io/messages"
	"golang.org/x/oauth2"
)

// Github issue tracker
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
	// TODO : Move this hardcoded stuff to api
	t.Labels = []string{"created", "todo", "doing", "review", "uat", "done"}
}

// List of issues for a given status
func (t *Github) List(input *messages.GetIssuesList) *Issues {
	githubIssues := t.getIssuesList(input)

	issues := make(Issues, len(githubIssues))

	for i, issue := range githubIssues {
		issues[i] = t.mapIssue(&issue)
	}

	return &issues
}

func (t *Github) getIssuesList(input *messages.GetIssuesList) []github.Issue {

	if input.Repo == "" {
		options := github.IssueListOptions{
			Filter: "all",
		}
		if input.Status != "" {
			options.Labels = []string{input.Status}
		}
		options.Page = 0
		options.PerPage = 100
		githubIssues, _, err := t.client.Issues.ListByOrg(input.Org, &options)
		if err != nil {
			log.Println(err.Error())
		}
		return githubIssues
	}

	options := github.IssueListByRepoOptions{}
	if input.Status != "" {
		options.Labels = []string{input.Status}
	}
	options.Page = 0
	options.PerPage = 100
	githubIssues, _, err := t.client.Issues.ListByRepo(input.Org, input.Repo, &options)
	if err != nil {
		log.Println(err.Error())
	}
	return githubIssues
}

// Details for an issue for the given issue id
func (t *Github) Details(i *messages.Issue) *Issue {
	gIssue, _, _ := t.client.Issues.Get(i.Org, i.Repo, i.Number)
	opt := github.IssueListCommentsOptions{}
	gComments, _, _ := t.client.Issues.ListComments(i.Org, i.Repo, i.Number, &opt)

	issue := t.mapIssue(gIssue)
	if issue != nil {
		issue.Comments = t.mapComments(&gComments)
	}

	return issue
}

// Details for an issue for the given issue id
func (t *Github) Create(i *Issue) *Issue {
	// TODO default label must be provided
	ir := github.IssueRequest{
		Title:  &i.Title,
		Body:   &i.Body,
		Labels: &[]string{"created"},
	}

	gi, _, err := t.client.Issues.Create(i.Owner, i.Repo, &ir)
	if err != nil {
		log.Println(err)
		return nil
	}

	i.Number = *gi.Number
	i.URL = *gi.HTMLURL

	return i
}

// Update an issue by id
func (t *Github) Update(i *Issue) []string {
	for _, status := range t.Labels {
		_, err := t.client.Issues.RemoveLabelForIssue(i.Owner, i.Repo, i.Number, status)
		if err != nil {
			println(err.Error())
		}
	}
	labels, _, _ := t.client.Issues.AddLabelsToIssue(i.Owner, i.Repo, i.Number, []string{i.Status})
	ls := make([]string, len(labels))
	for index, label := range labels {
		ls[index] = *label.Name
	}
	return ls
}

// Setup github labels to support active workflow
func (t *Github) Setup(labels []string, owner string, repo string) {
	for _, status := range labels {
		l := github.Label{
			URL:  &status,
			Name: &status,
		}
		t.client.Issues.CreateLabel(owner, repo, &l)
	}
}

// Comment on an issue
func (t *Github) Comment(id string, body string) {
}

func (t *Github) mapIssue(gi *github.Issue) *Issue {
	i := Issue{}
	if gi == nil {
		return nil
	}
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
	if gi.URL != nil {
		repo := t.getRepoFromURL(*gi.URL)
		parts := strings.Split(repo, "/")
		i.Repo = parts[1]
		i.Owner = parts[0]
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
