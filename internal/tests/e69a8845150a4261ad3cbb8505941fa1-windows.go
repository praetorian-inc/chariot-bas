package main

import (
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/praetorian-inc/chariot-bas/internal/endpoint"
)

const (
	domainLength = 12
	domainCount  = 10
)

var tlds = []string{".ru", ".ir", ".cn", ".su"}

func generateDomain() string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	sb := strings.Builder{}
	for i := 0; i < domainLength; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	tld := tlds[rand.Intn(len(tlds))]
	return sb.String() + tld
}

func test() {
	domains := make([]string, domainCount)
	for i := 0; i < domainCount; i++ {
		domains[i] = generateDomain()
	}

	for _, domain := range domains {
		net.LookupHost(domain)
	}

	endpoint.Stop(endpoint.Risk.Allowed)
}

func cleanup() {
}

func main() {
	rand.Seed(time.Now().UnixNano())
	endpoint.Start(test, cleanup)
}
