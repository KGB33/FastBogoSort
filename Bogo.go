package FastBogo

import (
	"io/ioutil"
	"log"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Sorts a splice of integers such that the the computational load will not exceed
// the given maxes.
// maxMem: Maximum Bytes of memory Allocated;
// maxCPU: Maximum Percent CPU Useage (0 to 100);
// delay: When one of the above values is exceed,
// how much time should we wait before calling garbage collection?
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
// 	maxMem: ~1GB,
// 	maxCPU: 80%,
// 	delay: 50ms,
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

// Uses the `/proc` file system to get real-time
// CPU stats. see `man 5 proc` for more info.
//
// Also see this SO post for more info
// https://stackoverflow.com/a/17783687
//
// NOTE: I'm assuming that the order is consistant.
func currCPUUse() int {
	rawData, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	data := string(rawData)
	cpuData := strings.Fields(strings.Split(data, "\n")[0]) // CPU is the 1st line
	if cpuData[0] != "cpu" {
		log.Fatal("Incorrectly parsed `/proc/stat`")
	}
	var total, idle float64
	for i, v := range cpuData[1:] {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Fatal(err)
		}
		total += val
		if i == 3 { // 4th field is idle
			idle = val
		}
	}
	return int(100 * ((total - idle) / total))
}

// For some reason the memory allocated can
// be greater than the total physical memory on
// the system. (Even with no swap)
// As such:
// TODO: Change this to use /proc/self/stat
// Much like the currCPUUse function
func currMemoryAlloc() int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(m.Alloc)
}
