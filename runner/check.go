package runner

import (
	"github.com/datamaglia/gimbal/wrapper"
)

func checkWrapper(w *wrapper.Wrapper) {
	w.Results = make([]*wrapper.Result, 0)

	checkAttempts(w)

	// If no connection was made then all other tests are invalid
	if w.Status() == wrapper.FAILURE {
		return
	}

	checkStatusCode(w)
	checkTimeElapsed(w)
	checkResponseHeaders(w)
}

func checkAttempts(w *wrapper.Wrapper) {
	result := new(wrapper.Result)
	result.CheckName = "Number of attempts"
	attempts := w.AttemptCount()
	maxAttempts := w.Spec.MaxAttempts
	if attempts > 1 {
		if attempts >= maxAttempts && !w.Success() {
			result.Status = wrapper.FAILURE
		} else {
			result.Status = wrapper.WARNING
		}
	} else {
		result.Status = wrapper.SUCCESS
	}
	result.Expected = maxAttempts
	result.Observed = attempts

	w.AddResult(result)
}

func checkStatusCode(w *wrapper.Wrapper) {
	result := new(wrapper.Result)
	result.CheckName = "Status code"
	expectedCode := w.Spec.StatusCode
	observedCode := w.LastAttempt().Resp.StatusCode
	if observedCode != expectedCode {
		result.Status = wrapper.FAILURE
	} else {
		result.Status = wrapper.SUCCESS
	}
	result.Expected = expectedCode
	result.Observed = observedCode

	w.AddResult(result)
}

func checkTimeElapsed(w *wrapper.Wrapper) {
	result := new(wrapper.Result)
	result.CheckName = "Time elapsed"
	expectedTimeElapsed := w.Spec.MaxTimeElapsed
	observedTimeElapsed := w.LastAttempt().TimeElapsed
	timeElapsedDelta := w.Spec.TimeElapsedDelta
	if observedTimeElapsed > expectedTimeElapsed {
		if (observedTimeElapsed - expectedTimeElapsed) < timeElapsedDelta {
			result.Status = wrapper.WARNING
		} else {
			result.Status = wrapper.FAILURE
		}
	} else {
		result.Status = wrapper.SUCCESS
	}
	result.Expected = expectedTimeElapsed
	result.Observed = observedTimeElapsed

	w.AddResult(result)
}

func checkResponseHeaders(w *wrapper.Wrapper) {
	expectedHeaders := *w.Spec.ResponseHeaders
	observedHeaders := w.LastAttempt().Resp.Header
	results := make([]*wrapper.Result, 0)
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
			result := new(wrapper.Result)
			result.CheckName = "Response headers"
			if found {
				result.Status = wrapper.SUCCESS
				result.Expected = expectedValue
				result.Observed = expectedValue
			} else {
				result.Status = wrapper.FAILURE
				result.Expected = expectedValue
				result.Observed = ""
			}
			results = append(results, result)
		}
	}

	for _, r := range results {
		w.AddResult(r)
	}
}
