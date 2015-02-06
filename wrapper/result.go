package wrapper

import (
	"github.com/datamaglia/gimbal/spec"
)

type ResultStatus int

const (
	FAILURE ResultStatus = iota
	SUCCESS
	WARNING
	UNKNOWN
)

type Result struct {
	Message  string
	Status   ResultStatus
	Expected interface{}
	Observed interface{}
}

type ResultSet struct {
	Spec    *spec.Spec
	Results []*Result
}

func (s *ResultSet) AddResult(r *Result) {
	s.Results = append(s.Results, r)
}

func (s *ResultSet) Success() bool {
	success := true
	for _, result := range s.Results {
		if result.Status != SUCCESS {
			success = false
		}
	}
	return success
}

func (s *ResultSet) Status() ResultStatus {
	failure := false
	warning := false
	for _, result := range s.Results {
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
