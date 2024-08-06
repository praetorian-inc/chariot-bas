// + build windows

package main

import (
	"os/exec"

	"github.com/praetorian-inc/chariot-bas/endpoint"
)

func test() {
	commands := []string{
		"vssadmin.exe delete shadows /all /quiet",
		"wmic.exe shadowcopy delete /nointeractive",
	}
	for _, cmd := range commands {
		command := exec.Command("cmd", "/C", cmd)
		err := command.Run()
		if err != nil {
			endpoint.Stop(endpoint.Protected.Blocked)
			return
		}
	}

	endpoint.Stop(endpoint.Risk.Allowed)
}

func cleanup() {
}

func main() {
	endpoint.Start(test, cleanup)
}
