package main

import "fmt"

func Reducer(initial int, f func(int, int) int) func(...int) int {
	return func(args ...int) int {
		for _, arg := range args {
			initial = f(initial, arg)
		}
		return initial
	}
}

func main() {
	sum := Reducer(0, func(a, b int) int { return a + b })
	distract := Reducer(100, func(a, b int) int { return a - b })

    fmt.Println(sum(1, 2, 3))       // 6
    fmt.Println(sum(5))          // 11
	fmt.Println(sum(100, 101, 102)) // 314
	
	
	fmt.Println(distract(1,5)) // 94
	fmt.Println(distract(20,20,20)) // 34
	fmt.Println(distract(34)) // 0
	fmt.Println(distract()) // 0
}


