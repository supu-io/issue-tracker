package main

import (
	"encoding/json"
	"log"
)

// Internal representation of an issue
type Issue struct {
	ID       string `json:"id"`
	Number   int    `json:"number"`
	Status   string `json:"status"`
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
	Assignee string `json:"assignee,omitempty"`
	Comments int    `json:"comments,omitempty"`
	URL      string `json:"url,omitempty"`
	Repo     string `json:"repo"`
	Owner    string `json:"rpo"`
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
	Config `json:"config"`
}
