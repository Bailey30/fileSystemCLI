package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path"
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
	topOffset      int
}

func NewDir() *Dir {
	return &Dir{
		path:           "path",
		files:          []*file{},
		position:       0,
		drawBeginning:  0,
		highlightIndex: 0,
		topOffset:      1,
	}
}

// create new instance of Dir - only being called once at the moment so path is always ""
func InitDir(path string) (*Dir, error) {
	dir := NewDir()

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
		highlightIndex: dir.topOffset,
		topOffset:      1,
	}, nil
}

func (d *Dir) ReadDir(path string, positionPath string) (string, []*file, int, error) {
	dir, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
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
		files = append(files, newFile)
		if filePath == positionPath {
			positionIndex = i
		}
	}

	return path, files, positionIndex, nil
}

// Traverses the file system backwards and forwards and update directory being displayed
// activeDir is the original dir unless there is a search input
// Original dir is the only thing that changes
func (activeDir *Dir) Traverse(direction Direction, ui *Ui, dir *Dir) error {

	ui.searchInput = ""

	var pathToRead, positionPath string

	if direction == Refresh {
		pathToRead = activeDir.path
	} else if direction == Jump {
		pathToRead = activeDir.path
	} else if direction == Backwards {
		pathToRead = filepath.Dir(activeDir.path)
		positionPath = activeDir.path // needed to get the location of the current folder in the folder you move back to
	} else if direction == Forwards {
		// Checks if there are an files or folders to open
		if activeDir.position >= 0 && activeDir.position < len(activeDir.files) {
			pathToRead = activeDir.files[activeDir.position].path
		} else {
			// if there are no files just stay in the same spot
			pathToRead = ""
		}
	}

	directoryPath, files, positionIndex, err := activeDir.ReadDir(pathToRead, positionPath)

	if err != nil {
		fmt.Println("directorypatherr", err)
	}

	if err != nil {
		return err
	}

	dir.path = directoryPath // directory path is the folder you are in at the end of the function
	dir.files = files
	dir.position = positionIndex // positionIndex represents the index of the file in the array

	if positionIndex > ui.dirHeight-3 {
		dir.drawBeginning = positionIndex - (ui.dirHeight / 2) + dir.topOffset
		dir.highlightIndex = 0 + (ui.dirHeight / 2)
	} else {
		dir.drawBeginning = 0
		dir.highlightIndex = positionIndex + dir.topOffset
	}

	return nil
}

// moves cursor up and down
func (dir *Dir) UpDown(direction Direction, ui *Ui) {
	if direction == Down && dir.position < len(dir.files)-1 {
		// director position cannot go higher than the length of the directory
		dir.position++
		if dir.highlightIndex <= ui.dirHeight {
			// highlight index cannot go higher than the height of the screen
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
	//
	if dir.position > height-dir.topOffset-1-1 && dir.highlightIndex == height-dir.topOffset {
		dir.drawBeginning++
		dir.highlightIndex-- // Highlight index has reached the end and should not go any further
	} else if dir.position < dir.drawBeginning+dir.topOffset && dir.highlightIndex == 0 {
		dir.drawBeginning--
		dir.highlightIndex++
	}
}

func (filteredDir *Dir) Filter(dir *Dir, searchInput string) {
	filtered := make(chan *file)

	go func() {
		for _, file := range dir.files {
			if strings.Contains(strings.ToLower(file.name), searchInput) {
				filtered <- file
			}
		}
		close(filtered)
	}()

	filteredFiles := []*file{}
	for file := range filtered {
		filteredFiles = append(filteredFiles, file)
	}

	if len(filteredFiles) > 0 {
		filteredDir.files = filteredFiles
		// filteredDir.path = filtered[0].path
	} else {
		filteredDir.files = nil
	}

	filteredDir.path = dir.path
	filteredDir.position = 0
	filteredDir.highlightIndex = 0 + dir.topOffset
	filteredDir.drawBeginning = 0

}

func JumpToDirectory(ui *Ui, dir *Dir) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	path := path.Join(homedir, "/", ui.searchInput[2:])

	consistent := filepath.ToSlash(path)
	dir.path = consistent

	jumpErr := dir.Traverse(Jump, ui, dir)
	if jumpErr != nil {
		fmt.Println(jumpErr.Error())
		ui.searchInput = "/"
		return
	}
}

func (dir *Dir) create(name string, ui *Ui) {
	split := strings.Split(name, ".")

	newPath := path.Join(dir.path, name)

	if len(split) == 1 {
		err := os.Mkdir(newPath, 0755)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, err := os.Create(newPath)
		if err != nil {
			fmt.Println(err)
		}
	}

	ui.searchInput = ""
	dir.Traverse(Refresh, ui, dir)
}

func (dir *Dir) Delete(ui *Ui) {
	ui.confirm = false
	ui.searchInput = "Confirm delete y/n"
}

func (dir *Dir) ConfirmDelete(ui *Ui) {
	filename := dir.files[dir.position].path

	err := os.Remove(filename)
	if err != nil {
		fmt.Println(err)
	}

	ui.confirm = true
	ui.searchInput = ""
	dir.Traverse(Refresh, ui, dir)
}

func (dir *Dir) CancelDelete(ui *Ui) {
	ui.confirm = true
	ui.searchInput = ""
}

func (dir *Dir) Open(window string) {
	filename := dir.files[dir.position].path
	splitname := strings.Split(filename, ".")
	if len(splitname) > 1 {
		var args []string
		if window == "new" {
			args = []string{"--new-window", filename}
		} else {
			args = []string{filename}
		}
		cmd := exec.Command("code", args...)
		fmt.Println(cmd)
		if err := cmd.Start(); err != nil {
			fmt.Println(err)
			return
		}
		if err := cmd.Wait(); err != nil {
			fmt.Println(err)
		}
	}
}
