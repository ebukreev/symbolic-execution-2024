package main

type Overflow struct{}

func (o *Overflow) ShortOverflow(x int, y int) int {
	if y > 10 || y <= 0 {
		return 0
	}
	if x+y < 0 && x > 0 {
		return -1
	}
	return x + y
}
