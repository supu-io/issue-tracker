package main

type Issue struct {
	ID       string `json:"issue"`
	Status   string `json:"status"`
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
	Assignee string `json:"assignee,omitempty"`
	Comments int    `json:"comments,omitempty"`
	URL      string `json:"url,omitempty"`
}
