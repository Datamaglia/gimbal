package spinner

import (
    "net/http"
)

type TestStatus int
const (
    FAILURE TestStatus = iota
    SUCCESS
    WARNING
    UNKNOWN
)

type RequestSpec struct {
    Url string
    Method string
    Headers http.Header
    Data string
}

type ResponseSpec struct {
    StatusCode int
    Headers http.Header
    TimeElapsed float64
    Attempts int
}

type TestSpec struct {
    Request *RequestSpec
    Response *ResponseSpec
    Options *TestOptions
}

func (asserted ResponseSpec) checkStatusCode(response *http.Response) TestStatus {
    if asserted.StatusCode == 0 {
        return UNKNOWN
    }
    success := response.StatusCode == asserted.StatusCode
    if success {
        return SUCCESS
    } else {
        return FAILURE
    }
}

func (asserted ResponseSpec) checkHeaders(response *http.Response) TestStatus {
    if asserted.Headers == nil {
        return UNKNOWN
    }

    success := true

    for rawAssertedHeader, assertedValues := range asserted.Headers {
        assertedHeader := http.CanonicalHeaderKey(rawAssertedHeader)
        receivedValues := response.Header[assertedHeader]
        for _, assertedValue := range assertedValues {
            found := false
            for _, receivedValue := range receivedValues {
                if assertedValue == receivedValue {
                    found = true
                }
            }
            if ! found {
                success = false
            }
        }
    }

    if success {
        return SUCCESS
    } else {
        return FAILURE
    }
}

func (asserted ResponseSpec) checkTimeElapsed(timeElapsed float64) TestStatus {
    if asserted.TimeElapsed == 0 {
        return UNKNOWN
    }
    success := timeElapsed <= asserted.TimeElapsed
    if success {
        return SUCCESS
    } else {
        return FAILURE
    }
}

func (asserted ResponseSpec) checkAttempts(attempts int) TestStatus {
    if asserted.Attempts == 0 {
        return UNKNOWN
    }
    success := attempts <= asserted.Attempts
    warning := attempts > 1
    if success {
        return SUCCESS
    } else if warning {
        return WARNING
    } else {
        return FAILURE
    }
}
