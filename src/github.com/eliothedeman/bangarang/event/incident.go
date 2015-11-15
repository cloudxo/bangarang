package event

import (
	"crypto/md5"
	"fmt"
	"time"
)

type IncidentFormatter func(*Incident) string

func DefaultIncidentFormatter(i *Incident) string {
	return fmt.Sprintf("%s on %s is %s. Triggered by %s", i.Tags.Get("service"), i.Tags.Get("host"), Status(i.Status), i.Policy)
}

// An incident is created whenever an event changes to a state that is not event.OK
type Incident struct {
	EventName   []byte `json:"event" msg:"event_name"`
	Time        int64  `json:"time" msg:"time"`
	Id          int64  `json:"id" msg:"id"`
	Active      bool   `json:"active" msg:"active"`
	Escalation  string `json:"escalation" msg:"escalation"`
	Description string `json:"description" msg:"description"`
	Policy      string `json:"policy" msg:"policy"`
	Status      int    `json:"status" "msg:"status"`
	indexName   []byte
	resChan     chan *Incident // this is used to call back to the policy that created this event
	Event
}

func (i *Incident) SetResolve(r chan *Incident) {
	i.resChan = r
}

func (i *Incident) GetResolve() chan *Incident {
	return i.resChan
}

func (i *Incident) IndexName() []byte {
	if len(i.indexName) == 0 {
		n := md5.New()
		n.Write([]byte(i.Policy + i.Event.Tags.String()))
		i.indexName = []byte(fmt.Sprintf("%x", n.Sum(nil)))
	}

	return i.indexName
}

func (i *Incident) GetEvent() *Event {
	return &i.Event
}

func (i *Incident) FormatDescription() string {
	return DefaultIncidentFormatter(i)
}

func NewIncident(policy string, escalation string, status int, e *Event) *Incident {
	in := &Incident{
		EventName:  []byte(e.IndexName()),
		Time:       time.Now().Unix(),
		Active:     true,
		Status:     status,
		Policy:     policy,
		Escalation: escalation,
	}

	in.Tags = e.Tags
	in.Metric = e.Metric
	in.Description = in.FormatDescription()

	return in
}
