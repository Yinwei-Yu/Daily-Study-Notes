package main

import (
	"fmt"
	"reflect"
)

// StructTag demonstrates how to use struct tags in Go
// It defines a struct with tags and prints the tag values.

type TagType struct {
	field1 bool   "An important answer"
	field2 string "The name of the thing"
	field3 int    "How much there are"
}

func main() {
	tt := TagType{true, "hello", 1}
	for i := range 3 {
		refTag(tt, i)
	}
}

func refTag(tt TagType, ix int) {
	// Use reflection to get the type of the struct
	t := reflect.TypeOf(tt)
	// Get the field by index
	field := t.Field(ix)
	// Print the field name and tag
	fmt.Printf("Field: %s, Tag: %s\n", field.Name, field.Tag)
}
