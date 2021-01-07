package FastBogo

import (
	"math/rand"
	"runtime"
	"sort"
	"time"
)

func BogoSortDispatcher(splice []int) []int {
	/*
		This function takes a splice of ints and returns a sorted splice.

		This function calls the garbage collector to clear up compleated goroutines
		once ~1gb of memory is allocated.
	*/
	result := make(chan []int, 1)
	for len(result) != 1 {
		if currMemoryAlloc() > 1e9 {
			runtime.GC()
			runtime.Gosched()
		} else {
			go BogoSort(append([]int{}, splice...), result)
		}
	}
	return <-result
}

func BogoSort(splice []int, out chan []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(splice), func(i, j int) { splice[i], splice[j] = splice[j], splice[i] })
	if sort.IntsAreSorted(splice) {
		select {
		case out <- splice:
		default:
		}
	}
}

func currMemoryAlloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}
