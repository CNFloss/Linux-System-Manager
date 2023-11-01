package main

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func readDir(uid string) []string {
	files, err := os.ReadDir(uid)
	if err != nil {
		fyne.LogError("Failed to read directory "+uid, err)
		return nil
	}

	var c []string // Create a nil slice to append to
	for _, file := range files {
		c = append(c, filepath.Join(uid, file.Name()))
	}
	return c
}

func main() {
	a := app.New()
	win := a.NewWindow("File Tree")

	rootDir := "/" // This is the root directory for our tree

	childUIDs := func(uid widget.TreeNodeID) []widget.TreeNodeID {
		if uid == "" {
			uid = rootDir
		}
		children := readDir(uid)
		ids := make([]widget.TreeNodeID, len(children))
		for i, child := range children {
			ids[i] = widget.TreeNodeID(child)
		}
		return ids
	}

	createNode := func(branch bool) fyne.CanvasObject {
		if branch {
			return widget.NewLabel(" (folder)")
		}
		return widget.NewLabel(" (file)")
	}

	updateNode := func(uid widget.TreeNodeID, branch bool, node fyne.CanvasObject) {
		node.(*widget.Label).SetText(filepath.Base(uid))
	}

	isBranch := func(uid widget.TreeNodeID) bool {
		if uid == "" {
			return true
		}
		info, err := os.Stat(uid)
		if err != nil {
			return false
		}
		return info.IsDir()
	}

	tree := widget.NewTree(childUIDs, isBranch, createNode, updateNode)
	tree.OnSelected = func(uid widget.TreeNodeID) {
		// When a folder is clicked, we toggle its expanded state
		if tree.IsBranchOpen(uid) {
			tree.CloseBranch(uid)
		} else {
			tree.OpenBranch(uid)
		}
	}

	win.SetContent(container.NewBorder(nil, nil, nil, nil, tree))
	win.Resize(fyne.NewSize(400, 600))
	win.ShowAndRun()
}
