package i8080

func insArg2(op uint8) uint8 {
	return (op >> 3) & 0x6 // 0, 2, 4, 6, ie every 2nd 8-bit register
}

func insArg3(op uint8) uint8 {
	return op & 0x7
}

func insArg3b(op uint8) uint8 {
	return (op >> 3) & 0x7
}

func insArg8(c *CPU) (ret uint8) {
	ret = c.Memory[c.PC]
	c.PC++
	return
}

func insArg16(c *CPU) (ret uint16) {
	ret = c.Read16(c.PC)
	c.PC += 2
	return
}

func insGetreg8(c *CPU, reg uint8) uint8 {
	switch reg {
	case M:
		return c.Memory[c.HL()]
	default:
		return c.Registers[reg]
	}
}

func insSetreg8(c *CPU, reg uint8, val uint8) {
	switch reg {
	case M:
		c.Memory[c.HL()] = val
	default:
		c.Registers[reg] = val
	}
}

var ops = [256]func(uint8, *CPU) uint64{
	instrNOP, instrLXI, instrSTAX, instrINX, instrINR, instrDCR, instrMVI, instrRLC, instrNOP, instrDAD, instrLDAX, instrDCX, instrINR, instrDCR, instrMVI, instrRRC,
	instrNOP, instrLXI, instrSTAX, instrINX, instrINR, instrDCR, instrMVI, instrRAL, instrNOP, instrDAD, instrLDAX, instrDCX, instrINR, instrDCR, instrMVI, instrRAR,
	instrNOP, instrLXI, instrSHLD, instrINX, instrINR, instrDCR, instrMVI, instrDAA, instrNOP, instrDAD, instrLHLD, instrDCX, instrINR, instrDCR, instrMVI, instrCMA,
	instrNOP, instrLXI, instrSTA, instrINX, instrINR, instrDCR, instrMVI, instrSTC, instrNOP, instrDAD, instrLDA, instrDCX, instrINR, instrDCR, instrMVI, instrCMC,
	instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV,
	instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV,
	instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV,
	instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrHLT, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV, instrMOV,
	instrADD, instrADD, instrADD, instrADD, instrADD, instrADD, instrADD, instrADD, instrADC, instrADC, instrADC, instrADC, instrADC, instrADC, instrADC, instrADC,
	instrSUB, instrSUB, instrSUB, instrSUB, instrSUB, instrSUB, instrSUB, instrSUB, instrSBB, instrSBB, instrSBB, instrSBB, instrSBB, instrSBB, instrSBB, instrSBB,
	instrANA, instrANA, instrANA, instrANA, instrANA, instrANA, instrANA, instrANA, instrXRA, instrXRA, instrXRA, instrXRA, instrXRA, instrXRA, instrXRA, instrXRA,
	instrORA, instrORA, instrORA, instrORA, instrORA, instrORA, instrORA, instrORA, instrCMP, instrCMP, instrCMP, instrCMP, instrCMP, instrCMP, instrCMP, instrCMP,
	instrCondRET, instrPOP, instrCondJMP, instrJMP, instrCondCALL, instrPUSH, instrADI, instrRST, instrCondRET, instrRET, instrCondJMP, instrJMP, instrCondCALL, instrCALL, instrACI, instrRST,
	instrCondRET, instrPOP, instrCondJMP, instrOUT, instrCondCALL, instrPUSH, instrSUI, instrRST, instrCondRET, instrRET, instrCondJMP, instrIN, instrCondCALL, instrBDOS, instrSBI, instrRST,
	instrCondRET, instrPOP, instrCondJMP, instrXTHL, instrCondCALL, instrPUSH, instrANI, instrRST, instrCondRET, instrPCHL, instrCondJMP, instrXCHG, instrCondCALL, instrCALL, instrXRI, instrRST,
	instrCondRET, instrPOP, instrCondJMP, instrDI, instrCondCALL, instrPUSH, instrORI, instrRST, instrCondRET, instrSPHL, instrCondJMP, instrEI, instrCondCALL, instrCALL, instrCPI, instrRST,
}
