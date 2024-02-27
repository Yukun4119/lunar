package demo

import "fmt"

func demo() {
	foo()
	a := 1
	if a == 1 {
		bar()
	} else {
		foo()
	}
}

func foo() {
	bar()
	fmt.Print("hello")
}

func bar() {
	fmt.Print("world")
}
