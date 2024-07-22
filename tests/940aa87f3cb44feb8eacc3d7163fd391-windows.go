package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/praetorian-inc/chariot-bas/endpoint"
)

var (
	tempDir = os.TempDir()
	sam     = filepath.Join(tempDir, "SAM")
	system  = filepath.Join(tempDir, "SYSTEM")
)

func test() {
	cmd := exec.Command("reg", "save", "HKLM\\SAM", sam)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to save SAM registry hive: %v", err)
		endpoint.Stop(endpoint.Protected.Blocked)
		return
	}

	cmd = exec.Command("reg", "save", "HKLM\\SYSTEM", system)
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to save SYSTEM registry hive: %v", err)
		endpoint.Stop(endpoint.Protected.Blocked)
		return
	}

	endpoint.Stop(endpoint.Risk.Allowed)
}

func cleanup() {
	if err := os.Remove(sam); err != nil {
		log.Printf("Failed to delete SAM file: %v", err)
	}
	if err := os.Remove(system); err != nil {
		log.Printf("Failed to delete SYSTEM file: %v", err)
	}
}

func main() {
	endpoint.Start(test, cleanup)
}
