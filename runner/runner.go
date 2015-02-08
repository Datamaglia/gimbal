package runner

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/datamaglia/gimbal/printer"
	"github.com/datamaglia/gimbal/spec"
	"github.com/datamaglia/gimbal/wrapper"
)

// Create a custom ReadCloser because that's the only way we can get a string
// into the request body (because Go doesn't have type unions).
type BodyWrapper struct {
	io.Reader
}

func (BodyWrapper) Close() error { return nil }

func RunSpec(s *spec.Spec) {
	requestQueue := make(chan chan *wrapper.Wrapper, 0)
	checkQueue := make(chan *wrapper.Wrapper, 0)
	outputQueue := make(chan *wrapper.Wrapper, 0)

	// TODO Consider making TIMEOUT configurable / overridable
	client := http.Client{Timeout: TIMEOUT * time.Second}

	// Create the request handler pool
	var requestGroup sync.WaitGroup
	for i := 0; i < s.MaxConcurrentRequests(); i++ {
		requestGroup.Add(1)
		go func() {
			for wrappers := range requestQueue {
				for nextWrapper := range wrappers {
					sendRequest(nextWrapper, &client)
					checkQueue <- nextWrapper
				}
			}
			requestGroup.Done()
		}()
	}

	// Create the checking handler pool
	var checkGroup sync.WaitGroup
	for i := 0; i < 5; i++ {
		checkGroup.Add(1)
		go func() {
			for nextWrapper := range checkQueue {
				checkWrapper(nextWrapper)
				outputQueue <- nextWrapper
			}
			checkGroup.Done()
		}()
	}

	// Create the output handler
	var outputGroup sync.WaitGroup
	outputGroup.Add(1)
	go func() {
		for nextWrapper := range outputQueue {
			printer.ResultsToConsole(nextWrapper)
		}
		outputGroup.Done()
	}()

	// Create and add queued
	runSpec(s, requestQueue)

	close(requestQueue)
	requestGroup.Wait()

	close(checkQueue)
	checkGroup.Wait()

	close(outputQueue)
	outputGroup.Wait()
}

func runSpec(s *spec.Spec, requestQueue chan<- (chan *wrapper.Wrapper)) {
	terminalChildren := s.TerminalChildren()

	if terminalChildren > 0 {
		wrappers := make(chan *wrapper.Wrapper, terminalChildren)
		for i := 0; i < s.ConcurrentRequests; i++ {
			requestQueue <- wrappers
		}
		for _, spec := range s.Specs {
			if spec.Terminal() {
				w := new(wrapper.Wrapper)
				w.Spec = spec
				wrappers <- w
			}
		}
		close(wrappers)
	}

	for _, spec := range s.Specs {
		if !spec.Terminal() {
			// TODO How kosher is recursion in Go?
			runSpec(s, requestQueue)
		}
	}
}

func sendRequest(w *wrapper.Wrapper, client *http.Client) {
	maxAttempts := w.Spec.MaxAttempts
	for w.AttemptCount() < maxAttempts {
		req := new(http.Request)
		req.Method = w.Spec.Method
		req.URL = w.Spec.Url()
		req.Header = *w.Spec.RequestHeaders
		req.Body = BodyWrapper{bytes.NewBufferString(w.Spec.RequestData)}

		start := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(start)

		att := wrapper.Attempt{req, resp, err, elapsed.Seconds()}
		w.AddAttempt(&att)

		if att.Err == nil {
			break
		}
	}
}
