package dam

import (
	"log"
	"net/http"
	"time"
)

func Example() {
	d := New(1000, 1)

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

func ExampleDam_Report() {
	d := New(1000, 1)

	d.Report = func(took time.Duration) {
		log.Printf("at=Protect took=%v\n", took)
	}

	index := d.Protect(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	d.Start()
	log.Fatal(http.ListenAndServe(":3000", index))
}

func ExampleDam_ReportRejection() {
	d := New(1000, 1)
	d.ReportRejection = func() {
		log.Printf("at=Protect status=rejected\n")
	}

	index := d.Protect(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	d.Start()
	log.Fatal(http.ListenAndServe(":3000", index))
}

func ExampleDam_RejectCode() {
	d := New(1000, 1)

	d.RejectCode = http.StatusServiceUnavailable

	index := d.Protect(
		// http handler func, on accept
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	d.Start()
	log.Fatal(http.ListenAndServe(":3000", index))
}

func ExampleDam_Protect() {
	d := New(1000, 1)

	index := d.Protect(
		// http handler func, on accept
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	d.Start()
	log.Fatal(http.ListenAndServe(":3000", index))
}
