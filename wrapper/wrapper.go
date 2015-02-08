package wrapper

import (
	"github.com/datamaglia/gimbal/spec"
)

// TODO: Better error logging and graceful exit instead of just panic(...)

type Wrapper struct {
	Index    int
	Spec     *spec.Spec
	Attempts []*Attempt
	Results  []*Result
}

func (w *Wrapper) AddAttempt(a *Attempt) {
	if w.Attempts == nil {
		w.Attempts = make([]*Attempt, 0)
	}
	w.Attempts = append(w.Attempts, a)
}

func (w *Wrapper) AddResult(r *Result) {
	w.Results = append(w.Results, r)
}

func (w *Wrapper) AttemptCount() int {
	if w.Attempts == nil {
		return 0
	}
	return len(w.Attempts)
}

func (w *Wrapper) LastAttempt() *Attempt {
	if w.Attempts == nil {
		panic("Wrapper has not yet been attempted")
	}
	return w.Attempts[len(w.Attempts)-1]
}

func (w *Wrapper) Success() bool {
	return w.Status() == SUCCESS
}

func (w *Wrapper) Status() ResultStatus {
	if w.Results == nil {
		panic("Wrapper has not yet been checked")
	}

	failure := false
	warning := false
	for _, result := range w.Results {
		if result.Status == WARNING {
			warning = true
		}
		if result.Status == FAILURE {
			failure = true
		}
	}
	if failure {
		return FAILURE
	}
	if warning {
		return WARNING
	}
	return SUCCESS
}
