package i8080

/*
	This file:

	INR, DCR
	RLC, RAL, DAA, STC, RRC, RAR, CMA, CMC
	ADD, ADC, SUB, SBB, ANA, XRA, ORA, CMP
	ADI, ACI, SUI, SBI, ANI, XRI, ORI, CPI
*/

// INR: 0x04, 0x0C, 0x14, 0x1C, 0x24, 0x2C, 0x34, 0x3C
func instrINR(op uint8, c *CPU) uint64 {
	reg := insArg3b(op)
	insSetreg8(c, reg, flaggedAdd8(c, insGetreg8(c, reg), 1, FlagCarry))
	if reg == M {
		return 10
	}
	return 5
}

// DCR: 0x05, 0x0D, 0x15, 0x1D, 0x25, 0x2D, 0x35, 0x3D
func instrDCR(op uint8, c *CPU) uint64 {
	reg := insArg3b(op)

	result := flaggedSub8(c, insGetreg8(c, reg), 1, FlagCarry)
	insSetreg8(c, reg, result)

	if (result & 0xf) != 0xf {
		c.Flags |= FlagAuxCarry
	} else {
		c.Flags &= ^FlagAuxCarry
	}

	if reg == M {
		return 10
	}
	return 5
}

// RLC: 0x07
func instrRLC(op uint8, c *CPU) uint64 {
	if (c.Registers[A] & 0x80) != 0 {
		c.Flags |= FlagCarry
	} else {
		c.Flags &= ^FlagCarry
	}

	return instrRAL(op, c)
}

// RAL: 0x17
func instrRAL(op uint8, c *CPU) uint64 {
	a := c.Registers[A]
	cy := (a & 0x80) != 0

	c.Registers[A] = a << 1

	if (c.Flags & FlagCarry) != 0 {
		c.Registers[A] |= 0x01
	}

	if cy {
		c.Flags |= FlagCarry
	} else {
		c.Flags &= ^FlagCarry
	}

	return 4
}

// DAA: 0x27
func instrDAA(op uint8, c *CPU) uint64 {
	nc := (c.Flags & FlagCarry) != 0

	var add uint8
	if (c.Flags&FlagAuxCarry) != 0 || (c.Registers[A]&0xF) > 0x9 {
		add += 0x6
	}

	if (c.Flags&FlagCarry) != 0 || (c.Registers[A]&0xF0) > 0x90 || (c.Registers[A]&0xF0 >= 0x90 && c.Registers[A]&0xF > 0x9) {
		add += 0x60
		nc = true
	}

	insSetreg8(c, A, flaggedAdd8(c, c.Registers[A], add, 0))

	if nc {
		c.Flags |= FlagCarry
	}

	return 4
}

// STC: 0x37
func instrSTC(op uint8, c *CPU) uint64 {
	c.Flags |= FlagCarry
	return 4
}

// RRC: 0x0F
func instrRRC(op uint8, c *CPU) uint64 {
	if (c.Registers[A] & 0x01) != 0 {
		c.Flags |= FlagCarry
	} else {
		c.Flags &= ^FlagCarry
	}

	return instrRAR(op, c)
}

// RAR: 0x1F
func instrRAR(op uint8, c *CPU) uint64 {
	a := c.Registers[A]
	cy := (a & 0x01) != 0

	c.Registers[A] = a >> 1

	if (c.Flags & FlagCarry) != 0 {
		c.Registers[A] |= 0x80
	}

	if cy {
		c.Flags |= FlagCarry
	} else {
		c.Flags &= ^FlagCarry
	}

	return 4
}

// CMA: 0x2F
func instrCMA(op uint8, c *CPU) uint64 {
	c.Registers[A] = ^c.Registers[A]
	return 4
}

// CMC: 0x3F

func instrCMC(op uint8, c *CPU) uint64 {
	c.Flags ^= FlagCarry
	return 4
}

// ADD: 0x80 to 0x87
func instrADD(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = flaggedAdd8(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	}
	return 4
}

// ADC: 0x88 to 0x8f
func instrADC(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = flaggedAdd8C(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	}
	return 4
}

// SUB: 0x90 to 0x97
func instrSUB(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = flaggedSub8(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	}
	return 4
}

// SBB: 0x98 to 0x9f
func instrSBB(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = flaggedSub8B(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	}
	return 4
}

// ANA: 0xA0 to 0xA7
func instrANA(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	a := c.Registers[A]
	b := insGetreg8(c, reg)
	c.Registers[A] = a & b
	setResultFlags(c, c.Registers[A], 0)

	if ((a | b) & 0x8) != 0 {
		c.Flags |= FlagAuxCarry
	}

	if reg == M {
		return 7
	}
	return 4
}

// XRA: 0xA8 to 0xAF
func instrXRA(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = c.Registers[A] ^ insGetreg8(c, reg)
	setResultFlags(c, c.Registers[A], 0)

	if reg == M {
		return 7
	}
	return 4
}

// ORA: 0xB0 to 0xB7
func instrORA(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = c.Registers[A] | insGetreg8(c, reg)
	setResultFlags(c, c.Registers[A], 0)

	if reg == M {
		return 7
	}
	return 4
}

// CMP: 0xB8 to 0xBF
func instrCMP(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	flaggedSub8(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	}
	return 4

}

// ADI, ACI, SUI, SBI, ANI, XRI, ORI, CPI

// ADI: 0xC6
func instrADI(op uint8, c *CPU) uint64 {
	c.Registers[A] = flaggedAdd8(c, c.Registers[A], insArg8(c), 0)
	return 7
}

// ACI: 0xCE
func instrACI(op uint8, c *CPU) uint64 {
	c.Registers[A] = flaggedAdd8C(c, c.Registers[A], insArg8(c), 0)
	return 7
}

// SUI: 0xD6
func instrSUI(op uint8, c *CPU) uint64 {
	c.Registers[A] = flaggedSub8(c, c.Registers[A], insArg8(c), 0)
	return 7
}

// SBI: 0xDE
func instrSBI(op uint8, c *CPU) uint64 {
	c.Registers[A] = flaggedSub8B(c, c.Registers[A], insArg8(c), 0)
	return 7
}

// ANI: 0xE6
func instrANI(op uint8, c *CPU) uint64 {
	a := c.Registers[A]
	b := insArg8(c)
	c.Registers[A] = a & b
	setResultFlags(c, c.Registers[A], 0)

	if ((a | b) & 0x8) != 0 {
		c.Flags |= FlagAuxCarry
	}
	return 7
}

// XRI: 0xEE
func instrXRI(op uint8, c *CPU) uint64 {
	c.Registers[A] = c.Registers[A] ^ insArg8(c)
	setResultFlags(c, c.Registers[A], 0)
	return 7
}

// ORI: 0xF6
func instrORI(op uint8, c *CPU) uint64 {
	c.Registers[A] = c.Registers[A] | insArg8(c)
	setResultFlags(c, c.Registers[A], 0)
	return 7
}

// CPI: 0xFE
func instrCPI(op uint8, c *CPU) uint64 {
	flaggedSub8(c, c.Registers[A], insArg8(c), 0)
	return 7
}
