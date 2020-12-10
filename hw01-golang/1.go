package main

import "fmt"

func Filter(p func(int) bool) func(...int) []int {
	return func(args ...int) []int {
		size := len(args)
		filteredOutput := make([]int, size, size*2)
		positives := 0
		for _, arg := range args {
			if p(arg) {
				positives += 1
				filteredOutput = append(filteredOutput, arg)
			}
		}
		return filteredOutput[len(filteredOutput)-positives:]
	}
}

func main() {
	odds := Filter(func(x int) bool { return x%2 == 1 })
    evens := Filter(func(x int) bool { return x%2 == 0 })

	fmt.Println(odds(1, 2, 3, 4, 5))   // [1 3 5]
    fmt.Println(evens(1, 2, 3, 4, 5))  // [2 4]
    fmt.Println(odds(6, 7, 8, 9, 10))  // [7 9]
    fmt.Println(evens(6, 7, 8, 9, 10)) // [6 8 10]
}