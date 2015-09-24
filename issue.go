package main

import (
	"github.com/google/go-github/github"
	"strconv"
)

type Issue struct {
	ID       string `json:"issue"`
	Status   string `json:"status"`
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
	Assignee string `json:"assignee,omitempty"`
	Comments int    `json:"comments,omitempty"`
	URL      string `json:"url,omitempty"`
}

func (i *Issue) githubImport(gi *github.Issue) {
}
