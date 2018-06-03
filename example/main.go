package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/GinjaNinja32/go-i8080"
)

func main() {
	rawMode := exec.Command("/bin/stty", "raw", "-echo")
	rawMode.Stdin = os.Stdin
	_ = rawMode.Run()
	_ = rawMode.Wait()

	defer func() {
		stty := exec.Command("/bin/stty", "-raw", "echo")
		stty.Stdin = os.Stdin
		_ = stty.Run()
		_ = stty.Wait()
	}()

	cpm, err := ioutil.ReadFile("CPM22.bin")
	if err != nil {
		panic(err)
	}

	disks := []i8080.Disk{}

	for _, f := range os.Args[1:] {
		readOnly := false
		if strings.HasPrefix(f, "ro:") {
			readOnly = true
			f = strings.TrimPrefix(f, "ro:")
		}

		file, err := os.OpenFile(f, os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		disks = append(disks, i8080.Disk{
			Data:     file,
			ReadOnly: readOnly,
		})
	}

	c := i8080.New(os.Stdin, os.Stdout, cpm, disks)
	c.Run()
}
