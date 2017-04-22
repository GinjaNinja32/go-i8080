package i8080

func i_nop(op uint8, c *CPU) int {
	return 4
}

func i_hlt(op uint8, c *CPU) int {
	return -1
}

func i_out(op uint8, c *CPU) int {
	// TODO
	return 10
}

func i_in(op uint8, c *CPU) int {
	// TODO
	return 10
}

func i_di(op uint8, c *CPU) int {
	// TODO
	return 4
}

func i_ei(op uint8, c *CPU) int {
	// TODO
	return 4
}
