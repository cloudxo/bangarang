package escalation

import (
	"encoding/json"
	"testing"

	"github.com/eliothedeman/bangarang/event"
)

const (
	test_policy_with_regex = `
	{
		"match": [
			{"key":"host", "value": "test\\.hello"}
		]
	}
`
)

type testingPasser struct {
	incidents map[string]*event.Incident
}

func (t *testingPasser) PassIncident(i *event.Incident) {
	if t.incidents == nil {
		t.incidents = map[string]*event.Incident{}
	}
	t.incidents[string(i.IndexName())] = i
}

func newTestPasser() event.IncidentPasser {
	return &testingPasser{}
}

func TestPolicyRegexParsing(t *testing.T) {
	p := &Policy{}

	err := json.Unmarshal([]byte(test_policy_with_regex), p)
	if err != nil {
		t.Error(err)
	}

	if p.Match.Get("host") != `test\.hello` {
		t.Error("regex not properly parsed")
	}

}

func TestMatchOr(t *testing.T) {
	p := &Policy{}
	e := &event.Event{}
	e.Tags = &event.TagSet{
		{"test_tag", "unknown"},
	}

	p.Match = &event.TagSet{
		{
			"test_tag", "unknown|shadow",
		},
	}
	p.Compile(newTestPasser())

	if !p.CheckMatch(e) {
		t.Fail()
	}

}

func TestMatchTagsMulti(t *testing.T) {
	p := &Policy{}
	e := &event.Event{}
	e.Tags = &event.TagSet{
		{
			"test_tag", "0",
		},
		{
			"other_tag", "what is this ice?",
		},
	}

	p.Match = &event.TagSet{
		{
			"test_tag", "[0-9]+",
		},
		{
			"other_tag", "ice",
		},
	}
	p.Compile(newTestPasser())

	if !p.CheckMatch(e) {
		t.Fail()
	}
}

func TestMatchTagsMultiNotMatch(t *testing.T) {
	p := &Policy{}
	e := &event.Event{}
	e.Tags = &event.TagSet{
		{
			"test_tag", "0",
		},
	}

	p.Match = &event.TagSet{
		{
			"test_tag", "[0-9]+",
		},
		{
			"other_tag", "ice",
		},
	}
	p.Compile(newTestPasser())

	if p.CheckMatch(e) {
		t.Fail()
	}
}

func TestMatchTagsSingle(t *testing.T) {
	p := &Policy{}
	e := &event.Event{}
	e.Tags = &event.TagSet{
		{
			"test_tag", "0",
		},
	}

	p.Match = &event.TagSet{
		{
			"test_tag", "[0-9]+",
		},
	}
	p.Compile(newTestPasser())

	if !p.CheckMatch(e) {
		t.Fail()
	}
}

func test_f(f float64) *float64 {
	return &f
}

func TestCompileWithCrit(t *testing.T) {
	p := &Policy{}
	p.Crit = &Condition{
		Greater: test_f(10.0),
		Less:    test_f(-0.1),
		Exactly: test_f(0.5),
	}

	p.Compile(newTestPasser())

}

func TestCompileSatisfies(t *testing.T) {
	p := &Policy{}
	p.Crit = &Condition{
		Greater: test_f(10.0),
		Less:    test_f(-0.1),
		Exactly: test_f(0.5),
	}

	p.Compile(newTestPasser())

	e := &event.Event{}

	e.Metric = 15
	if !p.Crit.Satisfies(e) {
		t.Fail()
	}

	e.Metric = 8
	if p.Crit.Satisfies(e) {
		t.Fail()
	}

}
