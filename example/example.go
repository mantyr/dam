package main

import (
	".."

	"fmt"
	"log"
	"net/http"
	"os"
	//"time"
)

func main() {
	d := dam.New(1000, 1)

	var log = log.New(os.Stdout, "[dam:example]: ", 0)

	// optionally add logging
	// d.Report = func(took time.Duration) {
	// 	log.Printf("at=Protect took=%v\n", took)
	// }

	// d.ReportRejection = func() {
	// 	log.Printf("at=Protect status=rejected\n")
	// }

	// optionally override the default RejectCode
	// d.RejectCode = http.StatusServiceUnavailable

	// basic http success handler
	index := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// writing custom middleware
	custom := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				// always do this in a goroutine, otherwise performance goes
				// to hell
				go d.Increment()

				log.Print("custom middleware complete")
			}()

			if d.Exceeded() {
				w.WriteHeader(429)
				fmt.Fprintf(w, "<html><body> Chill out man! </body></html>")
				log.Print("custom middleware rejected request")
				return
			}

			handler.ServeHTTP(w, r)
		})
	}

	d.Start()

	http.Handle("/", d.Protect(index))
	http.Handle("/custom", custom(index))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
