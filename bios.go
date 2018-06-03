package i8080

import (
	"fmt"
)

// Constants
const (
	MSIZE = 64

	SECTORS_PER_DISK = 26

	DPH_SIZE   = 16
	DPB_SIZE   = 16 // ?
	ALLOC_SIZE = 32
	CHECK_SIZE = 16

	DirBf_SIZE = 128

	DMA_SIZE = 128

	//CBASE    = (MSIZE - 17) * 1024
	//CPMB     = CBASE + 0x2900
	//BDOS     = CBASE + 0x3106
	//BIOSBase = CPMB + 0x1500

	//CCP = 0xDC00
)

func dph(c *CPU, disk int) uint16 {
	return 0
}

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

	DPH  = SECTRAN + 3
	DPH0 = DPH
	DPH1 = DPH + 1*DPH_SIZE
	DPH2 = DPH + 2*DPH_SIZE
	DPH3 = DPH + 3*DPH_SIZE

	SecTrans = DPH + 4*DPH_SIZE

	DPB = SecTrans + SECTORS_PER_DISK

	DirBf = DPB + DPB_SIZE

	Alloc  = DirBf + DirBf_SIZE
	Alloc0 = Alloc
	Alloc1 = Alloc + 1*ALLOC_SIZE
	Alloc2 = Alloc + 2*ALLOC_SIZE
	Alloc3 = Alloc + 3*ALLOC_SIZE

	Check  = Alloc + 4*ALLOC_SIZE
	Check0 = Check
	Check1 = Check + 1*CHECK_SIZE
	Check2 = Check + 2*CHECK_SIZE
	Check3 = Check + 3*CHECK_SIZE

	END = Check + 4*CHECK_SIZE
)

type bios struct {
	dmaAddress    uint16
	currentDisk   uint8
	currentTrack  uint16
	currentSector uint16

	cpmImage []byte
	disks    []Disk
}

type Disk struct {
	Data     []byte
	ReadOnly bool
}

func (c *CPU) initBIOS(cpmImage []byte, disks []Disk) {
	if cpmImage != nil {
		c.bios.cpmImage = cpmImage
	}
	if disks != nil {
		c.bios.disks = disks
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

	// Setup DPH
	dph := []byte{
		uint8(SecTrans & 0xFF), uint8(SecTrans >> 8),
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ?
		uint8(DirBf & 0xFF), uint8(DirBf >> 8),
		uint8(DPB & 0xFF), uint8(DPB >> 8),
		// 0x00, 0x00, // CHK
		// 0x00, 0x00, // ALLOC
	}

	copy(c.Memory[DPH0:], dph)
	copy(c.Memory[DPH1:], dph)
	copy(c.Memory[DPH2:], dph)
	copy(c.Memory[DPH3:], dph)

	c.Write16(uint16(DPH0+len(dph)), Check0)
	c.Write16(uint16(DPH0+len(dph)+2), Alloc0)
	c.Write16(uint16(DPH1+len(dph)), Check1)
	c.Write16(uint16(DPH1+len(dph)+2), Alloc1)
	c.Write16(uint16(DPH2+len(dph)), Check2)
	c.Write16(uint16(DPH2+len(dph)+2), Alloc2)
	c.Write16(uint16(DPH3+len(dph)), Check3)
	c.Write16(uint16(DPH3+len(dph)+2), Alloc3)

	copy(c.Memory[SecTrans:], []byte{
		1, 7, 13, 19, 25, 5, 11, 17, 23, 3, 9, 15, 21, 2, 8, 14, 20, 26, 6, 12, 18, 24, 4, 10, 16, 22,
	})

	copy(c.Memory[DPB:], []byte{
		26, 0, // sectors per track
		3,      // block shift factor ??
		7,      // block mask ??
		0,      // null mask ??
		242, 0, // disk size - 1 ??
		63, 0, // directory max ??
		192,   // alloc 0
		0,     // alloc 1
		15, 0, // check size
		2, 0, // track offset ??
	})
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
		if c.ioHasChar() {
			c.Registers[A] = 0xFF
		} else {
			c.Registers[A] = 0x00
		}
	case CONIN:
		c.Registers[A] = c.ioGetChar()
	case CONOUT:
		c.ioPutChar(c.Registers[C])
	case LIST:
		panic("LIST unimplemented")
	case PUNCH:
		panic("PUNCH unimplemented")
	case READER:
		panic("READER unimplemented")
	case HOME:
		c.bios.currentTrack = 0
	case SELDSK:
		requestedDisk := c.Registers[C]
		if int(requestedDisk) < len(c.bios.disks) {
			c.bios.currentDisk = requestedDisk
			c.SetR16(HL, DPH+uint16(requestedDisk)*DPH_SIZE)
		} else {
			c.SetR16(HL, 0)
		}
	case SETTRK:
		c.bios.currentTrack = c.GetR16(BC)
	case SETSEC:
		c.bios.currentSector = c.GetR16(BC)
	case SETDMA:
		c.bios.dmaAddress = c.GetR16(BC)
	case READ:
		offset := 128 * int(26*c.bios.currentTrack+c.bios.currentSector)
		for i := 0; i < DMA_SIZE; i++ {
			c.Memory[int(c.bios.dmaAddress)+i] = c.bios.disks[c.bios.currentDisk].Data[offset+i]
		}
		c.Registers[A] = 0
	case WRITE:
		if c.bios.disks[c.bios.currentDisk].ReadOnly {
			c.Registers[A] = 2
		} else {
			offset := 128 * int(26*c.bios.currentTrack+c.bios.currentSector)
			for i := 0; i < DMA_SIZE; i++ {
				c.bios.disks[c.bios.currentDisk].Data[offset+i] = c.Memory[int(c.bios.dmaAddress)+i]
			}
			c.Registers[A] = 0
		}
	case LISTST:
		panic("LISTST unimplemented")
	case SECTRAN:
		addr := c.GetR16(DE)
		c.SetR16(HL, uint16(c.Memory[addr+c.GetR16(BC)])-1)
	default:
		panic(fmt.Sprintf(
			"unimplemented BIOS call to %x, expected %x, %x, %x, ..., %x, %x, %x",
			c.PC-1,
			BOOT, WBOOT, CONST,
			SETDMA, READ, WRITE,
		))
	}

	instrRET(0, c)
	return 100
}
