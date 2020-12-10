package main

import "fmt"

func MapReducer(initial int, mapper func(int) int, reducer func(int, int) int) func(...int) int {
	return func(args ...int) int {
		for _, arg := range args {
			initial = reducer(initial, mapper(arg))
		}
		return initial
	}
}

func main() {
	powerSum := MapReducer(
        0,
        func(v int) int { return v * v },
        func(a, v int) int { return a + v },
    )

    fmt.Println(powerSum(1, 2, 3, 4)) // 30
    fmt.Println(powerSum(1, 2, 3, 4)) // 60
    fmt.Println(powerSum())           // 60
}


