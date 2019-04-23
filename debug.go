package i8080

import (
	"fmt"
)

// Debug outputs a debug trace for the CPU
func (c *CPU) Debug() string {
	mnem, args := c.disasmPC()

	disasm := fmt.Sprintf("ASM %04x %02x => %-8s%-8s", c.PC, c.Memory[c.PC], mnem, args)

	regs := fmt.Sprintf("REG BC=%02x %02x DE=%02x %02x HL=%02x %02x A=%02x SP=%04x PC=%04x",
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
	ptrs := fmt.Sprintf("PTR [PC]=%s [SP]=%s [HL]=%s", pc, sp, hl)

	flags := fmt.Sprintf("FLG %s", FlagsToString(uint8(c.Flags)))

	return fmt.Sprintf("%s %s %s %s", regs, ptrs, flags, disasm)
}

// FlagsToString converts a flag register value `ff` to a string describing the state of the flags
func FlagsToString(ff uint8) string {
	f := flags(ff)
	var S, Z, P, C string

	var short string

	if (f & FlagSign) != 0 {
		S = "M"
		short += "S"
	} else {
		S = "P"
		short += "_"
	}

	if (f & FlagZero) != 0 {
		Z = " Z"
		short += "Z"
	} else {
		Z = "NZ"
		short += "_"
	}

	if (f & FlagBit5) != 0 {
		short += "x"
	} else {
		short += "_"
	}

	if (f & FlagAuxCarry) != 0 {
		short += "A"
	} else {
		short += "_"
	}

	if (f & FlagBit3) != 0 {
		short += "x"
	} else {
		short += "_"
	}

	if (f & FlagParity) != 0 {
		P = "PE"
		short += "P"
	} else {
		P = "PO"
		short += "_"
	}

	if (f & FlagBit1) != 0 {
		short += "x"
	} else {
		short += "_"
	}

	if (f & FlagCarry) != 0 {
		C = " C"
		short += "C"
	} else {
		C = "NC"
		short += "_"
	}

	return fmt.Sprintf("%s (%02x; %s %s %s %s)", short, f, S, Z, P, C)
}
