package jira

import (
	gojira "github.com/andygrunwald/go-jira"
)

/*
 model view controller:
 - model = raw data
 - view = display to customer
		ie - screen on app, response from api
 - controller = converts data
*/


type Response struct {
	Issues []Issue
}

type Issue struct {
	URL     string `json:"url"`
	ID      string `json:"id"`
	Summary string `json:"title"`
}


func FromJiraIssue(issue gojira.Issue) Issue {
	return Issue{
		ID: issue.ID,
	}
}
