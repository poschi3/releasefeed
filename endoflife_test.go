package main

import (
	"testing"
)

func TestGetProduct(t *testing.T) {
	jiraCycles := getProduct("jira-software")
	if len(jiraCycles) < 99 {
		t.Errorf("To less results")
	}
}

func TestGetCycle(t *testing.T) {
	cycle := getCycle("jira-software", "3.7")
	if cycle.Latest != "3.7.4" {
		t.Errorf("Wrong version")
	}
}
