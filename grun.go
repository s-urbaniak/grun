package grun

// #cgo pkg-config: glib-2.0
// #include "grun.h"
import "C"
import (
	"runtime"
	"sync"
)

var (
	runChan = make(chan func())
	runDone = make(chan struct{})
)

var once sync.Once

// Run runs f in the default glib main loop and waits for f to return.
// It can be called from different goroutines.
func Run(f func()) {
	runChan <- f
	<-runDone
}

//export runFunc
func runFunc() {
	select {
	case f := <-runChan:
		f()
		runDone <- struct{}{}
	default:
	}
}

func init() {
	// force main.main runs on main thread
	// see https://github.com/golang/go/wiki/LockOSThread
	runtime.LockOSThread()
	C.add_runner()
}
