package alarm

import (
	"log"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/eliothedeman/bangarang/event"
)

func init() {
	log.SetFlags(log.Llongfile)
}

var (
	DEFAULT_GROUP_BY = map[string]string{
		"host":        "^(.*)$",
		"service":     "^(.*)$",
		"sub_service": "^(.*)$",
	}
)

type Policy struct {
	Match       map[string]string `json:"match"`
	NotMatch    map[string]string `json:"not_match"`
	GroupBy     map[string]string `json:"group_by"`
	Crit        *Condition        `json:"crit"`
	Warn        *Condition        `json:"warn"`
	Name        string            `json:"name"`
	r_match     map[string]*regexp.Regexp
	r_not_match map[string]*regexp.Regexp
}

// check to see if an event satisfies the policy
func (p *Policy) Matches(e *event.Event) bool {
	return p.CheckMatch(e) && p.CheckNotMatch(e)
}

// compile the regex patterns for this policy
func (p *Policy) Compile() {
	logrus.Infof("Compiling regex maches for %s", p.Name)

	if p.r_match == nil {
		p.r_match = make(map[string]*regexp.Regexp)
	}

	if p.r_not_match == nil {
		p.r_not_match = make(map[string]*regexp.Regexp)
	}

	// if we don't have at least three componants of the group by, establish them from the defaults
	if len(p.GroupBy) < 3 {

		if len(p.GroupBy) == 0 {
			p.GroupBy = DEFAULT_GROUP_BY

		} else {

			tmp := map[string]string{}
			for k, v := range DEFAULT_GROUP_BY {
				tmp[k] = v
			}

			for k, v := range p.GroupBy {
				tmp[k] = v
			}

			p.GroupBy = tmp
		}
	}

	for k, v := range p.Match {
		p.r_match[k] = regexp.MustCompile(v)
	}

	for k, v := range p.NotMatch {
		p.r_not_match[k] = regexp.MustCompile(v)
	}

	if p.Crit != nil {
		logrus.Infof("Initializing crit for %s", p.Name)
		p.Crit.init(p.GroupBy)
	}

	if p.Warn != nil {
		logrus.Infof("Initializing warn for %s", p.Name)
		p.Warn.init(p.GroupBy)
	}
}

func formatFileName(n string) string {
	s := strings.Split(n, "_")
	a := ""
	for _, k := range s {
		a = a + strings.Title(k)
	}
	return a
}

// return the action to take for a given event
func (p *Policy) ActionCrit(e *event.Event) string {
	if p.Crit != nil {
		if p.Crit.TrackEvent(e) {
			e.Status = event.CRITICAL
		} else {
			e.Status = event.OK
		}

		if p.Crit.StateChanged(e) {
			return p.Crit.Escalation
		}
	}

	e.Status = event.OK
	return ""
}

func (p *Policy) ActionWarn(e *event.Event) string {
	if p.Warn != nil {
		if p.Warn.TrackEvent(e) {
			e.Status = event.WARNING
		} else {
			e.Status = event.OK
		}
		if p.Warn.StateChanged(e) {
			return p.Warn.Escalation
		}
	}

	e.Status = event.OK
	return ""
}

func (p *Policy) CheckNotMatch(e *event.Event) bool {
	for k, m := range p.r_not_match {
		if m.MatchString(e.Get(k)) {
			return false
		}
	}
	return true
}

// check if any of p's matches are satisfied by the event
func (p *Policy) CheckMatch(e *event.Event) bool {
	for k, m := range p.r_match {
		// if the element does not match the regex pattern, the event does not fully match
		if !m.MatchString(e.Get(k)) {
			return false
		}
	}
	return true
}
