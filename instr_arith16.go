package i8080

/*
	This file:

	INX, DAD, DCX
*/

// INX: 0x03, 0x13, 0x23, 0x33
func i_inx(op uint8, c *CPU) int {
	reg := insArg2(op)

	switch reg {
	case R4:
		c.SP = c.SP + 1
	default:
		c.SetR16(reg, c.GetR16(reg)+1)
	}

	return 5
}

// DAD: 0x09, 0x19, 0x29, 0x39
func i_dad(op uint8, c *CPU) int {
	reg := insArg2(op)

	a := c.HL()
	var b uint16

	switch reg {
	case R4:
		b = c.SP
	default:
		b = c.GetR16(reg)
	}

	if uint32(a)+uint32(b) > 65535 {
		c.Flags |= F_CARRY
	} else {
		c.Flags &= ^F_CARRY
	}

	c.SetHL(a + b)

	return 10
}

// DCX: 0x0B, 0x1B, 0x2B, 0x3B
func i_dcx(op uint8, c *CPU) int {
	reg := insArg2(op)

	switch reg {
	case R4:
		c.SP = c.SP - 1
	default:
		c.SetR16(reg, c.GetR16(reg)-1)
	}

	return 5
}
