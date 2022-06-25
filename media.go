package main

import (
	"image"
	"os"

	"github.com/gizak/termui/v3/widgets"
)

func display_image(current_dir []string, display *widgets.Image, list *widgets.List) {
	new_dir := slice_2_string(current_dir) + list.Rows[list.SelectedRow]
	file, err := os.Open(new_dir)
	if err != nil {
		return
	}

	image, _, _ := image.Decode(file)
	display.Image = image
	defer file.Close()

}
