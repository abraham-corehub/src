package main

import (
	"image"
	"image/color"

	"golang.org/x/tour/pic"
)

//Img is a custom type to generate a custom image
type Img struct{}

//ColorModel implements the method in the image.Image interface
func (i Img) ColorModel() color.Model {
	return color.RGBAModel
}

//Bounds implements the method in the image.Image interface
func (i Img) Bounds() image.Rectangle {

	return image.Rect(0, 0, 10, 10)

}

//At implements the method in the image.Image interface
func (i Img) At(x, y int) color.Color {

	return color.RGBA{uint8(x), uint8(y), 0, 255}

}

func main() {
	m := Img{}
	pic.ShowImage(m)
}
