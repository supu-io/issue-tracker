package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

// Internal representation of an issue
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
type Comments []Comment

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

// A collection of issues
type Issues []*Issue

// Get the json representation for a collection of issues
func (i *Issues) toJSON() *[]byte {
	json, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
	}
	return &json
}

// This config needs to be received on any input
// event
type Config struct {
	Github *Github `json:"github, omitempty"`
}

// Representation for the input event issues.list
type IssuesList struct {
	Status string `json:"status"`
	Org    string `json:"org"`
	Config `json:"config"`
}

type IssuesDetails struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Config `json:"config"`
}

func (i *IssuesDetails) toIssue() *Issue {
	issue := Issue{}
	parts := strings.Split(i.ID, "/")
	issue.Owner = parts[1]
	issue.Repo = parts[2]
	issue.Status = i.Status
	number, _ := strconv.Atoi(parts[3])
	issue.Number = number

	return &issue
}
