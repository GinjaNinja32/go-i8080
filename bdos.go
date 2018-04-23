package i8080

import (
	"bufio"
	"fmt"
	"os"
)

var scratch = []byte{0}

type bdos struct {
	input *bufio.Reader
}

func (c *CPU) initBDOS() {
	c.bdos = bdos{
		input: bufio.NewReader(os.Stdin),
	}

	c.Memory[0x0005] = 0xDD // BDOS
	c.Memory[0x0006] = 0xC9 // RET
}

func instrBDOS(op uint8, c *CPU) uint64 {
	switch c.Registers[C] {
	case 0: // TERMCPM
		// TODO
	case 1: // READ
		n, err := c.bdos.input.Read(scratch)
		if n != 1 || err != nil {
			panic(fmt.Sprintf("failed to read: %d %s", n, err))
		}

		c.Registers[A] = scratch[0]
		c.Registers[L] = scratch[0]
	case 2: // WRITE
		fmt.Printf("%c", c.Registers[E])

	case 6: // RAWIO
		ready := c.bdos.input.Buffered() != 0

		switch c.Registers[E] {
		case 0xFF:
			if ready {
				n, err := c.bdos.input.Read(scratch)
				if n != 1 || err != nil {
					panic(fmt.Sprintf("failed to read: %d %s", n, err))
				}
				c.Registers[A] = scratch[0]
			} else {
				c.Registers[A] = 0x00
			}
		case 0xFE:
			if ready {
				c.Registers[A] = 0xFF
			} else {
				c.Registers[A] = 0x00
			}
		case 0xFD:
			n, err := c.bdos.input.Read(scratch)
			if n != 1 || err != nil {
				panic(fmt.Sprintf("failed to read: %d %s", n, err))
			}
			c.Registers[A] = scratch[0]
		case 0xFC:
			if ready {
				buf, err := c.bdos.input.Peek(1)
				if err != nil {
					panic(fmt.Sprintf("failed to read: %s", err))
				}
				c.Registers[A] = buf[0]
			} else {
				c.Registers[A] = 0x00
			}
		}
	case 8: // set I/O byte
		// nothing
	case 9: // WRITESTR
		addr := c.DE()
		for c.Memory[addr] != '$' {
			fmt.Printf("%c", c.Memory[addr])
			addr++
		}
	case 11: // STAT
		ready := c.bdos.input.Buffered() != 0
		if ready {
			c.Registers[A] = 0xFF
		} else {
			c.Registers[A] = 0x00
		}
	case 12: // BDOSVER
		mchType := uint8(0x00)
		version := uint8(0x22)

		c.Registers[B] = mchType
		c.Registers[H] = mchType
		c.Registers[A] = version
		c.Registers[L] = version
	default:
		fmt.Printf("{%d,%d}", c.Registers[C], c.Registers[A])
	}

	return 0
}
