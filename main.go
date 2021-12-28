package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func terminal_ui() {
	if err := ui.Init(); err != nil {
		log.Fatal("Can't init termui")
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Text = "Prototype"
	p.SetRect(0, 0, 25, 5)

	ui.Render(p)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}

func termui_test() {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "List"
	l.Rows = []string{
		"[0] github.com/gizak/termui/v3",
		"[1] [你好，世界](fg:blue)",
		"[2] [こんにちは世界](fg:red)",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)

	ui.Render(l)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(l)
	}

}

func dirs_to_list(files []fs.FileInfo) []string {
	list := []string{}
	for i := range files {
		list = append(list, files[i].Name())
	}
	return list
}

func enter_dir(current_dir string, files []fs.FileInfo, new_dir int) (string, []fs.FileInfo) {
	fmt.Println(new_dir)
	fmt.Println(current_dir)
	current_dir += files[new_dir].Name() + "/"

	files, _ = ioutil.ReadDir(current_dir)

	return current_dir, files
}

func update_list(list *widgets.List, current_dir string) {
	list.Title = current_dir
	files, _ := ioutil.ReadDir(current_dir)
	list.Rows = dirs_to_list(files)
	list.ScrollTop()

}

func is_dir(current_dir []string, list *widgets.List) (bool, []string) {
	new_dir := slice_2_string(current_dir) + list.Rows[list.SelectedRow] + "/"
	_, err := ioutil.ReadDir(new_dir)
	if err != nil {
		return false, []string{"fail"}
	}
	current_dir = append(current_dir, (list.Rows[list.SelectedRow] + "/"))
	return true, current_dir
}

func go_back(current_dir []string) []string {
	new := current_dir[0 : len(current_dir)-1]
	return new
}

func slice_2_string(slice []string) string {
	var temp_string string
	for i := range slice {
		temp_string += slice[i]
	}
	return temp_string
}

func main() {
	current_dir := []string{"/", "home/", "tgallaher/"}
	files, _ := ioutil.ReadDir(slice_2_string(current_dir))

	if err := ui.Init(); err != nil {
		log.Fatal("Could not init UI!")
	}
	defer ui.Close()

	list := widgets.NewList()

	list.Title = slice_2_string(current_dir)

	list.Rows = dirs_to_list(files)

	list.TextStyle = ui.NewStyle(ui.ColorGreen)
	list.WrapText = false
	list.SetRect(0, 0, 50, 25)

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
		case "d":
			dir, temp_dir := is_dir(current_dir, list)
			if dir {

				current_dir = temp_dir
				update_list(list, slice_2_string(current_dir))
			}

		case "a":
			if len(current_dir)-1 != 0 {
				current_dir = go_back(current_dir)
				update_list(list, slice_2_string(current_dir))
			}

		}

		ui.Render(list)

	}
}
