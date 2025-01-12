package main

import "fmt"

type person struct {
	firstName string
	lastName  string
	contact   contactInfo
}

type contactInfo struct {
	email string
	zip   int
}

func main() {
	var alex person
	alex = person{
		firstName: "Alex",
		lastName:  "Anderson",
		contact: contactInfo{
			email: "alex@gmail.com",
			zip:   123456,
		},
	}
	alex.updateName("Alexis")
	alex.print()

}

func (p *person) updateName(firstName string) {
	(*p).firstName = firstName
}

// func (p *person) updateNameUsingPointer(firstname string) {
// 	(*p).firstName = firstname
// }
func (p person) print() {
	fmt.Printf("%+v \n", p)
}
