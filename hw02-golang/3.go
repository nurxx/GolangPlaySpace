package main

import (
	"errors"
	"fmt"
	"time"
)

type Task interface {
	Execute(int) (int, error)
}

type lazyAdder struct {
	adder
	delay time.Duration
}

func (la lazyAdder) Execute(addend int) (int, error) {
	time.Sleep(la.delay * time.Millisecond)
	return la.adder.Execute(addend)
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

type TimedTask struct {
	task    Task
	timeout time.Duration
}

func (t TimedTask) Execute(x int) (int, error) {
	ch := make(chan int, 1)
	var err error
	var res int
	go func(t TimedTask) {
		res, err = t.task.Execute(x)
		ch <- res
	}(t)

	select {
	case res := <-ch:
		return res, nil
	case <-time.After(t.timeout):
		err = errors.New("timeout!")
		return 0, err
	}

}

func Timed(task Task, timeout time.Duration) Task {
	return TimedTask{task, timeout}
}

func main() {
	r1, e1 := Timed(lazyAdder{adder{20}, 50}, 2*time.Millisecond).Execute(2)
	if e1 != nil {
		fmt.Printf("Timed e1 returned an error...\n")
	} else {
		fmt.Printf("Timed r1 returned %d...\n", r1)
	}
	r2, e2 := Timed(lazyAdder{adder{20}, 50}, 300*time.Millisecond).Execute(2)
	if e2 != nil {
		fmt.Printf("Timed e2 returned an error...\n")
	} else {
		fmt.Printf("Timed r2 returned %d...\n", r2)
	}
}
