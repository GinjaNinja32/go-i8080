package i8080

import (
	"io"
)

// Disk constants
const (
	SECTORS_PER_TRACK = 26

	DPH_SIZE   = 16
	DPB_SIZE   = 16 // ?
	ALLOC_SIZE = 32
	CHECK_SIZE = 16

	DirBf_SIZE = 128

	DMA_SIZE = 128
)

// Memory layout constants
const (
	DPH  = SECTRAN + 3
	DPH0 = DPH
	DPH1 = DPH + 1*DPH_SIZE
	DPH2 = DPH + 2*DPH_SIZE
	DPH3 = DPH + 3*DPH_SIZE

	SecTrans = DPH + 4*DPH_SIZE

	DPB = SecTrans + SECTORS_PER_TRACK

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

type diskio struct {
	dmaAddress uint16
	disk       uint8
	track      uint16
	sector     uint16

	disks []Disk
}

type Disk struct {
	Data     io.ReadWriteSeeker
	ReadOnly bool
}

func (c *CPU) initDisks(disks []Disk) {
	if disks != nil {
		c.diskio.disks = disks
	}

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

func (c *CPU) diskDMA(dma uint16) {
	c.diskio.dmaAddress = dma
}

func (c *CPU) diskSelect(disk uint8) uint16 {
	if int(disk) < len(c.diskio.disks) {
		c.diskio.disk = disk
		return DPH + uint16(disk)*DPH_SIZE
	}
	return 0
}
func (c *CPU) diskTrack(track uint16) {
	c.diskio.track = track
}

func (c *CPU) diskSector(sector uint16) {
	c.diskio.sector = sector
}

func (c *CPU) diskRead() uint8 {
	offset := DMA_SIZE * int64(SECTORS_PER_TRACK*c.diskio.track+c.diskio.sector)

	disk := c.diskio.disks[c.diskio.disk].Data
	dma := c.Memory[c.diskio.dmaAddress:][:DMA_SIZE]

	_, err := disk.Seek(offset, io.SeekStart)
	if err != nil {
		panic(err)
	}

	n, err := disk.Read(dma)
	if err != nil {
		panic(err)
	}
	if n != DMA_SIZE {
		panic(n)
	}

	return 0
}

func (c *CPU) diskWrite() uint8 {
	if c.diskio.disks[c.diskio.disk].ReadOnly {
		return 2
	}

	disk := c.diskio.disks[c.diskio.disk].Data
	dma := c.Memory[c.diskio.dmaAddress:][:DMA_SIZE]

	offset := DMA_SIZE * int64(SECTORS_PER_TRACK*c.diskio.track+c.diskio.sector)

	_, err := disk.Seek(offset, io.SeekStart)
	if err != nil {
		panic(err)
	}

	n, err := disk.Write(dma)
	if err != nil {
		panic(err)
	}
	if n != DMA_SIZE {
		panic(n)
	}

	return 0
}

func (c *CPU) diskSectran(table uint16, entry uint16) uint16 {
	return uint16(c.Memory[table+entry]) - 1
}
