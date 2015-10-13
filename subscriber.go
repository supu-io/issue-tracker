package main

import (
	"encoding/json"
	"github.com/nats-io/nats"
	"log"
	"runtime"
)

// Collection of methods to subscribe the different events
type Subscriber struct{}

// Manages the subscriptions to different events.
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
		}
	})

	nc.Subscribe("issues.update", func(m *nats.Msg) {
		issues := s.issuesUpdate(m.Data)
		nc.Publish(m.Reply, *issues.toJSON())
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
	input := IssuesDetails{}
	err := json.Unmarshal(body, &input)
	if err != nil {
		log.Println(err)
	}

	g := input.Config.Github
	g.setup()
	issue := input.toIssue()
	if issue == nil {
		return nil
	}

	return g.Details(issue)
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
