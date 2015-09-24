deps:
	go get -u github.com/google/go-github/github
	go get golang.org/x/oauth2
build:
	go build
test:
	go test
