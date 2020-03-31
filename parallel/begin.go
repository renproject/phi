package parallel

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// ParBegin multiple functions onto background goroutines, where each function
// is run on its own goroutine. This function blocks until all given functions
// have been called and returned.
//
//	phi.ParBegin(
//		func() {
//			log.Println("goroutine #1")
//		},
//		func() {
//			log.Println("goroutine #2")
//		})
//
func ParBegin(fs ...func()) {
	var wg sync.WaitGroup
	for _, f := range fs {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(f)
	}
	wg.Wait()
}

// Begin multiple functions onto a number of background goroutines equal to the
// number of logical CPUs. This function blocks until all given functions have
// been called and returned.
//
//	phi.Begin(
//		func() {
//			// Assuming we only have 2x CPUs.
//			log.Println("goroutine #1 or #2")
//		},
//		func() {
//			log.Println("goroutine #1 or #2")
//		},
//		func() {
//			log.Println("goroutine #1 or #2")
//		},
//		func() {
//			log.Println("goroutine #1 or #2")
//		})
//
func Begin(fs ...func()) {
	// An atomic variable for indexing into the list of functions.
	i := uint64(0)

	// Create a wait group that expected "done" to be called once per CPU.
	numCPU := runtime.NumCPU()
	wg := sync.WaitGroup{}
	wg.Add(numCPU)

	// Spawn one goroutine per CPU, and call "done" on the wait group when the
	// goroutine returns.
	for cpu := 0; cpu < numCPU; cpu++ {
		go func(cpu int) {
			defer wg.Done()
			for {
				j := atomic.AddUint64(&i, 1)
				if j > uint64(len(fs)) {
					return
				}
				fs[j-1]()
			}
		}(cpu)
	}
	wg.Done()
}
