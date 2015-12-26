package grun

// #cgo pkg-config: glib-2.0
// #include "grun.h"
import "C"
import (
	"sync"
	"unsafe"
)

type done struct{}

var (
	mu      sync.Mutex // serializes calls to Run(f)
	runChan = make(chan func())
	runDone = make(chan done)
)

// Run runs f in the default glib main loop and waits for f to return.
// It can be called from any goroutine.
func Run(f func()) {
	mu.Lock()
	defer mu.Unlock()

	// add idle handler which will eventually invoke runFunc in the main loop thread,
	// assumed to be goroutine/thread-safe
	C.idle_add()

	runChan <- f // send f to main loop thread
	<-runDone    // wait for f to complete
}

// runFunc is the function that is being invoked by the glib main loop idle handler.
// It is assumed that this function runs in the same thread as the main loop.
// It returns 0 (gboolean FALSE) when the function received via the runChan channel is done
// or 1 (gboolean TRUE) otherwise.
// Note: this can busy-wait until f is received via the runChan channel.
//export run
func run(user_data unsafe.Pointer) (fin C.int) {
	fin = 1

	select {
	case f := <-runChan:
		f()
		runDone <- done{}
		fin = 0
	default:
	}

	return
}
