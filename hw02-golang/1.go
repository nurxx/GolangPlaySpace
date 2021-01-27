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

type Tasks struct {
	tasks []Task
}

func (t Tasks) Execute(x int) (int, error) {
	var err error
	if len(t.tasks) == 0 {
		return 0, errors.New("No tasks in pipeline")
	}
	for _, task := range t.tasks {
		x, err = task.Execute(x)
		if err != nil {
			return 0, err
		}
	}
	return x, nil
}

func Pipeline(tasks ...Task) Task {
	return Tasks{tasks}
}

func main() {
	if res, err := Pipeline(adder{50}, adder{60}).Execute(10); err != nil {
		fmt.Printf("The pipeline returned an error\n")
	} else {
		fmt.Printf("The pipeline returned %d\n", res)
	}
}
