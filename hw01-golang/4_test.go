package main

import (
	"math"
	"testing"
)

func TestMapReducer(t *testing.T) {
	tests := []struct {
		description string
		mapFunc     func(int) int
		reduceFunc  func(int, int) int
		initial     int
		args        []int
		expected    int
	}{
		{
			description: "sqrtMap summer",
			mapFunc: func(n int) int {
				return int(math.Sqrt(float64(n)))
			},
			reduceFunc: func(a, b int) int {
				return a + b
			},
			initial:  0,
			args:     []int{1, 4, 9, 16, 25, 36, 49, 64},
			expected: 1 + 2 + 3 + 4 + 5 + 6 + 7 + 8,
		},
		{
			description: "biggest left shift",
			mapFunc: func(n int) int {
				return n << 1
			},
			reduceFunc: func(a, b int) int {
				if a > b {
					return a
				}
				return b
			},
			initial:  0,
			args:     []int{5, 21, 10, 25, 22, 18, 18, 19, 23},
			expected: 50,
		},
		{
			description: "left shift all the way",
			mapFunc: func(n int) int {
				return n
			},
			reduceFunc: func(a, b int) int {
				return a << uint(b)
			},
			initial:  42,
			args:     []int{1, 2, 3, 4, 5, 6},
			expected: 88080384,
		},
		{
			description: "initial is used first",
			mapFunc: func(n int) int {
				return n * n
			},
			reduceFunc: func(a, b int) int {
				return a * b
			},
			initial:  20,
			args:     []int{},
			expected: 20,
		},
		{
			description: "left foldedness",
			mapFunc: func(n int) int {
				return n * n
			},
			reduceFunc: func(a, b int) int {
				return b / a
			},
			initial:  1,
			args:     []int{2, 4, 8, 16, 32, 64},
			expected: 64,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			mreducer := MapReducer(test.initial, test.mapFunc, test.reduceFunc)
			actual := mreducer(test.args...)
			if test.expected != actual {
				t.Errorf("Expected %d but got %d", test.expected, actual)
			}
		})
	}
}