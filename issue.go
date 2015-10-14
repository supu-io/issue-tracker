package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

// Issue Internal representation
type Issue struct {
	ID       string    `json:"id"`
	Number   int       `json:"number"`
	Status   string    `json:"status"`
	Title    string    `json:"title,omitempty"`
	Body     string    `json:"body,omitempty"`
	Assignee string    `json:"assignee,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
	URL      string    `json:"url,omitempty"`
	Repo     string    `json:"repo"`
	Owner    string    `json:"owner"`
}

// Comments of an issue
type Comments []Comment

// Comment ...
type Comment struct {
	ID        int       `json:"id,omitempty"`
	Body      string    `json:"body,omitempty"`
	User      string    `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	URL       string    `json:"url,omitempty"`
	HTMLURL   string    `json:"html_url,omitempty"`
	IssueURL  string    `json:"issue_url,omitempty"`
}

// Get the json representation for an issue
func (i *Issue) toJSON() *[]byte {
	json, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
	}
	return &json
}

// Issues collection
type Issues []*Issue

// Get the json representation for a collection of issues
func (i *Issues) toJSON() *[]byte {
	json, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
	}
	return &json
}

// Config ...
type Config struct {
	Github *Github `json:"github, omitempty"`
}

// IssuesList Representation
type IssuesList struct {
	Status string `json:"status"`
	Org    string `json:"org"`
	Repo   string `json:"repo,omitempty"`
	Config `json:"config"`
}

// IssuesDetails ...
type IssuesDetails struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Config `json:"config"`
}

func (i *IssuesDetails) toIssue() *Issue {
	issue := Issue{}
	if strings.Contains(i.ID, "/") == false {
		return nil
	}
	parts := strings.Split(i.ID, "/")
	issue.Owner = parts[0]
	issue.Repo = parts[1]
	issue.Status = i.Status
	number, _ := strconv.Atoi(parts[2])
	issue.Number = number

	return &issue
}

type IssuesCreate struct {
	Issue  *Issue `json:issue`
	Config `json:"config"`
}
