package utils

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func checkCommandInstalled(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func Dev() {
	// Check if gow is installed
	if !checkCommandInstalled("gow") {
		fmt.Println("Error: gow is not installed.")
		fmt.Println("go install github.com/mitranim/gow@latest")
		return
	}

	// Check if templ is installed
	if !checkCommandInstalled("templ") {
		fmt.Println("Error: templ is not installed.")
		fmt.Println("go install github.com/a-h/templ/cmd/templ@latest")
		return
	}

	var wg sync.WaitGroup

	// Add the number of goroutines to wait for
	wg.Add(2)

	// Running "gow run main.go" in a goroutine
	go func() {
		defer wg.Done()
		err := runCommand("gow", "run", ".")
		if err != nil {
			fmt.Println("Error running gow:", err)
		}
	}()

	// Running "templ generate -watch" in a goroutine
	go func() {
		defer wg.Done()
		err := runCommand("templ", "generate", "-watch")
		if err != nil {
			fmt.Println("Error running templ generate -watch:", err)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
}
