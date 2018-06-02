package i8080

import (
	"fmt"
)

// Debug outputs a debug trace for the CPU
func (c *CPU) Debug() string {
	mnem, args := c.disasmPC()

	disasm := fmt.Sprintf("ASM %4x %2x => %-8s%-8s", c.PC, c.Memory[c.PC], mnem, args)

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
	var Z, C, A, P, S string

	if (f & FlagZero) != 0 {
		Z = " Z"
	} else {
		Z = "NZ"
	}

	if (f & FlagCarry) != 0 {
		C = " C"
	} else {
		C = "NC"
	}

	if (f & FlagAuxCarry) != 0 {
		A = " A"
	} else {
		A = "NA"
	}

	if (f & FlagParity) != 0 {
		P = "PE"
	} else {
		P = "PO"
	}

	if (f & FlagSign) != 0 {
		S = "M"
	} else {
		S = "P"
	}

	unused := fmt.Sprintf("%08b", f)

	return fmt.Sprintf("%s %s %s %s %s (%s)", S, Z, A, P, C, unused)
}
