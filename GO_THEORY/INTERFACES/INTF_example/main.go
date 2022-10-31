package main

import (
	"company"
	"person"
)

func main() {
	pers := person.Person{}
	comp := company.Company{}

	comp.Hire(pers) // мы передаём переменную типа Person в функцию, аргументом которой является переменная Worker!
}
