package main

import (
	"fmt"
	"os"

	"github.com/bailey30/fileSystemCLI/pkg/cli"
)

func main() {
	logFile, err := os.Create("error.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	os.Stderr = logFile

	dir, err := cli.InitDir("")
	if err != nil {
		return
	}

	filteredDirectory, err := cli.InitDir("")

	error := cli.NewError("")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	cli.InitUi(dir, filteredDirectory, error)

	normalLogFile, err := os.Create("logs.log")
	if err != nil {
		panic(err)
	}
	os.Stdout = normalLogFile
	fmt.Fprintln(os.Stdout, "hello2")
}
