package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/themes/dark"
)

type item struct {
	Image []byte `json:"image"`
	Label int    `json:"label"`
}

func drawImageOnScreen(numberImage []byte) func(driver gxui.Driver) {
	return func(driver gxui.Driver) {
		gray := image.NewGray(image.Rect(0, 0, 28, 28))
		copy(gray.Pix, numberImage)
		rgba := image.NewRGBA(image.Rect(0, 0, gray.Bounds().Dx(), gray.Bounds().Dy()))
		draw.Draw(rgba, rgba.Bounds(), gray, gray.Bounds().Min, draw.Src)
		theme := dark.CreateTheme(driver)
		img := theme.CreateImage()
		window := theme.CreateWindow(28, 28, "MNIST")
		texture := driver.CreateTexture(rgba, 1.0)
		img.SetTexture(texture)
		window.AddChild(img)
		window.OnClose(driver.Terminate)
	}
}

func main() {
	fmt.Println("Reading the MNIST digit data ...")
	zipReader, err := zip.OpenReader("data/mnist_handwritten_test.json.zip")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer zipReader.Close()
	fmt.Println("Parsing ...")
	jsonReader, _ := zipReader.File[0].Open()
	raw, err := ioutil.ReadAll(jsonReader)
	var items []item
	err = json.Unmarshal(raw, &items)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Displaying ...")
	gl.StartDriver(drawImageOnScreen(items[0].Image))
}
