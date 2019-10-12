package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	var img [][]uint8
	for x := 0; x < dx; x++ {
		var arr []uint8
		for y := 0; y < dy; y++ {
			arr = append(arr, uint8((x+y)/2))
		}

		img = append(img, arr)
	}
	return img
}

func main() {

	pic.Show(Pic)

}
