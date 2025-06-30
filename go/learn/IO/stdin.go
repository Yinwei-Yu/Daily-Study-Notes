package main

import "fmt"

var (
	firstName, lastName string
	i                   int
	f                   float64
	s                   string
	input               = "56.12 / 5212 / Go"
	format              = "%f / %d / %s"
)

func main() {
	fmt.Println("Please enter your first name, last name")
	fmt.Scanln(&firstName, &lastName)
	fmt.Printf("Hi %s %s \n", firstName, lastName)
	fmt.Sscanf(input, format, &f, &i, &s)
	fmt.Println("From the string we read: ", f, i, s)
}
