package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/eliothedeman/bangarang/pipeline"
	"github.com/gorilla/mux"
)

// EscalationConfig handles the api methods for incidents
type EscalationConfig struct {
	pipeline *pipeline.Pipeline
}

func NewEscalationConfig(pipe *pipeline.Pipeline) *EscalationConfig {
	return &EscalationConfig{
		pipeline: pipe,
	}
}

// EndPoint return the endpoint of this method
func (p *EscalationConfig) EndPoint() string {
	return "/api/escalation/config/{id}"
}

// Get HTTP get method
func (p *EscalationConfig) Get(w http.ResponseWriter, r *http.Request) {
	conf := p.pipeline.GetConfig()
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		logrus.Error("Must append escalation id", r.URL.String())
		http.Error(w, "must append escalation id", http.StatusBadRequest)
		return
	}

	if id == "*" {
		buff, err := json.Marshal(conf.Escalations)
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(buff)
		return
	}

	coll := conf.Escalations.Collection()
	logrus.Error(coll)
}

// Delete the given event provider
func (p *EscalationConfig) Delete(w http.ResponseWriter, r *http.Request) {
}

// Post HTTP get method
func (p *EscalationConfig) Post(w http.ResponseWriter, r *http.Request) {
}
