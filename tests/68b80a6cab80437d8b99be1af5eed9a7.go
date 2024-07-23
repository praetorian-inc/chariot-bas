// + build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/praetorian-inc/chariot-bas/endpoint"
)

func test() {
	commands := []string{
		"driverquery",
		"sc query type= driver",
	}

	outputLength := 0
	for _, cmd := range commands {
		command := exec.Command("cmd", "/C", cmd)
		output, err := command.Output()
		outputLength += len(output)

		if err != nil {
			fmt.Println("Failed to execute:", cmd, "Error:", err)
			endpoint.Stop(endpoint.Protected.Blocked)
			return
		}

		drivers := strings.Split(string(output), "\n")
		for _, driver := range drivers {
			fmt.Println(driver)
		}
	}

	if outputLength > 0 {
		endpoint.Stop(endpoint.Risk.Allowed)
	}

	endpoint.Stop(endpoint.Protected.Blocked)
}

func cleanup() {
}

func main() {
	endpoint.Start(test, cleanup)
}
