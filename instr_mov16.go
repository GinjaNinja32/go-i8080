package i8080

import ()

/*
	This file:

	LXI
	SHLD, LHLD
	POP, PUSH
	XTHL
	SPHL
	XCHG
*/

// LXI: 0x01, 0x11, 0x21, 0x31
func instrLXI(op uint8, c *CPU) uint64 {
	reg := insArg2(op)
	val := insArg16(c)

	switch reg {
	case R4:
		c.SP = val
	default:
		c.SetR16(reg, val)
	}
	return 10
}

// SHLD: 0x22
func instrSHLD(op uint8, c *CPU) uint64 {
	addr := insArg16(c)
	c.Write16(addr, c.HL())
	return 16
}

// LHLD: 0x2A
func instrLHLD(op uint8, c *CPU) uint64 {
	addr := insArg16(c)
	c.SetHL(c.Read16(addr))
	return 16
}

// POP: 0xC1, 0xD1, 0xE1, 0xF1
func instrPOP(op uint8, c *CPU) uint64 {
	reg := insArg2(op)
	val := c.Pop()

	switch reg {
	case BC:
		c.SetBC(val)
	case DE:
		c.SetDE(val)
	case HL:
		c.SetHL(val)
	case R4:
		c.SetPSW(val)
	}

	return 10
}

// PUSH: 0xC5, 0xD5, 0xE5, 0xF5
func instrPUSH(op uint8, c *CPU) uint64 {
	reg := insArg2(op)
	var val uint16

	switch reg {
	case BC:
		val = c.BC()
	case DE:
		val = c.DE()
	case HL:
		val = c.HL()
	case R4:
		val = c.PSW()
	}

	c.Push(val)

	return 11
}

// XTHL: 0xE3
func instrXTHL(op uint8, c *CPU) uint64 {
	newHL := c.Pop()
	c.Push(c.HL())
	c.SetHL(newHL)
	return 18
}

// SPHL: 0xF9
func instrSPHL(op uint8, c *CPU) uint64 {
	c.SP = c.HL()
	return 5
}

// XCHG: 0xEB
func instrXCHG(op uint8, c *CPU) uint64 {
	newHL := c.HL()
	c.SetHL(c.DE())
	c.SetDE(newHL)
	return 5
}
