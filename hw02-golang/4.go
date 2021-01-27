package main

import (
	"errors"
	"fmt"
)

type Task interface {
	Execute(int) (int, error)
}

type adder struct {
	augend int
}

func (a adder) Execute(addend int) (int, error) {
	result := a.augend + addend
	if result > 127 {
		return 0, fmt.Errorf("Result %d exceeds the adder threshold", a)
	}
	return result, nil
}

type ConcurrentTasks struct {
	reduce func(results []int) int
	tasks  []Task
}

func (t ConcurrentTasks) Execute(x int) (int, error) {
	result := make([]int, len(t.tasks))
	ch := make(chan int)
	if len(t.tasks) == 0 {
		return 0, errors.New("No tasks")
	}
	var err error
	for i, task := range t.tasks {
		go func(i int, task Task) {
			var res int
			res, err = task.Execute(x)
			ch <- res
			result[i] = res
		}(i, task)
	}

	for i, _ := range t.tasks {
		y := <-ch
		fmt.Printf("res = %d\n", y)
		fmt.Printf("i = %d\n", i)
	}
	return t.reduce(result), err
}

func ConcurrentMapReduce(reduce func(results []int) int, tasks ...Task) Task {
	return ConcurrentTasks{reduce, tasks}
}

func main() {
	reduce := func(results []int) int {
		smallest := 128
		for _, v := range results {
			if v < smallest {
				smallest = v
			}
		}
		return smallest
	}

	mr := ConcurrentMapReduce(reduce, adder{30}, adder{50}, adder{20})
	if res, err := mr.Execute(5); err != nil {
		fmt.Printf("We got an error!\n")
	} else {
		fmt.Printf("The ConcurrentMapReduce returned %d\n", res)
	}
}
