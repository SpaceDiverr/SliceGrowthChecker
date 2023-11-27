package slicegrowthchecker

import (
	"fmt"
	"sync"
)

// Function to find out slice capacity growth for your current version of Go
//
// Gets executed concurrently.
//
// SliceGrowingForType generates a slice of type T and appends elements to it up to the specified breakpoint.
//
// In the stdout, it will print the capacity of the slice before and after whether it has changed.
//
// T: the type of elements in the slice
// breakpoint: the number of iterations to perform.
//
// Return type: none
func SliceGrowthMetricsForType[T any](baseCap, breakpoint uint) {

	if baseCap >= breakpoint {
		panic("slicegrowthchecker: 'baseCap' must be < 'breakpoint'")
	}

	var (
		wg       *sync.WaitGroup
		mu       *sync.Mutex

		s        []T
		toAppend T

		bp, bc = int(breakpoint), int(baseCap)
	)
	wg = &sync.WaitGroup{}
	mu = &sync.Mutex{}
	s = make([]T, bc)

	fmt.Printf("Type considered: %T\n\n", toAppend)

	for i := 0; i < bp; i++ {
		wg.Add(1)
		go appendAndLog(&s, toAppend, wg, mu)
		wg.Wait()
	}
}

func appendAndLog[T any](s *[]T, toAppend T, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	capBefore := cap(*s)
	mu.Lock()
	*s = append(*s, toAppend)
	mu.Unlock()
	capAfter := cap(*s)
	if capBefore != capAfter {
		fmt.Println("___________________")
		fmt.Printf("capBefore:  %v\n", capBefore)
		fmt.Printf("capAfter:   %v\n", capAfter)
		fmt.Printf("growthRate (capBefore -> capAfter): %.4f\n", 1.0+(float64(capAfter)-float64(capBefore))/float64(capBefore))
	}

}
