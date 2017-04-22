package i8080

// 16-bit registers helpers

const (
	BC = B
	DE = D
	HL = H
	R4 = M
)

func (c *CPU) GetR16(reg uint8) uint16 {
	switch reg {
	case BC:
		return c.BC()
	case DE:
		return c.DE()
	case HL:
		return c.HL()
	default:
		panic("invalid 16-bit register")
	}
}

func (c *CPU) SetR16(reg uint8, val uint16) {
	switch reg {
	case BC:
		c.SetBC(val)
	case DE:
		c.SetDE(val)
	case HL:
		c.SetHL(val)
	default:
		panic("invalid 16-bit register")
	}
}

func (c *CPU) BC() uint16 {
	return uint16(c.Registers[B])<<8 | uint16(c.Registers[C])
}

func (c *CPU) SetBC(v uint16) {
	c.Registers[B] = uint8(v >> 8)
	c.Registers[C] = uint8(v)
}

func (c *CPU) DE() uint16 {
	return uint16(c.Registers[D])<<8 | uint16(c.Registers[E])
}

func (c *CPU) SetDE(v uint16) {
	c.Registers[D] = uint8(v >> 8)
	c.Registers[E] = uint8(v)
}

func (c *CPU) HL() uint16 {
	return uint16(c.Registers[H])<<8 | uint16(c.Registers[L])
}

func (c *CPU) SetHL(v uint16) {
	c.Registers[H] = uint8(v >> 8)
	c.Registers[L] = uint8(v)
}

func (c *CPU) PSW() uint16 {
	return uint16(c.Registers[A])<<8 | uint16(c.Flags)
}

func (c *CPU) SetPSW(v uint16) {
	c.Registers[A] = uint8(v >> 8)
	c.Flags = flags(v)

	c.Flags |= F_BIT1_1
	c.Flags &= ^(F_BIT3_0 | F_BIT5_0)
}

// PUSH/POP and 16-bit memory helpers

func (c *CPU) Push(v uint16) {
	c.Write16(c.SP-2, v)
	c.SP -= 2
}

func (c *CPU) Pop() uint16 {
	c.SP += 2
	return c.Read16(c.SP - 2)
}

func (c *CPU) Write16(addr uint16, v uint16) {
	c.Memory[addr+1] = uint8(v >> 8)
	c.Memory[addr] = uint8(v)
}

func (c *CPU) Read16(addr uint16) uint16 {
	return uint16(c.Memory[addr+1])<<8 | uint16(c.Memory[addr])
}
