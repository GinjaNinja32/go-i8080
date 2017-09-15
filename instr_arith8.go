package i8080

/*
	This file:

	INR, DCR
	RLC, RAL, DAA, STC, RRC, RAR, CMA, CMC
	ADD, ADC, SUB, SBB, ANA, XRA, ORA, CMP
	ADI, ACI, SUI, SBI, ANI, XRI, ORI, CPI
*/

// INR: 0x04, 0x0C, 0x14, 0x1C, 0x24, 0x2C, 0x34, 0x3C
func i_inr(op uint8, c *CPU) uint64 {
	reg := insArg3b(op)
	insSetreg8(c, reg, flaggedAdd8(c, insGetreg8(c, reg), 1, F_CARRY))
	if reg == M {
		return 10
	} else {
		return 5
	}
}

// DCR: 0x05, 0x0D, 0x15, 0x1D, 0x25, 0x2D, 0x35, 0x3D
func i_dcr(op uint8, c *CPU) uint64 {
	reg := insArg3b(op)

	result := flaggedSub8(c, insGetreg8(c, reg), 1, F_CARRY)
	insSetreg8(c, reg, result)

	if (result & 0xf) != 0xf {
		c.Flags |= F_AUX_CARRY
	} else {
		c.Flags &= ^F_AUX_CARRY
	}

	if reg == M {
		return 10
	} else {
		return 5
	}
}

// RLC: 0x07
func i_rlc(op uint8, c *CPU) uint64 {
	if (c.Registers[A] & 0x80) != 0 {
		c.Flags |= F_CARRY
	} else {
		c.Flags &= ^F_CARRY
	}

	return i_ral(op, c)
}

// RAL: 0x17
func i_ral(op uint8, c *CPU) uint64 {
	a := c.Registers[A]
	cy := (a & 0x80) != 0

	c.Registers[A] = a << 1

	if (c.Flags & F_CARRY) != 0 {
		c.Registers[A] |= 0x01
	}

	if cy {
		c.Flags |= F_CARRY
	} else {
		c.Flags &= ^F_CARRY
	}

	return 4
}

// DAA: 0x27
func i_daa(op uint8, c *CPU) uint64 {
	nc := (c.Flags & F_CARRY) != 0

	var add uint8 = 0
	if (c.Flags&F_AUX_CARRY) != 0 || (c.Registers[A]&0xF) > 0x9 {
		add += 0x6
	}

	if (c.Flags&F_CARRY) != 0 || (c.Registers[A]&0xF0) > 0x90 || (c.Registers[A]&0xF0 >= 0x90 && c.Registers[A]&0xF > 0x9) {
		add += 0x60
		nc = true
	}

	insSetreg8(c, A, flaggedAdd8(c, c.Registers[A], add, 0))

	if nc {
		c.Flags |= F_CARRY
	}

	return 4
}

// STC: 0x37
func i_stc(op uint8, c *CPU) uint64 {
	c.Flags |= F_CARRY
	return 4
}

// RRC: 0x0F
func i_rrc(op uint8, c *CPU) uint64 {
	if (c.Registers[A] & 0x01) != 0 {
		c.Flags |= F_CARRY
	} else {
		c.Flags &= ^F_CARRY
	}

	return i_rar(op, c)
}

// RAR: 0x1F
func i_rar(op uint8, c *CPU) uint64 {
	a := c.Registers[A]
	cy := (a & 0x01) != 0

	c.Registers[A] = a >> 1

	if (c.Flags & F_CARRY) != 0 {
		c.Registers[A] |= 0x80
	}

	if cy {
		c.Flags |= F_CARRY
	} else {
		c.Flags &= ^F_CARRY
	}

	return 4
}

// CMA: 0x2F
func i_cma(op uint8, c *CPU) uint64 {
	c.Registers[A] = ^c.Registers[A]
	return 4
}

// CMC: 0x3F

func i_cmc(op uint8, c *CPU) uint64 {
	c.Flags ^= F_CARRY
	return 4
}

// ADD: 0x80 to 0x87
func i_add(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = flaggedAdd8(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	} else {
		return 4
	}
}

// ADC: 0x88 to 0x8f
func i_adc(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = flaggedAdd8C(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	} else {
		return 4
	}
}

// SUB: 0x90 to 0x97
func i_sub(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = flaggedSub8(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	} else {
		return 4
	}
}

// SBB: 0x98 to 0x9f
func i_sbb(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = flaggedSub8B(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	} else {
		return 4
	}
}

// ANA: 0xA0 to 0xA7
func i_ana(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	a := c.Registers[A]
	b := insGetreg8(c, reg)
	c.Registers[A] = a & b
	setResultFlags(c, c.Registers[A], 0)

	if ((a | b) & 0x8) != 0 {
		c.Flags |= F_AUX_CARRY
	}

	if reg == M {
		return 7
	} else {
		return 4
	}
}

// XRA: 0xA8 to 0xAF
func i_xra(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = c.Registers[A] ^ insGetreg8(c, reg)
	setResultFlags(c, c.Registers[A], 0)

	if reg == M {
		return 7
	} else {
		return 4
	}
}

// ORA: 0xB0 to 0xB7
func i_ora(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	c.Registers[A] = c.Registers[A] | insGetreg8(c, reg)
	setResultFlags(c, c.Registers[A], 0)

	if reg == M {
		return 7
	} else {
		return 4
	}
}

// CMP: 0xB8 to 0xBF
func i_cmp(op uint8, c *CPU) uint64 {
	reg := insArg3(op)
	flaggedSub8(c, c.Registers[A], insGetreg8(c, reg), 0)

	if reg == M {
		return 7
	} else {
		return 4
	}
}

// ADI, ACI, SUI, SBI, ANI, XRI, ORI, CPI

// ADI: 0xC6
func i_adi(op uint8, c *CPU) uint64 {
	c.Registers[A] = flaggedAdd8(c, c.Registers[A], insArg8(c), 0)
	return 7
}

// ACI: 0xCE
func i_aci(op uint8, c *CPU) uint64 {
	c.Registers[A] = flaggedAdd8C(c, c.Registers[A], insArg8(c), 0)
	return 7
}

// SUI: 0xD6
func i_sui(op uint8, c *CPU) uint64 {
	c.Registers[A] = flaggedSub8(c, c.Registers[A], insArg8(c), 0)
	return 7
}

// SBI: 0xDE
func i_sbi(op uint8, c *CPU) uint64 {
	c.Registers[A] = flaggedSub8B(c, c.Registers[A], insArg8(c), 0)
	return 7
}

// ANI: 0xE6
func i_ani(op uint8, c *CPU) uint64 {
	a := c.Registers[A]
	b := insArg8(c)
	c.Registers[A] = a & b
	setResultFlags(c, c.Registers[A], 0)

	if ((a | b) & 0x8) != 0 {
		c.Flags |= F_AUX_CARRY
	}
	return 7
}

// XRI: 0xEE
func i_xri(op uint8, c *CPU) uint64 {
	c.Registers[A] = c.Registers[A] ^ insArg8(c)
	setResultFlags(c, c.Registers[A], 0)
	return 7
}

// ORI: 0xF6
func i_ori(op uint8, c *CPU) uint64 {
	c.Registers[A] = c.Registers[A] | insArg8(c)
	setResultFlags(c, c.Registers[A], 0)
	return 7
}

// CPI: 0xFE
func i_cpi(op uint8, c *CPU) uint64 {
	flaggedSub8(c, c.Registers[A], insArg8(c), 0)
	return 7
}
