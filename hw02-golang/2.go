package main

import (
	"errors"
	"fmt"
	"time"
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

type lazyAdder struct {
	adder
	delay time.Duration
}

func (la lazyAdder) Execute(addend int) (int, error) {
	time.Sleep(la.delay * time.Millisecond)
	return la.adder.Execute(addend)
}

type FastestTasks struct {
	tasks []Task
}

func (t FastestTasks) Execute(x int) (int, error) {
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
		}(i, task)
	}

	x = <-ch
	return x, err
}

func Fastest(tasks ...Task) Task {
	return FastestTasks{tasks}
}

func main() {
	f := Fastest(
		lazyAdder{adder{20}, 5000},
		lazyAdder{adder{50}, 100000},
		adder{41},
	)
	if res, err := f.Execute(1); err != nil {
		fmt.Printf("Fastes returned an error\n")
	} else {
		fmt.Printf("Fastest returned %d\n", res)
	}
}
