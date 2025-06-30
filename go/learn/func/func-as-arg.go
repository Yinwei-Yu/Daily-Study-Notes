package main

import "strings"

//	"fmt"

func main() {
	// Indexfunc returns the index of the first rune in s that satisfies f.
	// If no rune satisfies f, IndexFunc returns -1.
	strings.IndexFunc("Hello", func(r rune) bool {
		return r == 'e'
	})
}


func FuncAsArg() {
	callBack(3, Add)
}

func Add(a,b int) int {
	return a + b
}

func callBack(y int,f func (int,int) int){
	f(y,2)
}

