package main

import "fmt"

type person struct {
	firstName string
	lastName  string
	contactInfo
}

type contactInfo struct {
	email string
	zip   int
}

func main() {
	alex := person{
		firstName: "Alex",
		lastName:  "Eiei",
		contactInfo: contactInfo{
			email: "example@domain.com",
			zip:   1150,
		},
	}

	alex.updateName("Mooping")
	alex.print()
}

func (p *person) updateName(newFirstName string) {
	p.firstName = newFirstName
}

func (p person) print() {
	fmt.Printf("%+v", p)
}
