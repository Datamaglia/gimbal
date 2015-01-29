package runner

func checkWrapper(w *Wrapper) {
    rs := new(ResultSet)
    w.ResultSet = rs

    // MaxAttempts
    attemptsResult := new(Result)
    attempts := w.Attempt()
    maxAttempts := w.Spec.MaxAttempts
    if attempts > 1 {
        if attempts >= maxAttempts && ! w.Success() {
            attemptsResult.Message = "Too many attempts"
            attemptsResult.Status = FAILURE
        } else {
            attemptsResult.Message = "More than one attempt"
            attemptsResult.Status = WARNING
        }
    } else {
        attemptsResult.Message = "One attempt"
        attemptsResult.Status = SUCCESS
    }
    attemptsResult.Expected = maxAttempts
    attemptsResult.Observed = attempts
    rs.AddResult(attemptsResult)

    // FIXME: The rest of these will fail if the http request never succeeded
    // because the response on the attempt struct will be nil

    // StatusCode
    statusResult := new(Result)
    expectedCode := w.Spec.StatusCode
    observedCode := w.LastAttempt().Resp.StatusCode
    if observedCode != expectedCode {
        statusResult.Message = "Status codes do not match"
        statusResult.Status = FAILURE
    } else {
        statusResult.Message = "Status codes match"
        statusResult.Status = SUCCESS
    }
    statusResult.Expected = expectedCode
    statusResult.Observed = observedCode
    rs.AddResult(statusResult)
}
