package main

func Filter(p func(int) bool) func(...int) []int {
	return func(args ...int) []int {
		size := len(args)
		filteredOutput := make([]int, size, size * 2)
		positives := 0
		for _, arg := range args {
			if p(arg) {
				positives += 1	
				filteredOutput = append(filteredOutput, arg)
			}
		}
		return filteredOutput[len(filteredOutput) - positives:]
	}
}

func Mapper(f func(int) int) func(...int) []int {
	return func(args ...int) []int {
		size := len(args)
		computedOutput := make([]int, size, size * 2)
		for _, arg := range args {
			computedOutput = append(computedOutput, f(arg))
		}
		return computedOutput[size:]
	}
}

func Reducer(initial int, f func(int, int) int) func(...int) int {
	return func(args ...int) int {
		for _, arg := range args {
			initial = f(initial, arg)
		}
		return initial 
	}	
}

func MapReducer(initial int, mapper func(int) int, reducer func(int, int) int) func(...int) int {
	return func(args ...int) int {
		for _, arg := range args {
			initial = reducer(initial, mapper(arg))
		}
		return initial 
	}
}
