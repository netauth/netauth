package health

import (
	"log"
	"sync/atomic"
)

// The server starts out unhealthy by default, since it hasn't loaded
// any data and hasn't expressed readiness to serve.
var healthy int32

// SetGood sets the overall status of the server to be OK.
func SetGood() {
	atomic.StoreInt32(&healthy, 1)
	log.Println("Server health status GOOD")
}

// SetBad sets the overall status of the server to be BAD.
func SetBad() {
	atomic.StoreInt32(&healthy, 0)
	log.Println("Server health status BAD")
}

// Get returns whether the server is OK or not.
func Get() bool {
	return atomic.LoadInt32(&healthy) != 0
}
