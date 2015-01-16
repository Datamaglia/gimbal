package spinner

import (
    "net/http"
    "strings"
    "time"
    "sync"
)

type TestWrapper struct {
    Spec *TestSpec
    Response *http.Response
    TimeElapsed float64
    Attempt int
    Err error
}

func getRequest(wrapper *TestWrapper) {
    start := time.Now()
    resp, err := http.Get(wrapper.Spec.Request.Url)
    elapsed := time.Since(start)

    wrapper.Response = resp
    wrapper.Err = err
    wrapper.TimeElapsed = elapsed.Seconds()
}

func postRequest(wrapper *TestWrapper) {
}

func requestHandler(reqChan <-chan *TestWrapper, respChan chan<- *TestWrapper,
        waitGroup *sync.WaitGroup) {
    for reqWrapper := range reqChan {
        reqWrapper.Attempt += 1

        switch strings.ToUpper(reqWrapper.Spec.Request.Method) {
        case "POST":
            postRequest(reqWrapper)
        default:
            getRequest(reqWrapper)
        }

        respChan <- reqWrapper
    }
}

func responseHandler(respChan <-chan *TestWrapper, reqChan chan<- *TestWrapper,
        outChan chan<- *TestWrapper, waitGroup *sync.WaitGroup) {
    for respWrapper := range respChan {
        if respWrapper.Err != nil {
            // Request failed
            if respWrapper.Attempt < respWrapper.Spec.Options.MaxAttempts {
                // Try again
                reqChan <- respWrapper
                continue
            }
        }

        outChan <- respWrapper
    }
}

func outputHandler(outChan <-chan *TestWrapper, waitGroup *sync.WaitGroup) {
    for respWrapper := range outChan {
        println(respWrapper.Spec.Request.Url)

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

// TODO We need to figure out a safe way to shut the whole thing down. The
// problem is that the responseHandler goroutine re-submits failed jobs, so the
// main goroutine has no idea when the work is done. One option would be to
// tell the outputHandler goroutine how many specs there are and once it has
// seen that many have it initiate a shutdown somehow (not sure exactly how
// since it only has access to the outChan, which the others don't read from).
// This isn't ideal, though, because it precludes processing streams of specs.

func ExecuteTestConfig(config *TestConfig) {
    reqChan := make(chan *TestWrapper)
    respChan := make(chan *TestWrapper)
    outChan := make(chan *TestWrapper)

    waitGroup := new(sync.WaitGroup)
    waitGroup.Add(3)

    go requestHandler(reqChan, respChan, waitGroup)
    go responseHandler(respChan, reqChan, outChan, waitGroup)
    go outputHandler(outChan, waitGroup)

    for _, spec := range config.Specs {
        wrapper := new(TestWrapper)
        wrapper.Spec = spec
        reqChan <- wrapper
    }

    waitGroup.Wait()
}

// TODO Make this a method on the TestStatus so it can change the way it prints
func printStatus(status TestStatus, attribute string) {
    if status == UNKNOWN {
        return
    }
    if status == SUCCESS {
        println("  \u2713", attribute)
    }
    if status == WARNING {
        println("  \u2713", attribute)
    }
    if status == FAILURE {
        println("  \u2718", attribute)
    }
}
