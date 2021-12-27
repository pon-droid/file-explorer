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

}

func main() {
	current_dir := "/home/tgallaher/"
	files, _ := ioutil.ReadDir(current_dir)

	if err := ui.Init(); err != nil {
		log.Fatal("Could not init UI!")
	}
	defer ui.Close()

	list := widgets.NewList()

	list.Title = current_dir

	list.Rows = dirs_to_list(files)

	list.TextStyle = ui.NewStyle(ui.ColorGreen)
	list.WrapText = false
	list.SetRect(0, 0, 50, 25)

	ui.Render(list)
	for e := range ui.PollEvents() {
		update_list(list, current_dir)
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
		case "e":
			current_dir = current_dir + list.Rows[list.SelectedRow] + "/"
		}

		ui.Render(list)

	}
}
