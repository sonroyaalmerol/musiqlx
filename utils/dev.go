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
	// Check if templ is installed
	if !checkCommandInstalled("templ") {
		fmt.Println("Error: templ is not installed.")
		fmt.Println("go install github.com/a-h/templ/cmd/templ@latest")
		return
	}

	var wg sync.WaitGroup

	// Add the number of goroutines to wait for
	wg.Add(1)

	// Running "templ generate -watch" in a goroutine
	go func() {
		defer wg.Done()
		err := runCommand("templ", "generate", "--watch", "--proxy=http://localhost:3000", "--cmd=go run .")
		if err != nil {
			fmt.Println("Error running templ generate --watch:", err)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
}
