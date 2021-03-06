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
		res := s.issuesUpdate(m.Data)
		nc.Publish(m.Reply, *ToJSON(res))
	})

	nc.Subscribe("issues.create", func(m *nats.Msg) {
		res := s.issuesCreate(m.Data)
		nc.Publish(m.Reply, *ToJSON(res))
	})

	nc.Subscribe("issue-tracker.setup", func(m *nats.Msg) {
		s.setup(nc, m)
	})
	runtime.Goexit()
}

// Manages an issues.list event
func (s *Subscriber) issuesList(body []byte) Issues {
	input := messages.GetIssuesList{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
	}

	g := getAdapter(input.Config)
	issues := g.List(&input)

	return *issues
}

func (s *Subscriber) issuesDetails(body []byte) *Issue {
	input := messages.GetIssue{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
	}

	g := getAdapter(input.Config)
	issue := input.Issue

	return g.Details(issue)
}

func (s *Subscriber) issuesCreate(body []byte) *messages.Issue {
	input := messages.CreateIssue{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
	}

	g := getAdapter(input.Config)
	issue := input.Issue

	return g.Create(issue)
}

func (s *Subscriber) issuesUpdate(body []byte) *messages.Issue {
	input := messages.UpdateIssue{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
	}
	g := getAdapter(input.Config)
	g.Update(input.Issue)

	return input.Issue
}

func (s *Subscriber) setup(nc *nats.Conn, m *nats.Msg) {
	body := m.Data
	input := messages.Setup{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		nc, _ := nats.Connect(nats.DefaultURL)
		nc.Publish(m.Reply, []byte(`{"error":"error setting up github labels"}`))
		return
	}
	g := getAdapter(input.Config)
	g.Labels = input.States
	g.Setup(input.States, input.Org, input.Repo)
	nc.Publish(m.Reply, []byte(`success`))
}

func getAdapter(config messages.Config) *Github {
	g := Github{
		Token: config.Github.Token,
	}
	g.setup()

	return &g
}
