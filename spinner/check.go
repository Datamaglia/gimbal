package spinner

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type TestStatus int
const (
    FAILURE TestStatus = iota
    SUCCESS
    WARNING
    UNKNOWN
)

type TestResult struct {
    Name string
    Url string
    Method string
    Status TestStatus
    Expected interface{}
    Observed interface{}
    Message string
}

func (a ResponseSpec) checkConnection(wrapper *TestWrapper) TestResult {
    result := TestResult{}

    result.Name = "Connection"
    result.Url = wrapper.Spec.Request.FullUrl()
    result.Method = wrapper.Request.Method
    result.Expected = "connected"

    if wrapper.Err != nil {
        result.Status = FAILURE
        result.Observed = "failed to connect"
    } else {
        result.Status = SUCCESS
        result.Observed = "connected"
    }

    return result
}

func (a ResponseSpec) checkStatusCode(resp *http.Response) TestResult {
    result := TestResult{}

    result.Name = "Status code"
    result.Method = resp.Request.Method
    result.Expected = a.StatusCode
    result.Observed = resp.StatusCode

    if a.StatusCode == 0 {
        result.Status = UNKNOWN
    } else {
        success := resp.StatusCode == a.StatusCode
        if success {
            result.Status = SUCCESS
        } else {
            result.Status = FAILURE
        }
    }

    return result
}

func (a ResponseSpec) checkHeaders(resp *http.Response) TestResult {
    result := TestResult{}

    result.Name = "Headers"
    result.Method = resp.Request.Method
    expectedJson, _ := json.MarshalIndent(a.Headers, "      ", "  ")
    result.Expected = fmt.Sprintf("%s", expectedJson)
    observedJson, _ := json.MarshalIndent(resp.Header, "      ", "  ")
    result.Observed = fmt.Sprintf("%s", observedJson)

    if a.Headers == nil {
        result.Status = UNKNOWN
    } else {
        success := true
        for rawAssertedHeader, assertedValues := range a.Headers {
            assertedHeader := http.CanonicalHeaderKey(rawAssertedHeader)
            receivedValues := resp.Header[assertedHeader]
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
            result.Status = SUCCESS
        } else {
            result.Status = FAILURE
        }
    }

    return result
}

func (a ResponseSpec) checkTimeElapsed(wrapper *TestWrapper) TestResult {
    result := TestResult{}

    result.Name = "Time elapsed"
    result.Method = wrapper.Response.Request.Method
    result.Expected = fmt.Sprintf("<%v", a.TimeElapsed)
    result.Observed = fmt.Sprintf("%v", wrapper.TimeElapsed)

    if a.TimeElapsed == 0 {
        result.Status = UNKNOWN
    } else {
        success := wrapper.TimeElapsed <= a.TimeElapsed
        failure := wrapper.TimeElapsed - a.TimeElapsed > wrapper.Spec.Options.TimeElapsedDelta
        switch {
        case success:
            result.Status = SUCCESS
        case failure:
            result.Status = FAILURE
        default:
            result.Status = WARNING
        }
    }

    return result
}

func (a ResponseSpec) checkAttempts(wrapper *TestWrapper) TestResult {
    result := TestResult{}

    result.Name = "Attempts"
    result.Method = wrapper.Response.Request.Method
    result.Expected = a.Attempts
    result.Observed = wrapper.Attempt

    if a.Attempts == 0 {
        result.Status = UNKNOWN
    } else {
        success := wrapper.Attempt <= a.Attempts
        warning := wrapper.Attempt > 1
        switch {
        case success:
            result.Status = SUCCESS
        case warning:
            result.Status = WARNING
        default:
            result.Status = FAILURE
        }
    }

    return result
}
