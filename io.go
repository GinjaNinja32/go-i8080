package i8080

import (
	"io"
)

type conio struct {
	in       io.Reader
	inStream chan byte

	out io.Writer
}

func (c *CPU) initIO(conin io.Reader, conout io.Writer) {
	c.conio.in = conin
	c.conio.out = conout
	c.conio.inStream = make(chan byte, 100)
	go c.processInput()
}

func (c *CPU) ioHasChar() bool {
	return len(c.conio.inStream) > 0
}

func (c *CPU) ioPutChar(char uint8) {
	n, err := c.conio.out.Write([]byte{char})
	if err != nil {
		panic(err)
	} else if n != 1 {
		panic(n)
	}
}

func (c *CPU) ioGetChar() uint8 {
	return <-c.conio.inStream
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
