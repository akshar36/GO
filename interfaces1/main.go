package main

import "fmt"

type shape interface {
	getArea() float64
}

func main() {
	sq := square{sideLen: 10}
	tr := triangle{base: 5, height: 10}
	printArea(sq)
	printArea(tr)
}

func printArea(s shape) {
	ar := s.getArea()
	fmt.Println(ar)
}
