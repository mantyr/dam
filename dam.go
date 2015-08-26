// dam provides a simple method off adding flood protection to an http application.
package dam

import (
	"net/http"
	"time"
)

type Dam struct {
	per     int
	limit   int
	current int
	curChan chan int

	// Report adds the ability to log after completing the middleware,
	// for example:
	//
	//  d := dam.New(100, 1)
	//  d.Report = func(took *time.Duration) {
	//      log.Printf("middleware took: %v", took)
	//  }
	Report          func(time.Duration)
	ReportRejection func()

	// RejectCode adds the ability to override the default rejection
	// http code -- http.StatusNotAcceptable, for example:
	//
	//  d := dam.New(100, 1)
	//  d.RejectCode = http.StatusServiceUnavailable
	RejectCode int
}

func New(limit, perSecond int) *Dam {
	d := &Dam{
		per:     perSecond,
		limit:   limit,
		curChan: make(chan int, 10),
	}

	// See: https://en.wikipedia.org/wiki/List_of_HTTP_status_codes
	// 429 Too Many Requests
	d.RejectCode = 429
	return d
}

func (dam *Dam) reset() {
	dam.curChan <- 0
}

func (dam *Dam) increment() {
	current := dam.current + 1
	dam.curChan <- current
}

// Exceeded returns true if the current rate count is greater then the rate limit.
func (dam *Dam) Exceeded() bool {
	return !(dam.current < dam.limit)
}

func (dam *Dam) update() {
	for {
		dam.current = <-dam.curChan
	}
}

func (dam *Dam) ticker() {
	ticker := time.Tick(time.Duration(dam.per) * time.Second)
	for _ = range ticker {
		dam.reset()
	}
}

// Start starts the flood protection helpers, which keep track of rate and reset
// the counter based on the specified interval.
func (dam *Dam) Start() {
	go dam.ticker()
	go dam.update()
}

// Protect is the HTTP Middleware for accepting or rejecting http requests based
// on the current rate limit.
func (dam *Dam) Protect(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		begin := time.Now()
		defer func() {
			go dam.increment()

			if dam.Report != nil {
				dam.Report(time.Since(begin))
			}
		}()

		if dam.Exceeded() {
			if dam.ReportRejection != nil {
				dam.ReportRejection()
			}
			w.WriteHeader(dam.RejectCode)
			return
		}

		handler.ServeHTTP(w, req)
	})
}
