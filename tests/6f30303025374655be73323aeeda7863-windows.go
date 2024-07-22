package main

import (
	"os/exec"

	"github.com/praetorian-inc/chariot-bas/endpoint"
)

func test() {
	commands := []string{
		"nltest /domain_trusts",
		"net group \"Domain Admins\" /domain",
		"net group \"Enterprise Admins\" /domain",
		"net user /domain",
		"systeminfo",
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
