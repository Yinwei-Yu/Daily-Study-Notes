package main

//closeure is also known as anonymous function

import (
    "fmt"
)

func anonymousFunc() {
	// Anonymous function
	for i := range 4 {
		// Anonymous function that takes an int and prints it
			g := func(i int) { fmt.Printf("%d ", i) }
			g(i)
			fmt.Printf(" - g is of type %T and has value %v\n", g, g)
		}
}

// Closure is a function that captures the surrounding state

func closure() {
	var f = Adder()
	// value of x is saved in the closure
	fmt.Print(f(1), " - ") //1
	fmt.Print(f(20), " - ") //21
	fmt.Print(f(300)) //321
}

func Adder() func(int) int {
	var x int
	return func(delta int) int {
		x += delta
		return x
	}
}