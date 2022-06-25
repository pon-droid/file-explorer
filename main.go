package main

import (
	"encoding/json"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func settings_menu(list *widgets.List, config *settings) {
	list.Title = "Styles"
	options := []string{"Black", "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White"}
	list.Rows = options
	for e := range ui.PollEvents() {
		switch e.ID {
		case "e":
			config.Text_style = ui.StandardColors[list.SelectedRow]
		}
		ui.Render(list)
	}

}

func display_loop(config settings) {
	current_dir := config.Default_dir
	files, _ := ioutil.ReadDir(slice_2_string(current_dir))

	if err := ui.Init(); err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	list := widgets.NewList()

	list.Title = slice_2_string(current_dir)

	list.Rows = dirs_to_list(files)

	list.TextStyle = ui.NewStyle(config.Text_style)
	list.WrapText = config.Wrap_text
	list.SetRect(config.X1, config.Y1, config.X2, config.Y2)

	//Make image

	display := widgets.NewImage(nil)

	display.SetRect(config.X2, config.Y1, 200, config.Y2)
	display.BorderStyle = ui.NewStyle(config.Text_style)

	ui.Render(list)
	for e := range ui.PollEvents() {

		switch e.ID {
		case "w":
			list.ScrollUp()

		case "s":
			list.ScrollDown()

		case "<Home>":
			list.Title = "hello"
			list.Title = fmt.Sprint(list.SelectedRow)
		case "f":
			return
		case "d", "<Enter>":
			display.Image = nil
			dir, temp_dir := is_dir(current_dir, list)
			if dir {

				current_dir = temp_dir
				update_list(list, slice_2_string(current_dir))
			}

		case "a":
			display.Image = nil
			//Massive lagspike when navigating with image in buffer
			if len(current_dir)-1 != 0 {
				current_dir = go_back(current_dir)
				update_list(list, slice_2_string(current_dir))
			}
		case "t":
			write_text(list, &current_dir)
		case "e":

			display_image(current_dir, display, list)
		case "<Escape>":
			backup := list
			settings_menu(list, &config)
			list = backup
		}

		ui.Render(list)
		ui.Render(display)
	}
}

//Has to be capital to export to JSON

type settings struct {
	X1, Y1, X2, Y2 int
	Text_style     ui.Color
	Wrap_text      bool
	Default_dir    []string
}

func init_config() {
	x, err := read_config()
	if !err {
		default_config()
	} else {
		display_loop(x)
	}
}

func default_config() {
	var config settings
	config.X1, config.Y1 = 0, 0
	config.X2, config.Y2 = 50, 50
	config.Wrap_text = false
	config.Text_style = ui.ColorGreen
	config.Default_dir = append(config.Default_dir, "/")

	write_config(config)
}

func write_config(config settings) {
	file, _ := json.Marshal(config)

	_ = ioutil.WriteFile("config.json", file, 0644)
}

func read_config() (settings, bool) {
	var config settings
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config, false
	}

	json.Unmarshal([]byte(file), &config)
	return config, true
}
func main() {

	init_config()
}
