package runner

import (
	"fmt"
)

func checkWrapper(w *Wrapper) {
	rs := new(ResultSet)
	var r *Result

	r = checkAttempts(w)
	rs.AddResult(r)

	// If the attempts failed, none of the other checks will work.
	if r.Status == FAILURE {
		return
	}

	r = checkStatusCode(w)
	rs.AddResult(r)

	r = checkTimeElapsed(w)
	rs.AddResult(r)

	for _, r = range checkResponseHeaders(w) {
		rs.AddResult(r)
	}

	w.ResultSet = rs
}

func checkAttempts(w *Wrapper) *Result {
	result := new(Result)
	attempts := w.Attempt()
	maxAttempts := w.Spec.MaxAttempts
	if attempts > 1 {
		if attempts >= maxAttempts && !w.Success() {
			result.Message = "Too many attempts"
			result.Status = FAILURE
		} else {
			result.Message = "More than one attempt"
			result.Status = WARNING
		}
	} else {
		result.Message = "One attempt"
		result.Status = SUCCESS
	}
	result.Expected = maxAttempts
	result.Observed = attempts

	return result
}

func checkStatusCode(w *Wrapper) *Result {
	result := new(Result)
	expectedCode := w.Spec.StatusCode
	observedCode := w.LastAttempt().Resp.StatusCode
	if observedCode != expectedCode {
		result.Message = "Status codes do not match"
		result.Status = FAILURE
	} else {
		result.Message = "Status codes match"
		result.Status = SUCCESS
	}
	result.Expected = expectedCode
	result.Observed = observedCode

	return result
}

func checkTimeElapsed(w *Wrapper) *Result {
	result := new(Result)
	expectedTimeElapsed := w.Spec.MaxTimeElapsed
	observedTimeElapsed := w.LastAttempt().TimeElapsed
	timeElapsedDelta := w.Spec.TimeElapsedDelta
	if observedTimeElapsed > expectedTimeElapsed {
		if (observedTimeElapsed - expectedTimeElapsed) < timeElapsedDelta {
			result.Message = "Request was slower than desired"
			result.Status = WARNING
		} else {
			result.Message = "Request was too slow"
			result.Status = FAILURE
		}
	} else {
		result.Message = "Request was fast enough"
		result.Status = SUCCESS
	}
	result.Expected = expectedTimeElapsed
	result.Observed = observedTimeElapsed

	return result
}

func checkResponseHeaders(w *Wrapper) []*Result {
	expectedHeaders := *w.Spec.ResponseHeaders
	observedHeaders := w.LastAttempt().Resp.Header
	results := make([]*Result, len(expectedHeaders))
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
			result := new(Result)
			if found {
				result.Message = fmt.Sprintf("Header %s was correct", header)
				result.Status = SUCCESS
				result.Expected = expectedValue
				result.Observed = expectedValue
			} else {
				result.Message = fmt.Sprintf("Header %s was incorrect", header)
				result.Status = FAILURE
				result.Expected = expectedValue
				result.Observed = ""
			}
			results = append(results, result)
		}
	}

	return results
}
