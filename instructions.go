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

func (c *CPU) dispatch(op uint8) func(uint8, *CPU) int {
	switch op & 0xf0 {
	case 0x40, 0x50, 0x60, 0x70: // MOV and HLT
		if op == 0x76 { // 0x76 would be MOV M, M
			return i_hlt
		} else {
			return i_mov
		}
	case 0x80:
		if (op & 0x8) != 0 {
			return i_adc
		} else {
			return i_add
		}
	case 0x90:
		if (op & 0x8) != 0 {
			return i_sbb
		} else {
			return i_sub
		}
	case 0xa0:
		switch op & 0x8 {
		case 0x0:
			return i_ana
		case 0x8:
			return i_xra
		}
	case 0xb0:
		switch op & 0x8 {
		case 0x0:
			return i_ora
		case 0x8:
			return i_cmp
		}
	}

	if op <= 0x3f {
		switch op & 0xf {
		case 0x0, 0x8:
			return i_nop
		case 0x1:
			return i_lxi
		case 0x3:
			return i_inx
		case 0x4, 0xc:
			return i_inr
		case 0x5, 0xd:
			return i_dcr
		case 0x6, 0xe:
			return i_mvi
		case 0x9:
			return i_dad
		case 0xb:
			return i_dcx
		case 0x2, 0x7, 0xa, 0xf:
			switch op {
			case 0x02, 0x12:
				return i_stax
			case 0x22:
				return i_shld
			case 0x32:
				return i_sta

			case 0x07:
				return i_rlc
			case 0x17:
				return i_ral
			case 0x27:
				return i_daa
			case 0x37:
				return i_stc

			case 0x0a, 0x1a:
				return i_ldax
			case 0x2a:
				return i_lhld
			case 0x3a:
				return i_lda

			case 0x0f:
				return i_rrc
			case 0x1f:
				return i_rar
			case 0x2f:
				return i_cma
			case 0x3f:
				return i_cmc
			}
		}
	} else {
		switch op & 0xf {
		case 0x0, 0x8:

			return i_cond_ret
		case 0x1:
			return i_pop
		case 0x2, 0xa:
			return i_cond_jmp
		case 0x4, 0xc:
			return i_cond_call
		case 0x5:
			return i_push
		case 0x7, 0xf:
			return i_rst
		case 0xd:
			return i_call

		case 0x3, 0x6, 0x9, 0xb, 0xe:
			switch op {
			case 0xc3:
				return i_jmp
			case 0xd3:
				return i_out
			case 0xe3:
				return i_xthl
			case 0xf3:
				return i_di

			case 0xc6:
				return i_adi
			case 0xd6:
				return i_sui
			case 0xe6:
				return i_ani
			case 0xf6:
				return i_ori

			case 0xc9, 0xd9:
				return i_ret
			case 0xe9:
				return i_pchl
			case 0xf9:
				return i_sphl

			case 0xcb:
				return i_jmp
			case 0xdb:
				return i_in
			case 0xeb:
				return i_xchg
			case 0xfb:
				return i_ei

			case 0xce:
				return i_aci
			case 0xde:
				return i_sbi
			case 0xee:
				return i_xri
			case 0xfe:
				return i_cpi
			}
		}
	}

	return func(op uint8, c *CPU) int {
		return -3
	}
}
