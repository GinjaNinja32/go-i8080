package i8080

import (
	"fmt"
	"io"
)

type conio struct {
	in       io.Reader
	inStream chan byte

	out io.Writer
}

func (c *CPU) initConsole(conin io.Reader, conout io.Writer) {
	c.conio.in = conin
	c.conio.out = conout
	c.conio.inStream = make(chan byte, 100)
	go c.processInput()
}

func (c *CPU) conHasChar() bool {
	return len(c.conio.inStream) > 0
}

func (c *CPU) conPutChar(char uint8) {
	if char == 8 {
		fmt.Printf("\033[1D") // \033[1D")
		return
	}
	n, err := c.conio.out.Write([]byte{char})
	if err != nil {
		panic(err)
	} else if n != 1 {
		panic(n)
	}
}

func (c *CPU) conGetChar() uint8 {
	ch := <-c.conio.inStream
	if ch == 0x7F {
		// convert DEL to ^H
		return 0x08
	}
	return ch
}

func (c *CPU) processInput() {
	buf := [4096]uint8{}
	for {
		n, err := c.conio.in.Read(buf[:])
		if err == io.EOF {
			return
		} else if err != nil {
			panic(err)
		}

		for i := 0; i < n; i++ {
			c.conio.inStream <- buf[i]
		}
	}
}
