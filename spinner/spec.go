package spinner

import (
    "fmt"
    "net/http"
    "strings"
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
    return fmt.Sprintf("%v://%v:%v%v", proto, r.Host, r.Port, r.Uri)
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
    case "GET", "POST", "PATCH", "PUT", "DELETE":
    default:
        panic("Invalid method")
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

    if r.Url == "" {
        if defaults.Url != "" {
            r.Url = defaults.Url
        }
    }
    // Skip the rest of the checks if we have a full URL
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

    if r.TimeElapsed == 0 {
        if defaults.TimeElapsed > 0 {
            r.TimeElapsed = defaults.TimeElapsed
        }
    }

    if r.Attempts == 0 {
        if defaults.Attempts != 0 {
            r.Attempts = defaults.Attempts
        }
    }
}
