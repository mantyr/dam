# dam
Simple (Go) HTTP Flood Protection

## example

```go
package main

import (
	"gopkg.in/jmervine/dam.v1"

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

	index := d.Protect(
		// http handler func, on accept
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	d.Start()
	log.Fatal(http.ListenAndServe(":3000", index))
}
```

## perf info

Basic perf info using `example/example.go`

```
httperf --client=0/1 --server=localhost --port=3000 --uri=/ --send-buffer=4096 --recv-buffer=16384 --num-conns=10000 --num-calls=10
httperf: warning: open file limit > FD_SETSIZE; limiting max. # of open files to FD_SETSIZE
Maximum connect burst length: 1

Total: connections 10000 requests 100000 replies 100000 test-duration 3.090 s

Connection rate: 3235.9 conn/s (0.3 ms/conn, <=1 concurrent connections)
Connection time [ms]: min 0.2 avg 0.3 max 1.7 median 0.5 stddev 0.1
Connection time [ms]: connect 0.0
Connection length [replies/conn]: 10.000

Request rate: 32358.6 req/s (0.0 ms/req)
Request size [B]: 62.0

Reply rate [replies/s]: min 0.0 avg 0.0 max 0.0 stddev 0.0 (0 samples)
Reply time [ms]: response 0.0 transfer 0.0
Reply size [B]: header 130.0 content 0.0 footer 0.0 (total 130.0)
Reply status: 1xx=0 2xx=4003 3xx=0 4xx=95997 5xx=0

CPU time [s]: user 1.01 system 2.08 (user 32.6% system 67.3% total 99.9%)
Net I/O: 6079.9 KB/s (49.8*10^6 bps)

Errors: total 0 client-timo 0 socket-timo 0 connrefused 0 connreset 0
Errors: fd-unavail 0 addrunavail 0 ftab-full 0 other 0
```
