deps:
	go get -u github.com/google/go-github/github
	go get -u golang.org/x/oauth2
	go get -u github.com/nats-io/nats
build:
	go build
test:
	go test
