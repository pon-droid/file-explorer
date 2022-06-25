package main

import (
	"io/fs"
	"io/ioutil"
	"sort"

	"github.com/agnivade/levenshtein"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func write_text(list *widgets.List, current_dir *[]string) {
	list.Title = ""
	ui.Render(list)
	for e := range ui.PollEvents() {
		switch e.ID {
		case "<Escape>":
			list.Title = slice_2_string(*current_dir)
			return
		case "<Backspace>", "<C-<Backspace>>":
			if len(list.Title) > 0 {
				list.Title = list.Title[:(len(list.Title) - 1)]
				ui.Render(list)
			}
		case "<Enter>":
			dir, temp_dir := is_dir(*current_dir, list)
			if dir {

				*current_dir = temp_dir
				update_list(list, slice_2_string(*current_dir))
			}
			return
		case "<Space>":
			list.Title = list.Title + " "
		default:
			list.Title = list.Title + e.ID
			ui.Render(list)
		}
		files, _ := ioutil.ReadDir(slice_2_string(*current_dir))
		list.Rows = filter(list.Title, files)
		list.ScrollTop()
		ui.Render(list)

	}

}

type fuzzy_search struct {
	name string
	dist int
}

type fuzz_list []fuzzy_search

func (f fuzz_list) Len() int { return len(f) }

func (f fuzz_list) Less(i, j int) bool { return f[i].dist < f[j].dist }

func (f fuzz_list) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func filter(input string, files []fs.FileInfo) []string {
	new_list := make(fuzz_list, len(files))
	for i := range files {
		dist := levenshtein.ComputeDistance(input, files[i].Name())
		new_list[i].name = files[i].Name()
		new_list[i].dist = dist
	}
	sort.Sort(new_list)
	string_slice := make([]string, len(new_list))
	for i := range new_list {
		string_slice[i] = new_list[i].name
	}
	return string_slice
}
