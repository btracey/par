package par

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// For executes a for loop in parallel.
//    for i := 0; i < n; i++ {
//          f(i)
//    }
// The function f must be thread-safe. grain sets the number of iterations to
// perform per goroutine. Functions with faster evaluation times should use a
// larger grain size to amortize the cost of a goroutine.
func For(n, grain int, f func(i int)) {
	P := runtime.GOMAXPROCS(0)
	idx := uint64(0)
	var wg sync.WaitGroup
	wg.Add(P)
	for p := 0; p < P; p++ {
		go func() {
			for {
				start := int(atomic.AddUint64(&idx, uint64(grain))) - grain
				if start >= n {
					break
				}
				end := start + grain
				if end > n {
					end = n
				}
				for i := start; i < end; i++ {
					f(i)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
