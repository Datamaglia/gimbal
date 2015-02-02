package runner

import (
	"fmt"

	"github.com/datamaglia/gimbal/printer"
)

func checkWrapper(w *Wrapper) {
	rs := new(printer.ResultSet)
	rs.Spec = w.Spec
	var r *printer.Result

	r = checkAttempts(w)
	rs.AddResult(r)

	// If the attempts failed, none of the other checks will work.
	if r.Status == printer.FAILURE {
		return
	}

	r = checkStatusCode(w)
	rs.AddResult(r)

	r = checkTimeElapsed(w)
	rs.AddResult(r)

	for _, r = range checkResponseHeaders(w) {
		rs.AddResult(r)
	}

	w.ResultSets = append(w.ResultSets, rs)
}

func checkAttempts(w *Wrapper) *printer.Result {
	result := new(printer.Result)
	attempts := w.Attempt()
	maxAttempts := w.Spec.MaxAttempts
	if attempts > 1 {
		if attempts >= maxAttempts && !w.Success() {
			result.Message = "Too many attempts"
			result.Status = printer.FAILURE
		} else {
			result.Message = "More than one attempt"
			result.Status = printer.WARNING
		}
	} else {
		result.Message = "One attempt"
		result.Status = printer.SUCCESS
	}
	result.Expected = maxAttempts
	result.Observed = attempts

	return result
}

func checkStatusCode(w *Wrapper) *printer.Result {
	result := new(printer.Result)
	expectedCode := w.Spec.StatusCode
	observedCode := w.LastAttempt().Resp.StatusCode
	if observedCode != expectedCode {
		result.Message = "Status codes do not match"
		result.Status = printer.FAILURE
	} else {
		result.Message = "Status codes match"
		result.Status = printer.SUCCESS
	}
	result.Expected = expectedCode
	result.Observed = observedCode

	return result
}

func checkTimeElapsed(w *Wrapper) *printer.Result {
	result := new(printer.Result)
	expectedTimeElapsed := w.Spec.MaxTimeElapsed
	observedTimeElapsed := w.LastAttempt().TimeElapsed
	timeElapsedDelta := w.Spec.TimeElapsedDelta
	if observedTimeElapsed > expectedTimeElapsed {
		if (observedTimeElapsed - expectedTimeElapsed) < timeElapsedDelta {
			result.Message = "Request was slower than desired"
			result.Status = printer.WARNING
		} else {
			result.Message = "Request was too slow"
			result.Status = printer.FAILURE
		}
	} else {
		result.Message = "Request was fast enough"
		result.Status = printer.SUCCESS
	}
	result.Expected = expectedTimeElapsed
	result.Observed = observedTimeElapsed

	return result
}

func checkResponseHeaders(w *Wrapper) []*printer.Result {
	expectedHeaders := *w.Spec.ResponseHeaders
	observedHeaders := w.LastAttempt().Resp.Header
	results := make([]*printer.Result, len(expectedHeaders))
	for header, expectedValues := range expectedHeaders {
		observedValues := observedHeaders[header]
		for _, expectedValue := range expectedValues {
			found := false
			for _, observedValue := range observedValues {
				if expectedValue == observedValue {
					found = true
				}
			}
			// Add a result depending on whether the result was found
			result := new(printer.Result)
			if found {
				result.Message = fmt.Sprintf("Header %s was correct", header)
				result.Status = printer.SUCCESS
				result.Expected = expectedValue
				result.Observed = expectedValue
			} else {
				result.Message = fmt.Sprintf("Header %s was incorrect", header)
				result.Status = printer.FAILURE
				result.Expected = expectedValue
				result.Observed = ""
			}
			results = append(results, result)
		}
	}

	return results
}
