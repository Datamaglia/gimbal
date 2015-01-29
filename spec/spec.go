package spec

import (
    "fmt"
    "net/http"
    "net/url"
)

type Spec struct {

    // Config

    Name string
    ConcurrentRequests int
    OutputLevel string
    MaxAttempts int

    // Request

    Host string
    Uri string
    Port int
    Method string
    Ssl bool
    RequestHeaders *http.Header
    RequestData string

    // Response

    StatusCode int
    MaxTimeElapsed float64
    TimeElapsedDelta float64
    ResponseHeaders *http.Header
    ExactResponseHeaders bool
    ResponseData string

    // Children

    Specs []*Spec

    // Metadata

    terminals int
    maxConcurrentRequests int
    url *url.URL

}

// Counts the terminal ancestors of this spec and returns the result.
func (s *Spec) TerminalAncestors() int {
    if s.terminals == 0 {
        if s.Specs == nil {
            s.terminals = 1
        } else {
            total := 0
            for _, child := range s.Specs {
                total += child.TerminalAncestors()
            }
            s.terminals = total
        }
    }
    return s.terminals
}

// Counts the terminal children of this spec and returns the result.
func (s *Spec) TerminalChildren() int {
    total := 0
    for _, child := range s.Specs {
        if child.Specs == nil {
            total += 1
        }
    }
    return total
}

// Returns whether or not the spec is a terminal spec.
func (s *Spec) Terminal() bool {
    return s.Specs == nil
}

// Returns the largest concurrentRequests value in the tree.
func (s *Spec) MaxConcurrentRequests() int {
    if s.maxConcurrentRequests == 0 {
        if s.Specs == nil {
            s.maxConcurrentRequests = s.ConcurrentRequests
        } else {
            globalMax := 0
            for _, child := range s.Specs {
                childMax := child.MaxConcurrentRequests()
                if childMax > globalMax {
                    globalMax = childMax
                }
            }
            s.maxConcurrentRequests = globalMax
        }
    }
    return s.maxConcurrentRequests
}

func (s *Spec) UrlString() string {
    var proto string
    if s.Ssl {
        proto = "https"
    } else {
        proto = "http"
    }
    //return fmt.Sprintf("%v://%v:%v%v", proto, s.Host, s.Port, s.Uri)
    // FIXME: Why doesn't the port work?!?!
    return fmt.Sprintf("%v://%v%v", proto, s.Host, s.Uri)
}

func (s *Spec) Url() *url.URL {
    if s.url == nil {
        s.setUrl()
    }
    return s.url
}

func (s *Spec) setUrl() {
    url, err := url.Parse(s.UrlString())
    // TODO Improve this error message to include info about url and spec
    if err != nil {
        panic("Invalid URL")
    }
    s.url = url
}

// Set default values based on system defaults.
func (s *Spec) setDefaults() {
    if s.ConcurrentRequests == 0 {
        s.ConcurrentRequests = CONCURRENT_REQUESTS
    }
    if s.OutputLevel == "" {
        s.OutputLevel = OUTPUT_LEVEL
    }
    if s.MaxAttempts == 0 {
        s.MaxAttempts = MAX_ATTEMPTS
    }
    if s.Uri == "" {
        s.Uri = URI
    }
    if s.Port == 0 {
        if s.Ssl {
            s.Port = SSL_PORT
        } else {
            s.Port = NON_SSL_PORT
        }
    }
    if s.Method == "" {
        s.Method = METHOD
    }
    if s.RequestHeaders == nil {
        s.RequestHeaders = new(http.Header)
    }
    if s.StatusCode == 0 {
        s.StatusCode = STATUS_CODE
    }
    if s.MaxTimeElapsed == 0 {
        s.MaxTimeElapsed = MAX_TIME_ELAPSED
    }
    if s.ResponseHeaders == nil {
        s.ResponseHeaders = new(http.Header)
    }

    for _, spec := range s.Specs {
        spec.inheritDefaults(s)
    }
}

// Set default values based on a parent spec.
func (s *Spec) inheritDefaults(d *Spec) {
    if s.Name == "" {
        s.Name = d.Name
    } else {
        if d.Name != "" {
            s.Name = d.Name + " :: " + s.Name
        }
    }
    if s.ConcurrentRequests == 0 {
        s.ConcurrentRequests = d.ConcurrentRequests
    }
    if s.OutputLevel == "" {
        s.OutputLevel = d.OutputLevel
    }
    if s.MaxAttempts == 0 {
        s.MaxAttempts = d.MaxAttempts
    }
    if s.Host == "" {
        s.Host = d.Host
    }
    if s.Uri == "" {
        s.Uri = d.Uri
    }
    if s.Port == 0 {
        s.Port = d.Port
    }
    if s.Method == "" {
        s.Method = d.Method
    }
    if s.RequestHeaders == nil {
        s.RequestHeaders = d.RequestHeaders
    }
    if s.StatusCode == 0 {
        s.StatusCode = d.StatusCode
    }
    if s.MaxTimeElapsed == 0 {
        s.MaxTimeElapsed = d.MaxTimeElapsed
    }
    if s.ResponseHeaders == nil {
        s.ResponseHeaders = d.ResponseHeaders
    }

    for _, spec := range s.Specs {
        spec.inheritDefaults(s)
    }
}
