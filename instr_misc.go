package i8080

func i_nop(op uint8, c *CPU) uint64 {
	return 4
}

func i_hlt(op uint8, c *CPU) uint64 {
	panic("hlt")
}

func i_out(op uint8, c *CPU) uint64 {
	// TODO
	return 10
}

func i_in(op uint8, c *CPU) uint64 {
	// TODO
	return 10
}

func i_di(op uint8, c *CPU) uint64 {
	// TODO
	return 4
}

func i_ei(op uint8, c *CPU) uint64 {
	// TODO
	return 4
}
