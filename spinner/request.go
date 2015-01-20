package spinner

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "sync"
    "time"

    "github.com/wsxiaoys/terminal/color"
)

type TestWrapper struct {
    Spec *TestSpec
    Request *http.Request
    Response *http.Response
    TimeElapsed float64
    Attempt int
    Err error
}

// Create a custom ReadCloser because that's the only way we can get a string
// into the request body (because Go doesn't have type unions).
type BodyWrapper struct {
    io.Reader
}
func (BodyWrapper) Close() error { return nil }

func requestHandler(reqChan <-chan *TestWrapper, outChan chan<- *TestWrapper,
        waitGroup *sync.WaitGroup) {
    defer waitGroup.Done()

    for reqWrapper := range reqChan {
        client := http.Client{Timeout: 10 * time.Second}

        for reqWrapper.Attempt < reqWrapper.Spec.Options.MaxAttempts {
            reqWrapper.Attempt += 1

            req := new(http.Request)
            req.Method = reqWrapper.Spec.Request.Method

            // TODO This conversion should be moved to the config read stage
            reqUrl, err := url.Parse(reqWrapper.Spec.Request.FullUrl())
            if err != nil {
                panic("Invalid URL")
            }
            req.URL = reqUrl

            req.Header = reqWrapper.Spec.Request.Headers
            req.Body = BodyWrapper{bytes.NewBufferString(reqWrapper.Spec.Request.Data)}

            start := time.Now()
            resp, err := client.Do(req)
            elapsed := time.Since(start)

            reqWrapper.Request = req
            reqWrapper.Response = resp
            reqWrapper.Err = err
            reqWrapper.TimeElapsed = elapsed.Seconds()

            if reqWrapper.Err == nil {
                // Succeeded, so we don't need any more attempts
                break
            }

            // Wait some small but significant amount of time before hitting
            // the server again
            time.Sleep(100 * time.Millisecond)
        }

        outChan <- reqWrapper
    }
}

func outputHandler(config *TestConfig, outChan <-chan *TestWrapper, waitGroup *sync.WaitGroup) {
    defer waitGroup.Done()

    for wrapper := range outChan {
        respSpec := wrapper.Spec.Response
        results := make([]TestResult, 0)

        if wrapper.Err != nil {
            results = append(results, respSpec.checkConnection(wrapper))
        } else {
            resp := wrapper.Response

            results = append(results, respSpec.checkConnection(wrapper))
            results = append(results, respSpec.checkStatusCode(resp))
            results = append(results, respSpec.checkHeaders(resp))
            results = append(results, respSpec.checkTimeElapsed(wrapper))
            results = append(results, respSpec.checkAttempts(wrapper))
        }

        for _, result := range results {
            if result.Status == WARNING {
                config.TotalWarnings += 1
            }
            if result.Status == FAILURE {
                config.TotalFailures += 1
            }
            if result.Status != UNKNOWN {
                config.TotalTests += 1
            }

            config.PrintResult(result)
        }
    }
}

func ExecuteTestConfig(config *TestConfig) {
    concurrentRequests := config.Settings.ConcurrentRequests

    reqChan := make(chan *TestWrapper, 10)
    outChan := make(chan *TestWrapper)

    reqWaitGroup := new(sync.WaitGroup)
    reqWaitGroup.Add(concurrentRequests)

    for i := 0; i < concurrentRequests; i++ {
        go requestHandler(reqChan, outChan, reqWaitGroup)
    }

    outWaitGroup := new(sync.WaitGroup)
    outWaitGroup.Add(1)

    go outputHandler(config, outChan, outWaitGroup)

    for _, spec := range config.Specs {
        wrapper := new(TestWrapper)
        wrapper.Spec = spec
        reqChan <- wrapper
    }

    close(reqChan)
    reqWaitGroup.Wait()

    close(outChan)
    outWaitGroup.Wait()

    if ! config.Settings.SuppressOutput {
        config.PrintSummary()
    }
}

func (t *TestConfig) PrintResult(result TestResult) {
    if t.Settings.SuppressOutput {
        return
    }

    if result.Url != "" {
        fmt.Printf("%v (%v)\n", result.Url, result.Method)
    }

    status := result.Status

    if status == UNKNOWN {
        return
    }
    if status == SUCCESS && ! t.Settings.SuppressSuccess {
        color.Printf("@g  \u2713 %v\n", result.Name)
    }
    if status == WARNING && ! t.Settings.SuppressWarning {
        color.Printf("@y  \u2713 %v\n", result.Name)
        color.Printf("@y    Expected:\n      %v\n", result.Expected)
        color.Printf("@y    Observed:\n      %v\n", result.Observed)
    }
    if status == FAILURE && ! t.Settings.SuppressFailure {
        color.Printf("@r  \u2718 %v\n", result.Name)
        color.Printf("@r    Expected:\n      %v\n", result.Expected)
        color.Printf("@r    Observed:\n      %v\n", result.Observed)
    }
}

func (t *TestConfig) PrintSummary() {
    r := t.TotalTests
    w := t.TotalWarnings
    f := t.TotalFailures
    p := r - w - f
    fmt.Printf("\n")
    color.Printf("Summary: %d tests run, @g%d passed@|, @y%d warnings@|, @r%d failures@|", r, p, w, f)
}
