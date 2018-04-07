package health

import (
	"log"
	"sync/atomic"
)

// The server starts out unhealthy by default, since it hasn't loaded
// any data and hasn't expressed readiness to serve.
var healthy int32 = 0

func SetGood() {
	atomic.StoreInt32(&healthy, 1)
	log.Println("Server health status GOOD")
}

func SetBad() {
	atomic.StoreInt32(&healthy, 0)
	log.Println("Server health status BAD")
}

func Get() bool {
	return atomic.LoadInt32(&healthy) != 0
}
