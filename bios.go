package i8080

import (
	"fmt"
)

// Offsets
const (
	CPM  = 0xE400
	BDOS = CPM + 0x806

	BIOS = WBOOT - 3

	// CP/M 1.x
	BOOT   = 0xFA00
	WBOOT  = BOOT + 3
	CONST  = WBOOT + 3
	CONIN  = WBOOT + 6
	CONOUT = WBOOT + 9
	LIST   = WBOOT + 12
	PUNCH  = WBOOT + 15
	READER = WBOOT + 18
	HOME   = WBOOT + 21
	SELDSK = WBOOT + 24
	SETTRK = WBOOT + 27
	SETSEC = WBOOT + 30
	SETDMA = WBOOT + 33
	READ   = WBOOT + 36
	WRITE  = WBOOT + 39

	// CP/M 2.x
	LISTST  = WBOOT + 42
	SECTRAN = WBOOT + 45

	// DPH is immediately above SECTRAN jump
	// see diskio.go for constants
)

type bios struct {
	cpmImage []byte
}

func (c *CPU) initBIOS(cpmImage []byte, disks []Disk) {
	if cpmImage != nil {
		c.bios.cpmImage = cpmImage
	}

	copy(c.Memory[CPM:], c.bios.cpmImage)
	c.PC = CPM

	// Setup 0xDD hook to call BIOS
	for i := BOOT; i <= SECTRAN; i += 3 {
		c.Memory[i] = 0xDD
	}

	// not sure why this is needed, but it prevents CP/M trying to select weird disks
	c.Registers[C] = 0

	// Setup jump vectors for CP/M
	c.Memory[0] = 0xC3 // JMP
	c.Write16(1, WBOOT)

	c.Memory[3] = 0x00
	c.Memory[4] = 0x00

	c.Memory[5] = 0xC3 // JMP
	c.Write16(6, BDOS)

	c.initDisks(disks)
}

func instrBIOS(op uint8, c *CPU) uint64 {
	switch c.PC - 1 {
	case BOOT:
		panic("BOOT unimplemented")
	case WBOOT:
		c.initBIOS(nil, nil)
		c.PC = CPM
		return 100
	case CONST:
		if c.conHasChar() {
			c.Registers[A] = 0xFF
		} else {
			c.Registers[A] = 0x00
		}
	case CONIN:
		c.Registers[A] = c.conGetChar()
	case CONOUT:
		c.conPutChar(c.Registers[C])
	case LIST:
		panic("LIST unimplemented")
	case PUNCH:
		panic("PUNCH unimplemented")
	case READER:
		panic("READER unimplemented")
	case HOME:
		c.diskTrack(0)
	case SELDSK:
		c.SetR16(HL, c.diskSelect(c.Registers[C]))
	case SETTRK:
		c.diskTrack(c.GetR16(BC))
	case SETSEC:
		c.diskSector(c.GetR16(BC))
	case SETDMA:
		c.diskDMA(c.GetR16(BC))
	case READ:
		c.Registers[A] = c.diskRead()
	case WRITE:
		c.Registers[A] = c.diskWrite()
	case LISTST:
		panic("LISTST unimplemented")
	case SECTRAN:
		c.SetR16(HL, c.diskSectran(c.GetR16(DE), c.GetR16(BC)))
	default:
		panic(fmt.Sprintf(
			"unexpected BIOS call to %x, expected %x, %x, %x, ..., %x, %x, %x",
			c.PC-1,
			BOOT, WBOOT, CONST,
			SETDMA, READ, WRITE,
		))
	}

	instrRET(0, c)
	return 100
}
