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
	c.PC += 1
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

var ops = [256]func(uint8, *CPU) int{
	i_nop, i_lxi, i_stax, i_inx, i_inr, i_dcr, i_mvi, i_rlc, i_nop, i_dad, i_ldax, i_dcx, i_inr, i_dcr, i_mvi, i_rrc,
	i_nop, i_lxi, i_stax, i_inx, i_inr, i_dcr, i_mvi, i_ral, i_nop, i_dad, i_ldax, i_dcx, i_inr, i_dcr, i_mvi, i_rar,
	i_nop, i_lxi, i_shld, i_inx, i_inr, i_dcr, i_mvi, i_daa, i_nop, i_dad, i_lhld, i_dcx, i_inr, i_dcr, i_mvi, i_cma,
	i_nop, i_lxi, i_sta, i_inx, i_inr, i_dcr, i_mvi, i_stc, i_nop, i_dad, i_lda, i_dcx, i_inr, i_dcr, i_mvi, i_cmc,
	i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov,
	i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov,
	i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov,
	i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_hlt, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov, i_mov,
	i_add, i_add, i_add, i_add, i_add, i_add, i_add, i_add, i_adc, i_adc, i_adc, i_adc, i_adc, i_adc, i_adc, i_adc,
	i_sub, i_sub, i_sub, i_sub, i_sub, i_sub, i_sub, i_sub, i_sbb, i_sbb, i_sbb, i_sbb, i_sbb, i_sbb, i_sbb, i_sbb,
	i_ana, i_ana, i_ana, i_ana, i_ana, i_ana, i_ana, i_ana, i_xra, i_xra, i_xra, i_xra, i_xra, i_xra, i_xra, i_xra,
	i_ora, i_ora, i_ora, i_ora, i_ora, i_ora, i_ora, i_ora, i_cmp, i_cmp, i_cmp, i_cmp, i_cmp, i_cmp, i_cmp, i_cmp,
	i_cond_ret, i_pop, i_cond_jmp, i_jmp, i_cond_call, i_push, i_adi, i_rst, i_cond_ret, i_ret, i_cond_jmp, i_jmp, i_cond_call, i_call, i_aci, i_rst,
	i_cond_ret, i_pop, i_cond_jmp, i_out, i_cond_call, i_push, i_sui, i_rst, i_cond_ret, i_ret, i_cond_jmp, i_in, i_cond_call, i_call, i_sbi, i_rst,
	i_cond_ret, i_pop, i_cond_jmp, i_xthl, i_cond_call, i_push, i_ani, i_rst, i_cond_ret, i_pchl, i_cond_jmp, i_xchg, i_cond_call, i_call, i_xri, i_rst,
	i_cond_ret, i_pop, i_cond_jmp, i_di, i_cond_call, i_push, i_ori, i_rst, i_cond_ret, i_sphl, i_cond_jmp, i_ei, i_cond_call, i_call, i_cpi, i_rst,
}
