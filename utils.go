package main

import (
	"io/fs"
	"io/ioutil"

	"github.com/gizak/termui/v3/widgets"
)

func dirs_to_list(files []fs.FileInfo) []string {
	list := make([]string, len(files))
	for i := range files {
		list[i] = files[i].Name()
	}
	return list
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
		return false, []string{"no"}
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
