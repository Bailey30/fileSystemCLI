package main

import (
	"fmt"
	"os"

	"github.com/bailey30/fileSystemCLI/pkg/screen"
)

func main() {
	logFile, err := os.Create("error.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	os.Stderr = logFile
	dirInstance := &screen.Dir{}
	dir, err := dirInstance.NewDir("")
	error := screen.NewError("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	screen.GetUi(dir, error)

	normalLogFile, err := os.Create("logs.log")
	if err != nil {
		panic(err)
	}
	os.Stdout = normalLogFile
	fmt.Fprintln(os.Stdout, "hello2")
}