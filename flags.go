package i8080

type flags uint8

const (
	F_CARRY     flags = 1 << iota // C
	F_BIT1_1                      // unused, always 1
	F_PARITY                      // P
	F_BIT3_0                      // unused, always 0
	F_AUX_CARRY                   // A or AC
	F_BIT5_0                      // unused, always 0
	F_ZERO                        // Z
	F_SIGN                        // S
)

// Clears the carry flag, then calls f(), restoring the carry flag if required
// Intended for instructions like INR and DCR, which set all flags except carry
// With this function, INR can be implemented as a flaggedAdd8 with exclude=F_CARRY
func withClearedCarry(c *CPU, exclude flags, f func() uint8) uint8 {
	oldCy := (c.Flags & F_CARRY) != 0
	c.Flags &= ^F_CARRY

	ret := f()

	if (exclude & F_CARRY) != 0 {
		if oldCy {
			c.Flags |= F_CARRY
		} else {
			c.Flags &= ^F_CARRY
		}
	}

	return ret
}

func flaggedAdd8(c *CPU, a, b uint8, exclude flags) uint8 {
	// Piggyback the "add with carry" logic, after clearing the carry flag
	return withClearedCarry(c, exclude, func() uint8 {
		return flaggedAdd8C(c, a, b, exclude)
	})
}

func flaggedAdd8C(c *CPU, a, b uint8, exclude flags) uint8 {
	var cy uint8 = 0
	if (c.Flags & F_CARRY) != 0 {
		cy = 1
	}

	ret := a + b + cy

	setResultFlags(c, ret, exclude)

	if (exclude&F_CARRY) == 0 && uint16(a)+uint16(b)+uint16(cy) >= 256 { // If a carry happened out of the 8-bit value
		c.Flags |= F_CARRY
	}

	index := (a&0x8)>>1 | (b&0x8)>>2 | (ret&0x8)>>3

	if addHCT[index&0x7] {
		c.Flags |= F_AUX_CARRY
	}

	//if (exclude&F_AUX_CARRY) == 0 && (a&0xf)+(b&0xf)+cy >= 16 { // If a carry happened out of the lower 4 bits
	//	c.Flags |= F_AUX_CARRY
	//}

	return ret
}

func flaggedSub8(c *CPU, a, b uint8, exclude flags) uint8 {
	// Piggyback the "subtract with borrow" logic, after clearing the carry flag
	return withClearedCarry(c, exclude, func() uint8 {
		return flaggedSub8B(c, a, b, exclude)
	})
}

var addHCT = []bool{false, false, true, false, true, false, true, true}
var subHCT = []bool{false, true, true, true, false, false, false, true}

func flaggedSub8B(c *CPU, a, b uint8, exclude flags) uint8 {
	var bw uint8 = 0
	if (c.Flags & F_CARRY) != 0 {
		bw = 1
	}

	ret := a - b - bw

	setResultFlags(c, ret, exclude)

	if (exclude&F_CARRY) == 0 && int16(a)-int16(b)-int16(bw) < 0 { // If a borrow happened out of the 8-bit value
		c.Flags |= F_CARRY
	}

	index := (a&0x8)>>1 | (b&0x8)>>2 | (ret&0x8)>>3

	if !subHCT[index&0x7] {
		c.Flags |= F_AUX_CARRY
	}

	//if (exclude&F_AUX_CARRY) == 0 && int16(a&0xf)-int16(b&0xf) < 0 { // If a borrow happened out of the upper 4 bits
	//	c.Flags |= F_AUX_CARRY
	//}

	return ret
}

func setResultFlags(c *CPU, result uint8, exclude flags) {
	c.Flags &= (^(F_CARRY | F_AUX_CARRY | F_SIGN | F_PARITY | F_ZERO) | exclude)

	if (exclude&F_SIGN) == 0 && (result&0x80) != 0 { // If the resulting value is a negative two's complement value
		c.Flags |= F_SIGN
	}

	if (exclude&F_PARITY) == 0 && (popcnt(result)&0x01) == 0 { // If the resulting value has odd parity
		c.Flags |= F_PARITY
	}

	if (exclude&F_ZERO) == 0 && result == 0 { // If the resulting value is zero
		c.Flags |= F_ZERO
	}
}

func popcnt(x uint8) (ret int) {
	for x != 0 {
		ret += int(x & 0x1)
		x = x >> 1
	}
	return
}
