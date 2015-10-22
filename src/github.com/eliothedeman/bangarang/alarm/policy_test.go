package alarm

import (
	"encoding/json"
	"testing"

	"github.com/eliothedeman/bangarang/event"
)

const (
	test_policy_with_regex = `
	{
		"match": {
			"host": "test\\.hello"
		}
	}
`
)

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
	e.Tags = event.TagSet{
		{"test_tag", "unknown"},
	}

	p.Match = event.TagSet{
		{
			"test_tag", "unknown|shadow",
		},
	}
	p.Compile()

	if !p.CheckMatch(e) {
		t.Fail()
	}

}

func TestMatchTags(t *testing.T) {
	p := &Policy{}
	e := &event.Event{}
	e.Tags = event.TagSet{
		{
			"test_tag", "0",
		},
	}

	p.Match = event.TagSet{
		{
			"test_tag", "[0-9]+",
		},
	}
	p.Compile()

	if !p.CheckMatch(e) {
		t.Fail()
	}
}

func TestMatchTagsExtra(t *testing.T) {
	p := &Policy{}
	e := &event.Event{}
	e.Tags = event.TagSet{
		{"test_tag", "0"},
		{"extra_tag", "w234"},
	}

	p.Match = event.TagSet{
		{
			"test_tag", "[0-9]+",
		},
	}
	p.Compile()

	if !p.CheckMatch(e) {
		t.Fail()
	}
}

func TestMatchStructFiled(t *testing.T) {
	p := &Policy{}
	e := newTestEvent("my_host", "", 0)

	p.Match = event.TagSet{
		{
			"host", "my.*",
		},
	}
	p.Compile()

	if !p.CheckMatch(e) {
		t.Fail()
	}

	e.Tags = append(e.Tags, event.KeyVal{"host", ""})
	if p.CheckMatch(e) {
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

	p.Compile()

}

func TestCompileSatisfies(t *testing.T) {
	p := &Policy{}
	p.Crit = &Condition{
		Greater: test_f(10.0),
		Less:    test_f(-0.1),
		Exactly: test_f(0.5),
	}

	p.Compile()

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
