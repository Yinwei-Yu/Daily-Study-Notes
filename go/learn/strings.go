package main

import (
	"fmt"
	"strings"
)

func strings_test() {
	var str = "This is a string"
	fmt.Printf("T/F? Does the string \"%s\" have prefix %s?", str,"This")
	fmt.Printf("%t\n", strings.HasPrefix(str,"This"))
	
	strings.Contains(str, "hello")
	strings.LastIndex(str,"nihao")
	strings.IndexRune(str, rune('k'))
	
}
