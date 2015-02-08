package wrapper

import (
	"net/http"
)

type Attempt struct {
	Req         *http.Request
	Resp        *http.Response
	Err         error
	TimeElapsed float64

	// TODO: Timestamp
}

func (a *Attempt) Connected() bool {
	return a.Err == nil
}
