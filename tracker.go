package main

// Tracker interface
type Tracker interface {
	setup()

	// Get a list of issues by status
	List(status string) *Issues

	// Get issue details
	Details(id string) *Issue

	// Update an issue
	Update(issue *Issue) []string

	// Add a comment on an issue
	Comment(id string, body string)
}
