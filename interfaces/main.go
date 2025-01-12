package main

import "fmt"

// Interfaces are implicit, we dont have to explicitly say all the other types are members of the interface

type bot interface {
	getGreeting() string
}

type englishBot struct {
	greeting string
}

type testBot string

type spanishBot struct {
	greeting string
}

func main() {
	eb := englishBot{greeting: "test"}
	sb := spanishBot{greeting: "test"}
	tb := testBot("test")
	printGreeting(eb)
	printGreeting(sb)
	printGreeting(tb)

}

// See here, you are not passing the struct (even though you are calling with struct type) 
// but the interface (all the types are members of the interface because they have a common function)

func printGreeting(b bot) { 
	// calling with bot interface
	fmt.Println(b.getGreeting())
}

// custom definition (with same function name) for each bot

func (eb englishBot) getGreeting() string {
	eb.greeting = "Hello!"
	return eb.greeting
}

func (sb spanishBot) getGreeting() string {
	sb.greeting = "Hola!"
	return sb.greeting
}
func (t testBot) getGreeting() string {
	return "Test!"
}
