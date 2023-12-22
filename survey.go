package main

// Column represents a column of a survey dataset.
type Column struct {
	code     string
	text     string
	minValue int
	maxValue int
	nullable bool
	ofType   string
}
