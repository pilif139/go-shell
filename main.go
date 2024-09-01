package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func execInput(input string) error {
	input = strings.TrimSuffix(input, "\r\n")

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return fmt.Errorf(red("path required"))
		}
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", append([]string{"/C", args[0]}, args[1:]...)...)
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func getPath() (string, error) {
	cmd := exec.Command("cmd", "/C", "cd")

	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	path := string(output)
	path = strings.TrimSuffix(path, "\r\n")
	path = filepath.Base(path)

	return path, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		path, err := getPath()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(path, " -> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\n~~Goodbye!~~")
				os.Exit(0)
			}
			fmt.Fprintln(os.Stderr, "error reading input:", err)
		}

		if err := execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, "Error executing command:", err)
		}
	}
}
