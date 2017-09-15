package i8080

/*
	This file:

	RNZ, RNC, RPO, RP, RZ, RC, RPE, RM
	JNZ, JNC, JPO, JP, JZ, JC, JPE, JM
	CNZ, CNC, CPO, CP, CZ, CC, CPE, CM

	RET
	JMP
	CALL

	PCHL
	RST
*/

func insCond(op uint8, c *CPU) bool {
	switch insArg3b(op) {
	case 0: // NZ
		return (c.Flags & F_ZERO) == 0
	case 1: // Z
		return (c.Flags & F_ZERO) != 0
	case 2: // NC
		return (c.Flags & F_CARRY) == 0
	case 3: // C
		return (c.Flags & F_CARRY) != 0
	case 4: // PO
		return (c.Flags & F_PARITY) == 0
	case 5: // PE
		return (c.Flags & F_PARITY) != 0
	case 6: // P
		return (c.Flags & F_SIGN) == 0
	case 7: // M
		return (c.Flags & F_SIGN) != 0
	default:
		panic("impossible")
	}
}

// RNZ, RNC, RPO, RP, RZ, RC, RPE, RM
// 0xC0, 0xC8, 0xD0, 0xD8, 0xE0, 0xE8, 0xF0, 0xF8
func i_cond_ret(op uint8, c *CPU) uint64 {
	if insCond(op, c) {
		c.PC = c.Pop()
		return 11
	} else {
		return 5
	}
}

// RET: 0xC9, 0xD9
func i_ret(op uint8, c *CPU) uint64 {
	c.PC = c.Pop()
	return 10
}

// JNZ, JNC, JPO, JP, JZ, JC, JPE, JM
// 0xC2, 0xCA, 0xD2, 0xDA, 0xE2, 0xEA, 0xF2, 0xFA
func i_cond_jmp(op uint8, c *CPU) uint64 {
	addr := insArg16(c)
	if insCond(op, c) {
		c.PC = addr
	}
	return 10
}

// JMP: 0xC3, 0xCB
func i_jmp(op uint8, c *CPU) uint64 {
	c.PC = insArg16(c)
	return 10
}

// CNZ, CNC, CPO, CP, CZ, CC, CPE, CM
// 0xC4, 0xCC, 0xD4, 0xDC, 0xE4, 0xEC, 0xF4, 0xFC
func i_cond_call(op uint8, c *CPU) uint64 {
	addr := insArg16(c)
	if insCond(op, c) {
		c.Push(c.PC)
		c.PC = addr
		return 17
	}
	return 11
}

// CALL: 0xCD, 0xDD, 0xED, 0xFD
func i_call(op uint8, c *CPU) uint64 {
	addr := insArg16(c)
	c.Push(c.PC)
	c.PC = addr

	return 17
}

// PCHL: 0xE9
func i_pchl(op uint8, c *CPU) uint64 {
	c.PC = c.HL()
	return 5
}

// RST: 0xC7, 0xCF, 0xD7, 0xDF, 0xE7, 0xEF, 0xF7, 0xFF
func i_rst(op uint8, c *CPU) uint64 {
	site := uint16(insArg3b(op))
	addr := site * 0x8

	c.Push(c.PC)
	c.PC = addr
	return 11
}
