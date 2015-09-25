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
	issues := g.List(input.Status)

	return *issues

}
