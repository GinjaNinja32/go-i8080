package i8080

func instrNOP(op uint8, c *CPU) uint64 {
	return 4
}

func instrHLT(op uint8, c *CPU) uint64 {
	panic("hlt")
}

func instrOUT(op uint8, c *CPU) uint64 {
	// TODO
	return 10
}

func instrIN(op uint8, c *CPU) uint64 {
	// TODO
	return 10
}

func instrDI(op uint8, c *CPU) uint64 {
	// TODO
	return 4
}

func instrEI(op uint8, c *CPU) uint64 {
	// TODO
	return 4
}
