package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/cloudfoundry/gosigar"
	"github.com/eliothedeman/bangarang/pipeline"
)

var startTime time.Time

func init() {
	// track when the program started
	startTime = time.Now()
}

// handles the api methods for incidents
type SystemStats struct {
	pipeline *pipeline.Pipeline
}

func NewSystemStats(pipe *pipeline.Pipeline) *SystemStats {
	return &SystemStats{
		pipeline: pipe,
	}
}

func (e *SystemStats) EndPoint() string {
	return "/api/stats/system"
}

func getMem() map[string]uint64 {
	m := sigar.Mem{}
	err := m.Get()
	a := map[string]uint64{}
	if err != nil {
		logrus.Error(err)
		return a
	}

	a["used"] = m.Used
	a["free"] = m.Free
	a["total"] = m.Total
	return a
}

func getLoad() map[string]float64 {
	l := sigar.LoadAverage{}
	err := l.Get()
	m := map[string]float64{}
	if err != nil {
		logrus.Error(err)
		return m
	}

	m["one"] = l.One
	m["five"] = l.Five
	m["fifteen"] = l.Fifteen
	return m
}

func getUptime() time.Duration {
	now := time.Now()
	return now.Sub(startTime)
}

func (e *SystemStats) Get(w http.ResponseWriter, r *http.Request) {
	m := map[string]interface{}{}

	m["memory"] = getMem()
	m["load"] = getLoad()
	m["uptime"] = getUptime().Seconds()

	buff, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(buff)
}
