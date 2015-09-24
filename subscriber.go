package main

import (
	"encoding/json"
	"github.com/nats-io/nats"
)

type Subscriber struct {
	nc *nats.Conn
	c  *nats.EncodedConn
}

func (s *Subscriber) Subscribe() {
	s.nc, _ = nats.Connect(nats.DefaultURL)
	s.c, _ = nats.NewEncodedConn(s.nc, nats.JSON_ENCODER)
	defer s.c.Close()

	c.Subscribe("issues.list", func(m *nats.Msg) {

	})
}
