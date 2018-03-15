package main

import (
	"image"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/themes/dark"
	"github.com/jackdreilly/biomorph"
)

func View(im image.Image, driver gxui.Driver) {
	theme := dark.CreateTheme(driver)
	img := theme.CreateImage()
	window := theme.CreateWindow(biomorph.ImageSize, biomorph.ImageSize, "Image viewer")
	texture := driver.CreateTexture(im, 1.0)
	img.SetTexture(texture)
	window.AddChild(img)
	window.OnClose(driver.Terminate)
}

func main() {
	c := biomorph.NewCreature(biomorph.NewTreeSpecies())
	c.SetValuesFromMap(map[string]float64{
		"num_branches":    2,
		"branch_length":   15,
		"num_gens":        4,
		"branch_angle":    0.5,
		"branch_increase": 1.0,
		"angle_increase":  1.0,
	})
	im := biomorph.DrawTreeCreature(c)
	gl.StartDriver(func(driver gxui.Driver) {
		View(im, driver)
	})
}
