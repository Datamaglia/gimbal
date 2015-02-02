package runner

import (
	"net/http"

	"github.com/datamaglia/gimbal/printer"
	"github.com/datamaglia/gimbal/spec"
)

type Attempt struct {
	Req         *http.Request
	Resp        *http.Response
	Err         error
	TimeElapsed float64
}

type Wrapper struct {
	Spec       *spec.Spec
	Attempts   []*Attempt
	ResultSets []*printer.ResultSet
}

func (w *Wrapper) AddAttempt(a *Attempt) {
	if w.Attempts == nil {
		w.Attempts = make([]*Attempt, 0)
	}
	w.Attempts = append(w.Attempts, a)
}

func (w *Wrapper) LastAttempt() *Attempt {
	return w.Attempts[len(w.Attempts)-1]
}

func (w *Wrapper) Attempt() int {
	if w.Attempts == nil {
		return 0
	}
	return len(w.Attempts)
}

func (w *Wrapper) Success() bool {
	return w.Attempts[len(w.Attempts)-1].Err == nil
}
