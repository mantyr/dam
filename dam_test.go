package dam

import (
	"log"
	"net/http"
	"os"
)

func ExampleDam_New() {
	d := New(1000, 1)

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

	index := d.Protect(
		// http handler func, on accept
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	d.Start()
	log.Fatal(http.ListenAndServe(":3000", index))
}
