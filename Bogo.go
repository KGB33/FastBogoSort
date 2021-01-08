package FastBogo

import (
	"math/rand"
	"runtime"
	"sort"
	"time"
)

// Sorts a splice of integers suchthat the the computational load will not exceed
// the given maxes.
//
// 		- maxMem: Maximum Bytes of memory Allocated
// 		- maxCPU: Maximum Percent CPU Useage (0 to 100)
// 		- delay: When one of the above values is exceed, how much time should
// 				we wait before calling garbage collection?
//
// 	This function calls the garbage collector to clear up compleated goroutines
// 	once maxCPU or maxMem is exceed.
func CustomSortDispatcher(a []int, maxMem int, maxCPU int, delay time.Duration) []int {
	// Best case
	if sort.IntsAreSorted(a) {
		return a
	}
	// Realistic Case
	result := make(chan []int, 1)
	for len(result) != 1 {
		if currMemoryAlloc() > maxMem || currCPUUse() > maxCPU {
			time.Sleep(delay)
			runtime.GC()
			runtime.Gosched()
		} else {
			go Sort(append([]int{}, a...), result)
		}
	}
	return <-result
}

// 	Calls CustomSortDispatcher with the following defaults:
// 		- maxMem: ~1GB
// 		- maxCPU: 80%
// 		- delay: 50ms
func SortDispatcher(a []int) []int {
	return CustomSortDispatcher(a, 1e9, 80, 50*time.Millisecond)
}

// Shuffles the input. If the shuffled input is
// sorted, it pushes it to the channel.
func Sort(a []int, out chan []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	if sort.IntsAreSorted(a) {
		select {
		case out <- a:
		default:
		}
	}
}

func currCPUUse() int {
	// TODO
	return 0
}

func currMemoryAlloc() int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(m.Alloc)
}
