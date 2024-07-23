// + build windo

package main

import (
	"crypto/rand"
	"math/big"
	"os/exec"

	"github.com/praetorian-inc/chariot-bas/endpoint"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generatePassword(length int) (string, error) {
	password := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))
	for i := range password {
		index, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		password[i] = charset[index.Int64()]
	}
	return string(password), nil
}

func test() {
	username := "securesvc"
	password, err := generatePassword(12)
	if err != nil {
		endpoint.Stop(endpoint.Errors.Unexpected)
		return
	}

	commands := [][]string{
		{"net", "user", username, password, "/add"},
		{"net", "localgroup", "administrators", username, "/add"},
	}

	for _, cmd := range commands {
		command := exec.Command(cmd[0], cmd[1:]...)
		if err := command.Run(); err != nil {
			endpoint.Stop(endpoint.Protected.Blocked)
			return
		}
	}

	endpoint.Stop(endpoint.Risk.Allowed)
}

func cleanup() {
	exec.Command("net", "user", "securesvc", "/delete").Run()
}

func main() {
	endpoint.Start(test, cleanup)
}
