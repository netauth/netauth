// Package startup allows other packages to register startup tasks
// that will be run after early initialization completes.
package startup

// A Callback is registered in init(), and must not attempt to log or
// initialize.  They allow the order in which factories are called to
// be handled in the right order.
type Callback func()

var (
	callbacks []Callback
)

// RegisterCallback registers a callback for later execution.
func RegisterCallback(cb Callback) {
	callbacks = append(callbacks, cb)
}

// DoCallbacks executes all callbacks currently registered.
func DoCallbacks() {
	for _, cb := range callbacks {
		cb()
	}
}
