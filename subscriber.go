package main

import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/nats-io/nats"
	"github.com/supu-io/messages"
)

// Subscriber : collection of methods to subscribe the different events
type Subscriber struct{}

// Subscribe to all events
func (s *Subscriber) Subscribe() {

	log.Println("Listening...")

	nc, _ := nats.Connect(nats.DefaultURL)
	nc.Subscribe("issues.list", func(m *nats.Msg) {
		issues := s.issuesList(m.Data)
		nc.Publish(m.Reply, *issues.toJSON())
	})

	nc.Subscribe("issues.details", func(m *nats.Msg) {
		issue := s.issuesDetails(m.Data)
		if issue != nil {
			nc.Publish(m.Reply, *issue.toJSON())
		} else {
			nc.Publish(m.Reply, []byte(`{"error":"non existing issue"}`))
			return
		}
	})

	nc.Subscribe("issues.update", func(m *nats.Msg) {
		issues := s.issuesUpdate(m.Data)
		nc.Publish(m.Reply, *issues.toJSON())
	})

	nc.Subscribe("issues.create", func(m *nats.Msg) {
		res := s.issuesCreate(m.Data)
		nc.Publish(m.Reply, *res.toJSON())
	})

	nc.Subscribe("issue-tracker.setup", func(m *nats.Msg) {
		s.setup(nc, m)
	})
	runtime.Goexit()
}

// Manages an issues.list event
func (s *Subscriber) issuesList(body []byte) Issues {
	input := IssuesList{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
	}

	g := input.Config.Github
	g.setup()
	issues := g.List(&input)

	return *issues
}

func (s *Subscriber) issuesDetails(body []byte) *Issue {
	input := messages.GetIssue{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
	}

	g := Github{
		Token: input.Config.Github.Token,
	}
	g.setup()
	issue := input.Issue

	return g.Details(&issue)
}

func (s *Subscriber) issuesCreate(body []byte) *Issue {
	input := IssuesCreate{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
	}

	g := input.Config.Github
	g.setup()
	issue := input.Issue

	return g.Create(issue)
}

func (s *Subscriber) issuesUpdate(body []byte) *Issue {
	input := IssuesDetails{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
	}

	g := input.Config.Github
	g.setup()
	issue := input.toIssue()
	g.Update(issue)

	return issue
}

func (s *Subscriber) setup(nc *nats.Conn, m *nats.Msg) {
	type SetupMsg struct {
		Org    string   `json:"org"`
		Repo   string   `json:"repo"`
		Labels []string `json:"states"`
		Config `json:"config"`
	}
	body := m.Data
	input := SetupMsg{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		nc, _ := nats.Connect(nats.DefaultURL)
		nc.Publish(m.Reply, []byte(`{"error":"error setting up github labels"}`))
		return
	}
	g := input.Config.Github
	g.setup()
	g.Labels = input.Labels
	g.Setup(input.Labels, input.Org, input.Repo)
	nc.Publish(m.Reply, []byte(`success`))
}
