package screen

import (
	"fmt"
	"os"
	"path/filepath"
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

// Traverses the file system backwards and forwards and update directory being dispayed
func (dir *Dir) UpdateDirectoryContents(direction Direction, ui *Ui) error {

	var currentPath, nextDir, positionPath string

	// Checks if there are an files or folders to open
	if dir.position >= 0 && dir.position < len(dir.files) {
		currentPath = dir.files[dir.position].path

	} else {
		currentPath = dir.path
	}

	if direction == Backwards {
		nextDir = filepath.Dir(dir.path)
		positionPath = dir.path
	} else {
		nextDir = currentPath
	}

	directoryPath, files, positionIndex, err := dir.ReadDir(nextDir, positionPath)

	if err != nil {
		return err
	}
	dir.path = directoryPath // directory path is the folder you are in at the end of the function
	dir.files = files
	dir.position = positionIndex // positionIndex represents the index of the file in the array

	fmt.Println("directorypath", directoryPath)

	if positionIndex > ui.ymax {
		dir.drawBeginning = positionIndex
		dir.highlightIndex = 0
	} else {
		dir.drawBeginning = 0
		dir.highlightIndex = positionIndex
	}

	return nil
}

func (dir *Dir) UpdateDirectoryPosition(direction Direction, ui *Ui) {
	if direction == Down && dir.position < len(dir.files)-1 {
		dir.position++
		if dir.highlightIndex < ui.ymax-2 {
			dir.highlightIndex++
		}
	} else if direction == Up && dir.position > 0 {
		dir.position--
		if dir.highlightIndex > 0 {
			dir.highlightIndex--
		}
	}

	dir.handleScroll(ui.ymax)

}

func (dir *Dir) handleScroll(height int) {

	if dir.position > height-3 && dir.highlightIndex == height-2 {
		dir.drawBeginning++
		dir.highlightIndex-- // Highlight index has reached the end and should do any further
	} else if dir.position < dir.drawBeginning && dir.highlightIndex == 0 {
		dir.drawBeginning--
	}

}
