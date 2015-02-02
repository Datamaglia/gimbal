package printer

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
