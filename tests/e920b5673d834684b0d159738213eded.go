// + build windows,amd64

package main

import (
	_ "embed"
	"fmt"
	"github.com/praetorian-inc/goffloader/src/coff"
	"strings"

	"github.com/praetorian-inc/chariot-bas/endpoint"
)

//go:embed static/bofs/whoami.x64.o.gz
var whoamiBytes []byte

func whoami_test() {
	decompressedBytes, _ := endpoint.Decompress(whoamiBytes)

	whoamiOutput, _ := coff.Load(decompressedBytes, make([]byte, 1))

	fmt.Println(whoamiOutput)

	if strings.Contains(whoamiOutput, "UserName") &&
		strings.Contains(whoamiOutput, "GROUP INFORMATION") &&
		strings.Contains(whoamiOutput, "Privilege Name") {
		endpoint.Stop(endpoint.Risk.Allowed)
	} else {
		endpoint.Stop(endpoint.Protected.Blocked)
	}

}

func whoami_cleanup() {
	return
}

func main() {
	endpoint.Start(whoami_test, whoami_cleanup)
}
