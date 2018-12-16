package db

var (
	callbacks map[string]Callback
)

func init() {
	callbacks = make(map[string]Callback)
}

// RegisterCallback takes a callback name and handle and registers
// them for later calling.
func RegisterCallback(name string, c Callback) {
	if _, ok := callbacks[name]; ok {
		// Already here...
		return
	}
	callbacks[name] = c
}

// FireEvent fires an event to all callbacks.
func FireEvent(e Event) {
	for _, c := range callbacks {
		c(e)
	}
}
