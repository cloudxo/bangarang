package pipeline

import (
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/eliothedeman/bangarang/event"
)

var (
	EMPTY_SUB_SERVICE = ""
)

// Provides stat tracking for events
type Tracker struct {
	inChan            chan *event.Event
	started           atomic.Value
	queryChan         chan QueryFunc
	total             *counter
	totalIncidents    counter
	incidentResolvers map[string]chan *event.Incident
	hosts             map[string]*counter
	hostTimes         map[string]time.Time
	services          map[string]*counter
	subServices       map[string]*counter
}

func (t *Tracker) Started() bool {
	if t.started.Load() == nil {
		return false
	} else {
		return t.started.Load().(bool)
	}
}

// create and return a new *Tracker
func NewTracker() *Tracker {
	t := &Tracker{
		inChan:            make(chan *event.Event),
		queryChan:         make(chan QueryFunc),
		total:             &counter{},
		incidentResolvers: make(map[string]chan *event.Incident),
		hosts:             make(map[string]*counter),
		hostTimes:         make(map[string]time.Time),
		services:          make(map[string]*counter),
		subServices:       make(map[string]*counter),
	}

	return t
}

// holds information about the current state of an event tracker
type TrackerReport struct {
	Total          uint64            `json:"total_events"`
	ByHost         map[string]uint64 `json:"by_host"`
	LastSeenByHost map[string]int64  `json:"last_seen_by_host"`
	ByService      map[string]uint64 `json:"by_service"`
	BySubService   map[string]uint64 `json:"by_sub_service"`
}

func NewReport() *TrackerReport {
	return &TrackerReport{
		ByHost:         make(map[string]uint64),
		ByService:      make(map[string]uint64),
		BySubService:   make(map[string]uint64),
		LastSeenByHost: make(map[string]int64),
	}
}

// TrackIncident will allow the tracker to keep state about an incident
func (t *Tracker) TrackIncident(i *event.Incident) {
	if i.GetResolve() != nil {
		t.Query(func(r *Tracker) {

			// Don't keep track of "OK" incident resolvers, as ok's can't be resolved
			if i.Status != event.OK {
				r.incidentResolvers[string(i.IndexName())] = i.GetResolve()
			}
			r.totalIncidents.inc()
		})
	}
}

func (t *Tracker) GetIncidentResolver(i *event.Incident) chan *event.Incident {
	var res chan *event.Incident
	t.Query(func(r *Tracker) {
		res, _ = r.incidentResolvers[string(i.IndexName())]
	})
	return res
}

// return a report of the current state of the tracker
func (t *Tracker) GetStats() *TrackerReport {
	logrus.Info("Fetching stats")
	r := NewReport()
	t.Query(func(t *Tracker) {
		r.Total = t.total.get()
		for k, v := range t.hosts {
			r.ByHost[k] = v.get()
		}
		for k, v := range t.services {
			r.ByService[k] = v.get()
		}
		for k, v := range t.subServices {
			r.BySubService[k] = v.get()
		}
		for k, v := range t.hostTimes {
			r.LastSeenByHost[k] = v.Unix()
		}
	})

	return r
}

func (t *Tracker) GetServices() []string {
	var services []string
	t.query(func(t *Tracker) {
		services = make([]string, len(t.services))
		x := 0
		for k, _ := range t.services {
			services[x] = k
			x += 1
		}
	})

	return services
}

// GetHosts returns all of the host names we have seen thus far
func (t *Tracker) GetHosts() []string {
	var hosts []string
	t.Query(func(t *Tracker) {
		hosts = make([]string, len(t.hostTimes))
		x := 0
		for k, _ := range t.hostTimes {
			hosts[x] = k
			x += 1
		}
	})

	return hosts
}

func (t *Tracker) RemoveHost(host string) {
	t.query(func(t *Tracker) {
		delete(t.hosts, host)
		delete(t.hostTimes, host)
	})
}

// return a map of hostnames to the last time we have heard from them
func (t *Tracker) HostTimes() map[string]time.Time {
	m := make(map[string]time.Time)
	t.Query(func(t *Tracker) {
		for k, v := range t.hostTimes {
			m[k] = v
		}
	})

	return m
}

// An function that is given access to the tracker without locks
type QueryFunc func(t *Tracker)

// Start the tracker. This should be done in it's own goroutine
func (t *Tracker) Start() {
	t.started.Store(true)
	logrus.Info("Starting event tracker")
	var e *event.Event
	var f QueryFunc
	for {
		select {
		case e = <-t.inChan:
			t.trackEvent(e)
		case f = <-t.queryChan:
			t.query(f)
		}
	}
}

// Process the stats for a given event
func (t *Tracker) TrackEvent(e *event.Event) {
	t.inChan <- e
}

// Perform a QueryFunc on the tracker syncronously
func (t *Tracker) Query(f QueryFunc) {
	done := make(chan struct{})
	t.queryChan <- func(t *Tracker) {
		f(t)
		done <- struct{}{}
	}
	<-done
}

func (t *Tracker) query(f QueryFunc) {
	f(t)
}

func (t *Tracker) trackEvent(e *event.Event) {
	// e.Id = t.total.get()

	// don't track keep alives
	if e.Service == KEEP_ALIVE_SERVICE_NAME {
		return
	}
	t.total.inc()

	// update the last time we have seen this host
	t.hostTimes[e.Host] = time.Now()

	// increment host counter
	host, ok := t.hosts[e.Host]
	if !ok {
		host = &counter{}
		t.hosts[e.Host] = host
	}
	host.inc()

	// increment service counter
	service, ok := t.services[e.Service]
	if !ok {
		service = &counter{}
		t.services[e.Service] = service
	}
	service.inc()

	// if the event has a sub_service, increment the sub_service counter
	if e.SubService != EMPTY_SUB_SERVICE {
		subService, ok := t.subServices[e.SubService]
		if !ok {
			subService = &counter{}
			t.subServices[e.SubService] = subService
		}

		subService.inc()
	}
}

type counter struct {
	c uint64
}

func (c *counter) inc() {
	atomic.AddUint64(&c.c, 1)
}
func (c *counter) set(val uint64) {
	atomic.StoreUint64(&c.c, val)
}

func (c *counter) get() uint64 {
	return atomic.LoadUint64(&c.c)
}
