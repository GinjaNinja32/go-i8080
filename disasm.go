package i8080

import (
	"fmt"
	"strings"
)

var disasmTable = [256]string{
	"NOP", "LXI B,d16", "STAX B", "INX B", "INR B", "DCR B", "MVI B,d8", "RLC", "*NOP", "DAD B", "LDAX B", "DCX B", "INR C", "DCR C", "MVI C,d8", "RRC", "*NOP", "LXI D,d16", "STAX D", "INX D", "INR D", "DCR D", "MVI D,d8", "RAL", "*NOP", "DAD D", "LDAX D", "DCX D", "INR E", "DCR E", "MVI E,d8", "RAR", "*NOP", "LXI H,d16", "SHLD a16", "INX H", "INR H", "DCR H", "MVI H,d8", "DAA", "*NOP", "DAD H", "LHLD a16", "DCX H", "INR L", "DCR L", "MVI L,d8", "CMA", "*NOP", "LXI SP,d16", "STA a16", "INX SP", "INR M", "DCR M", "MVI M,d8", "STC", "*NOP", "DAD SP", "LDA a16", "DCX SP", "INR A", "DCR A", "MVI A,d8", "CMC", "MOV B,B", "MOV B,C", "MOV B,D", "MOV B,E", "MOV B,H", "MOV B,L", "MOV B,M", "MOV B,A", "MOV C,B", "MOV C,C", "MOV C,D", "MOV C,E", "MOV C,H", "MOV C,L", "MOV C,M", "MOV C,A", "MOV D,B", "MOV D,C", "MOV D,D", "MOV D,E", "MOV D,H", "MOV D,L", "MOV D,M", "MOV D,A", "MOV E,B", "MOV E,C", "MOV E,D", "MOV E,E", "MOV E,H", "MOV E,L", "MOV E,M", "MOV E,A", "MOV H,B", "MOV H,C", "MOV H,D", "MOV H,E", "MOV H,H", "MOV H,L", "MOV H,M", "MOV H,A", "MOV L,B", "MOV L,C", "MOV L,D", "MOV L,E", "MOV L,H", "MOV L,L", "MOV L,M", "MOV L,A", "MOV M,B", "MOV M,C", "MOV M,D", "MOV M,E", "MOV M,H", "MOV M,L", "HLT", "MOV M,A", "MOV A,B", "MOV A,C", "MOV A,D", "MOV A,E", "MOV A,H", "MOV A,L", "MOV A,M", "MOV A,A", "ADD B", "ADD C", "ADD D", "ADD E", "ADD H", "ADD L", "ADD M", "ADD A", "ADC B", "ADC C", "ADC D", "ADC E", "ADC H", "ADC L", "ADC M", "ADC A", "SUB B", "SUB C", "SUB D", "SUB E", "SUB H", "SUB L", "SUB M", "SUB A", "SBB B", "SBB C", "SBB D", "SBB E", "SBB H", "SBB L", "SBB M", "SBB A", "ANA B", "ANA C", "ANA D", "ANA E", "ANA H", "ANA L", "ANA M", "ANA A", "XRA B", "XRA C", "XRA D", "XRA E", "XRA H", "XRA L", "XRA M", "XRA A", "ORA B", "ORA C", "ORA D", "ORA E", "ORA H", "ORA L", "ORA M", "ORA A", "CMP B", "CMP C", "CMP D", "CMP E", "CMP H", "CMP L", "CMP M", "CMP A", "RNZ", "POP B", "JNZ a16", "JMP a16", "CNZ a16", "PUSH B", "ADI d8", "RST 0", "RZ", "RET", "JZ a16", "*JMP a16", "CZ a16", "CALL a16", "ACI d8", "RST 1", "RNC", "POP D", "JNC a16", "OUT d8", "CNC a16", "PUSH D", "SUI d8", "RST 2", "RC", "*RET", "JC a16", "IN d8", "CC a16", "*CALL a16", "SBI d8", "RST 3", "RPO", "POP H", "JPO a16", "XTHL", "CPO a16", "PUSH H", "ANI d8", "RST 4", "RPE", "PCHL", "JPE a16", "XCHG", "CPE a16", "*CALL a16", "XRI d8", "RST 5", "RP", "POP PSW", "JP a16", "DI", "CP a16", "PUSH PSW", "ORI d8", "RST 6", "RM", "SPHL", "JM a16", "EI", "CM a16", "*CALL a16", "CPI d8", "RST 7",
}

func (c *CPU) disasmPC() (string, string) {
	op := c.Memory[c.PC]

	if op == 0xDD {
		// BIOS calls
		return "BIOS", fmt.Sprintf("%d", int(c.PC)-WBOOT)
	}

	str := disasmTable[op]

	str = strings.Replace(str, "d8", fmt.Sprintf("#%02x", c.Memory[c.PC+1]), -1)
	str = strings.Replace(str, "d16", fmt.Sprintf("#%04x", c.Read16(c.PC+1)), -1)
	str = strings.Replace(str, "a16", fmt.Sprintf("#%04x", c.Read16(c.PC+1)), -1)

	strs := strings.Split(str, " ")

	if len(strs) == 1 {
		return strs[0], ""
	}
	return strs[0], strs[1]
}
