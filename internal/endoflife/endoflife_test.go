package endoflife

import "testing"

func TestGetProduct(t *testing.T) {
	jiraCycles, err := GetProduct("jira-software")
	if err != nil {
		t.Error(err)
	}
	if len(jiraCycles) < 99 {
		t.Errorf("To less results")
	}
}

func TestGetCycle(t *testing.T) {
	cycle, err := GetCycle("jira-software", "3.7")
	if err != nil {
		t.Error(err)
	}
	if cycle.Latest != "3.7.4" {
		t.Errorf("Wrong version %s x", cycle.Latest)
	}
}
