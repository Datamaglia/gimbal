package spinner

import (
    "fmt"
    "net/http"
    "strings"
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
    Host string
    SSL bool
    Uri string
    Port int
    Method string
    Headers http.Header
    Data string
}

// Builds a URL from component parts including host, SSL, URI, and port or, if
// set, returns the value of the `Url` attribute.
func (r *RequestSpec) FullUrl() string {
    if r.Url != "" {
        return r.Url
    }
    var proto string
    if r.SSL {
        proto = "https"
    } else {
        proto = "http"
    }
    return fmt.Sprintf("%v://%v%v:%v", proto, r.Host, r.Uri, r.Port)
}

func (r *RequestSpec) Update(defaults *RequestSpec) {
    if r.Method == "" {
        if defaults.Method != "" {
            r.Method = defaults.Method
        } else {
            r.Method = "GET"
        }
    }
    r.Method = strings.ToUpper(r.Method)
    // TODO Add more methods
    switch r.Method {
    case "GET", "POST", "PUT", "DELETE":
    default:
        panic("Invalid method")
    }

    if r.Url != "" {
        return
    }

    if r.Host == "" {
        if defaults.Host != "" {
            r.Host = defaults.Host
        } else {
            panic("Host not found")
        }
    }

    if defaults.SSL {
        r.SSL = true
    }

    if r.Uri == "" {
        if defaults.Uri != "" {
            r.Uri = defaults.Uri
        } else {
            r.Uri = "/"
        }
    }

    if r.Port == 0 {
        if defaults.Port != 0 {
            r.Port = defaults.Port
        } else {
            if r.SSL {
                r.Port = 443
            } else {
                r.Port = 80
            }
        }
    }

    // TODO Insert headers present in default but not in request
    if r.Headers == nil {
        if defaults.Headers != nil {
            r.Headers = defaults.Headers
        } else {
            r.Headers = http.Header{}
        }
    }

    if r.Data == "" {
        if defaults.Data != "" {
            r.Data = defaults.Data
        }
    }
}

type ResponseSpec struct {
    StatusCode int
    Headers http.Header
    TimeElapsed float64
    Attempts int
}

func (r *ResponseSpec) Update(defaults *ResponseSpec) {
    if r.StatusCode == 0 {
        if defaults.StatusCode != 0 {
            r.StatusCode = defaults.StatusCode
        }
    }

    // TODO Be smarter about merging the headers
    if defaults.Headers != nil {
        if r.Headers == nil {
            r.Headers = defaults.Headers
        } else {
            for header, value := range defaults.Headers {
                if r.Headers[header] == nil {
                    r.Headers[header] = value
                }
            }
        }
    }

    if r.TimeElapsed == 0.0 {
        if defaults.TimeElapsed > 0.0 {
            r.TimeElapsed = defaults.TimeElapsed
        }
    }

    if r.Attempts == 0 {
        if defaults.Attempts != 0 {
            r.Attempts = defaults.Attempts
        }
    }
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
