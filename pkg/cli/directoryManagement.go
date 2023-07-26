package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type file struct {
	name string
	path string
}

type Dir struct {
	path           string
	files          []*file
	position       int // position in directory
	drawBeginning  int
	highlightIndex int
}

// create new instance of Dir - only being called once at the moment so path is always ""
func (dir *Dir) NewDir(path string) (*Dir, error) {

	var nextPath string

	if path == "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}
		nextPath = pwd
	} else {
		nextPath = path
	}

	normalLogFile, err := os.Create("logs.log")
	if err != nil {
		panic(err)
	}
	os.Stdout = normalLogFile
	fmt.Fprintln(os.Stdout, nextPath)

	_, files, positionIndex, err := dir.ReadDir(nextPath, nextPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	return &Dir{
		path:           nextPath,
		files:          files,
		position:       positionIndex,
		drawBeginning:  0,
		highlightIndex: 0,
	}, nil
}

func (d *Dir) ReadDir(path string, positionPath string) (string, []*file, int, error) {
	dir, err := os.Open(path)
	if err != nil {
		return "", nil, 0, err
	}
	defer dir.Close()
	contents, err := dir.Readdir(-1)
	if err != nil {
		return "", nil, 0, err
	}

	files := []*file{}
	positionIndex := 0

	for i, content := range contents {
		filePath := filepath.Join(path, content.Name())
		newFile := &file{name: content.Name(), path: filePath}
		// fmt.Println(newFile)
		files = append(files, newFile)
		if filePath == positionPath {
			positionIndex = i
		}
	}
	// fmt.Println("files", files)
	return path, files, positionIndex, nil
}

// Traverses the file system backwards and forwards and update directory being displayed
// activeDir is the original dir unless there is a search input
// Original dir is the only thing that changes
func (activeDir *Dir) UpdateDirectoryContents(direction Direction, ui *Ui, dir *Dir) error {

	ui.searchInput = ""

	var currentPath, nextDir, positionPath string

	// Checks if there are an files or folders to open
	if activeDir.position >= 0 && activeDir.position < len(activeDir.files) {
		currentPath = activeDir.files[activeDir.position].path

	} else {
		currentPath = activeDir.path
	}

	if direction == Backwards {
		nextDir = filepath.Dir(activeDir.path)
		positionPath = activeDir.path
	} else {
		nextDir = currentPath
	}

	directoryPath, files, positionIndex, err := activeDir.ReadDir(nextDir, positionPath)

	if err != nil {
		return err
	}
	dir.path = directoryPath // directory path is the folder you are in at the end of the function
	dir.files = files
	dir.position = positionIndex // positionIndex represents the index of the file in the array

	if positionIndex > ui.dirHeight-2 {
		dir.drawBeginning = positionIndex
		dir.highlightIndex = 0
	} else {
		dir.drawBeginning = 0
		dir.highlightIndex = positionIndex
	}

	return nil
}

// moves cursor up and down
func (dir *Dir) UpdateDirectoryPosition(direction Direction, ui *Ui) {
	if direction == Down && dir.position < len(dir.files)-1 {
		dir.position++
		if dir.highlightIndex < ui.dirHeight-2 {
			dir.highlightIndex++
		}
	} else if direction == Up && dir.position > 0 {
		dir.position--
		if dir.highlightIndex > 0 {
			dir.highlightIndex--
		}
	}

	dir.handleScroll(ui.dirHeight)

}

func (dir *Dir) handleScroll(height int) {
	if dir.position > height-3 && dir.highlightIndex == height-2 {
		dir.drawBeginning++
		dir.highlightIndex-- // Highlight index has reached the end and should do any further
	} else if dir.position < dir.drawBeginning && dir.highlightIndex == 0 {
		dir.drawBeginning--
	}

}

func (filteredDir *Dir) Filter(dir *Dir, searchInput string) {
	filtered := make([]*file, 0)

	for _, file := range dir.files {
		if strings.Contains(strings.ToLower(file.name), searchInput) {
			filtered = append(filtered, file)
		}
	}

	if len(filtered) > 0 {
		filteredDir.files = filtered
		filteredDir.path = filtered[0].path
	} else {
		filteredDir.files = nil
		filteredDir.path = ""
	}
	filteredDir.position = 0
	filteredDir.highlightIndex = 0
	filteredDir.drawBeginning = 0

}
