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

func requestHandler(reqChan <-chan *TestWrapper, outChan chan<- *TestWrapper,
        waitGroup *sync.WaitGroup) {
    defer waitGroup.Done()

    for reqWrapper := range reqChan {
        for reqWrapper.Attempt < reqWrapper.Spec.Options.MaxAttempts {
            reqWrapper.Attempt += 1

            switch strings.ToUpper(reqWrapper.Spec.Request.Method) {
            case "POST":
                postRequest(reqWrapper)
            default:
                getRequest(reqWrapper)
            }

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
        println("  \u2713", attribute)
    }
    if status == WARNING {
        println("  \u2713", attribute)
    }
    if status == FAILURE {
        println("  \u2718", attribute)
    }
}
