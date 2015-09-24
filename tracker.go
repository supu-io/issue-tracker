package main

type Tracker interface {
	setup()

	// Get a list of issues by status
	List(status string) *[]*Issue

	// Get issue details
	Details(id string) *Issue

	// Update an issue
	Update(id string, issue *Issue)

	// Add a comment on an issue
	Comment(id string, body string)
}
