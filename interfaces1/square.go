package main

type square struct {
	sideLen float64
}

func (s square) getArea() float64 {
	return s.sideLen * s.sideLen
}
