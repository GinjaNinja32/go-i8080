package i8080

// 16-bit registers helpers

// 16-bit register aliases
const (
	BC = B
	DE = D
	HL = H
	R4 = M
)

// GetR16 gets a 16-bit register pair by index
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

// SetR16 sets a 16-bit register pair by index
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

// BC returns the value of the BC register pair
func (c *CPU) BC() uint16 {
	return uint16(c.Registers[B])<<8 | uint16(c.Registers[C])
}

// SetBC sets the value of the BC register pair
func (c *CPU) SetBC(v uint16) {
	c.Registers[B] = uint8(v >> 8)
	c.Registers[C] = uint8(v)
}

// DE returns the value of the DE register pair
func (c *CPU) DE() uint16 {
	return uint16(c.Registers[D])<<8 | uint16(c.Registers[E])
}

// SetDE sets the value of the DE register pair
func (c *CPU) SetDE(v uint16) {
	c.Registers[D] = uint8(v >> 8)
	c.Registers[E] = uint8(v)
}

// HL returns the value of the HL register pair
func (c *CPU) HL() uint16 {
	return uint16(c.Registers[H])<<8 | uint16(c.Registers[L])
}

// SetHL sets the value of the HL register pair
func (c *CPU) SetHL(v uint16) {
	c.Registers[H] = uint8(v >> 8)
	c.Registers[L] = uint8(v)
}

// PSW gets the value of the PSW register pair
func (c *CPU) PSW() uint16 {
	return uint16(c.Registers[A])<<8 | uint16(c.Flags)
}

// SetPSW sets the value of the PSW register pair
// The bits of the flag register that are always 1 or always 0 will be forced to those values
func (c *CPU) SetPSW(v uint16) {
	c.Registers[A] = uint8(v >> 8)
	c.Flags = flags(v)

	c.Flags |= FlagBit1
	c.Flags &= ^(FlagBit3 | FlagBit5)
}

// PUSH/POP and 16-bit memory helpers

// Push pushes a value to the stack
func (c *CPU) Push(v uint16) {
	c.SP -= 2
	c.Write16(c.SP, v)
}

// Pop pops a value from the stack
func (c *CPU) Pop() uint16 {
	v := c.Read16(c.SP)
	c.SP += 2
	return v
}

// Write16 writes a 16-bit value `v` at address `addr`
func (c *CPU) Write16(addr uint16, v uint16) {
	c.Memory[addr+1] = uint8(v >> 8)
	c.Memory[addr] = uint8(v)
}

// Read16 reads a 16-bit value from address `addr`
func (c *CPU) Read16(addr uint16) uint16 {
	return uint16(c.Memory[addr+1])<<8 | uint16(c.Memory[addr])
}
