package i8080

import (
	"fmt"
)

const (
	B uint8 = iota
	C
	D
	E
	H
	L
	M
	A
)

type CPU struct {
	Memory    [65536]uint8
	Registers [8]uint8
	Flags     flags

	SP uint16
	PC uint16

	OutputStr string
}

func (c *CPU) Debug() string {
	mnem, args := c.disasmPC()

	disasm := fmt.Sprintf("ASM: %4x => %-8s%-8s", c.PC, mnem, args)

	regs := fmt.Sprintf("REGISTERS: BC=%02x %02x DE=%02x %02x HL=%02x %02x A=%02x SP=%04x PC=%04x",
		c.Registers[B], c.Registers[C], c.Registers[D], c.Registers[E], c.Registers[H], c.Registers[L], c.Registers[A], c.SP, c.PC)

	pc, sp, hl := "", "", ""
	for i := c.PC; i < c.PC+8; i++ {
		pc += fmt.Sprintf(" %02x", c.Memory[i&0xFFFF])
	}
	for i := c.SP; i < c.SP+8; i++ {
		sp += fmt.Sprintf(" %02x", c.Memory[i&0xFFFF])
	}
	for i := uint32(c.HL()); i < uint32(c.HL())+2; i++ {
		hl += fmt.Sprintf(" %02x", c.Memory[i&0xFFFF])
	}
	ptrs := fmt.Sprintf("POINTERS: [PC]=%s [SP]=%s [HL]=%s", pc, sp, hl)

	flags := fmt.Sprintf("FLAGS: %s", FlagsToString(uint8(c.Flags)))

	output := fmt.Sprintf("OUTPUT: %s$", c.OutputStr)

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", regs, ptrs, flags, output, disasm)
}

func FlagsToString(ff uint8) string {
	f := flags(ff)
	var Z, C, A, P, S string

	if (f & F_ZERO) != 0 {
		Z = "Z"
	} else {
		Z = "NZ"
	}

	if (f & F_CARRY) != 0 {
		C = "C"
	} else {
		C = "NC"
	}

	if (f & F_AUX_CARRY) != 0 {
		A = "A"
	} else {
		A = "NA"
	}

	if (f & F_PARITY) != 0 {
		P = "PE"
	} else {
		P = "PO"
	}

	if (f & F_SIGN) != 0 {
		S = "M"
	} else {
		S = "P"
	}

	unused := fmt.Sprintf("%08b", f)

	return fmt.Sprintf("%s %s %s %s %s (%s)", S, Z, A, P, C, unused)
}

func New() *CPU {
	return &CPU{
		Flags: F_BIT1_1,
	}
}

var n int = 0

func (c *CPU) Step() (cycles int) {
	op := c.Memory[c.PC]
	opFunc := c.dispatch(op)
	c.PC++
	cycles = opFunc(op, c)

	if c.PC == 0x0005 { // calling BDOS
		switch c.Registers[C] {
		case 9: // full msg?
			addr := c.DE()
			for c.Memory[addr] != '$' {
				fmt.Printf("%c", c.Memory[addr])
				addr++
			}
		case 2:
			fmt.Printf("%c", c.Registers[E])
		default:
			fmt.Printf("{%d,%d}", c.Registers[C], c.Registers[A])
		}
	}

	return
}

func (c *CPU) Output(s string) {
	//c.OutputStr += s
	fmt.Printf("%s", s)
}

func (c *CPU) Run() (cycles uint64) {
	for {
		cycles += uint64(c.Step())

		if c.PC == 0 {
			fmt.Printf("\n\n")
			return
		}
	}
}
