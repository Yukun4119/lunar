package demo

import "fmt"

func demo() {
	a := 1
	if a == 1 {
		foo()
	} else {
		bar()
	}
}

func foo() {
	fmt.Print("hello")
}

func bar() {
	fmt.Print("world")
}
