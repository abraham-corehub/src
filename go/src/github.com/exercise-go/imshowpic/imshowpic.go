package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	var img [][]uint8
	for x := 0; x < dx; x++ {
		var s []uint8
		for y := 0; y < dy; y++ {
			s = append(s, uint8(x*y))
		}
		img = append(img, s)
	}
	return img
}

func main() {
	pic.Show(Pic)
}
