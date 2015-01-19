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
        client := http.Client{}

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

func outputHandler(outChan <-chan *TestWrapper, waitGroup *sync.WaitGroup) {
    defer waitGroup.Done()

    for respWrapper := range outChan {
        fmt.Printf("%v (%v)\n", respWrapper.Spec.Request.FullUrl(),
                respWrapper.Spec.Request.Method)

        if respWrapper.Err != nil {
            printStatus(FAILURE, "Connection")
        } else {
            response := respWrapper.Response
            responseSpec := respWrapper.Spec.Response

            printStatus(SUCCESS, "Connection")
            printStatus(responseSpec.checkStatusCode(response), "Status code")
            printStatus(responseSpec.checkHeaders(response), "Headers")
            printStatus(responseSpec.checkTimeElapsed(respWrapper.TimeElapsed), "Time elapsed")
            printStatus(responseSpec.checkAttempts(respWrapper.Attempt), "Attempts")
        }
    }
}

func ExecuteTestConfig(config *TestConfig) {
    var concurrentRequests int
    if config.Settings.ConcurrentRequests < 1 {
        concurrentRequests = 1
    } else {
        concurrentRequests = config.Settings.ConcurrentRequests
    }

    reqChan := make(chan *TestWrapper, 10)
    outChan := make(chan *TestWrapper)

    reqWaitGroup := new(sync.WaitGroup)
    reqWaitGroup.Add(concurrentRequests)

    for i := 0; i < concurrentRequests; i++ {
        go requestHandler(reqChan, outChan, reqWaitGroup)
    }

    outWaitGroup := new(sync.WaitGroup)
    outWaitGroup.Add(1)

    go outputHandler(outChan, outWaitGroup)

    for _, spec := range config.Specs {
        wrapper := new(TestWrapper)
        wrapper.Spec = spec
        reqChan <- wrapper
    }

    close(reqChan)
    reqWaitGroup.Wait()

    close(outChan)
    outWaitGroup.Wait()
}

// TODO Make this a method on the TestStatus so it can change the way it prints
func printStatus(status TestStatus, attribute string) {
    if status == UNKNOWN {
        return
    }
    if status == SUCCESS {
        color.Println("@g  \u2713", attribute)
    }
    if status == WARNING {
        color.Println("@y  \u2713", attribute)
    }
    if status == FAILURE {
        color.Println("@r  \u2718", attribute)
    }
}
