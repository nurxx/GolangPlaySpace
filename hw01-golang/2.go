package main

import "fmt"

func Mapper(f func(int) int) func(...int) []int {
	return func(args ...int) []int {
		size := len(args)
		computedOutput := make([]int, size, size*2)
		for _, arg := range args {
			computedOutput = append(computedOutput, f(arg))
		}
		return computedOutput[size:]
	}
}

func main() {
	double := Mapper(func(a int) int { return 2 * a })
	decrement := Mapper(func(b int) int { return b - 1 })

    fmt.Println(double(1, 2, 3)) // [2, 4, 6]
	fmt.Println(double(4, 5, 6)) // [8, 10, 12]
	
	fmt.Println(decrement(1, 2, 3)) // [0, 1, 2]
    fmt.Println(decrement(4, 5, 6)) // [3, 4, 5]
}