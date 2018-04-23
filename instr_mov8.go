package i8080

/*
	This file:

	STAX, STA
	LDAX, LDA
	MVI, MOV
*/

// STAX: 0x02, 0x12
func instrSTAX(op uint8, c *CPU) uint64 {
	dst := insArg2(op) // actually 1 bit, but the 2nd bit of insArg2 is always 0 in the two STAX instructions

	if dst == 0 { // STAX B
		c.Memory[c.BC()] = c.Registers[A]
	} else { // STAX D
		c.Memory[c.DE()] = c.Registers[A]
	}

	return 7
}

// STA: 0x32
func instrSTA(op uint8, c *CPU) uint64 {
	c.Memory[insArg16(c)] = c.Registers[A]
	return 13
}

// LDAX: 0x0A, 0x1A
func instrLDAX(op uint8, c *CPU) uint64 {
	dst := insArg2(op) // actually 1 bit, but the 2nd bit of insArg2 is always 0 in the two LDAX instructions

	if dst == 0 { // LDAX B
		c.Registers[A] = c.Memory[c.BC()]
	} else { // LDAX D
		c.Registers[A] = c.Memory[c.DE()]
	}

	return 7
}

// LDA: 0x3A
func instrLDA(op uint8, c *CPU) uint64 {
	c.Registers[A] = c.Memory[insArg16(c)]
	return 13
}

// MVI: 0x06, 0x0E, 0x16, 0x1E, 0x26, 0x2E, 0x36, 0x3E
func instrMVI(op uint8, c *CPU) uint64 {
	dst := insArg3b(op)

	insSetreg8(c, dst, insArg8(c))

	if dst == M {
		return 10
	}
	return 7
}

// MOV: 0x40 to 0x75, 0x77 to 0x7f
func instrMOV(op uint8, c *CPU) uint64 {
	src := insArg3(op)
	dst := insArg3b(op)

	insSetreg8(c, dst, insGetreg8(c, src))

	if src == M || dst == M {
		return 7
	}
	return 5
}
