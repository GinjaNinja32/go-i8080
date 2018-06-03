package i8080

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Inverse speed constants; Speed2MHz per instruction is a speed of 2MHz, etc
const (
	SpeedDebug    = 10 * time.Millisecond
	Speed2Mhz     = 500 * time.Nanosecond
	Speed3_125Mhz = 320 * time.Nanosecond
)

// Register index constants
const (
	B uint8 = iota
	C
	D
	E
	H
	L
	M
	A
)

// CPU implements an emulated Intel 8080 CPU
type CPU struct {
	Memory    [65536]uint8
	Registers [8]uint8
	Flags     flags

	SP uint16
	PC uint16

	ClockTime time.Duration // nanoseconds per clock tick

	bios
	conio
}

// New creates a new emulated Intel 8080 CPU
func New(conin io.Reader, conout io.Writer, cpmImage []byte, disks []Disk) (c *CPU) {
	c = &CPU{
		Flags:     FlagBit1,
		ClockTime: Speed2Mhz,
	}

	c.initBIOS(cpmImage, disks)
	c.initIO(conin, conout)

	return
}

// Output is called when a string is output
func (c *CPU) Output(s string) {
	//c.OutputStr += s
	fmt.Printf("%s", s)
}

const tickBudget = 10 * time.Millisecond

// Run runs the CPU and returns how many CPU cycles were executed before a halt
func (c *CPU) Run() (cycles uint64) {
	debug := false
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)
	go func() {
		for range sigChan {
			debug = !debug
		}
	}()

	ticker := time.NewTicker(tickBudget)

	defer func() {
		f := recover()
		if f != nil && f != "hlt" {
			panic(f)
		}
	}()
	defer ticker.Stop()

	var timeUsed time.Duration

	nops := 0

	for {
		<-ticker.C // wait for next tick

		for timeUsed < tickBudget {
			if debug {
				fmt.Printf("%8d %s\r\n", cycles, c.Debug())
			}
			op := c.Memory[c.PC]
			if op == 0x00 {
				nops++
			} else {
				nops = 0
			}
			c.PC++
			cyclesThisOp := ops[op](op, c)
			cycles += cyclesThisOp
			timeUsed += time.Duration(cyclesThisOp) * c.ClockTime

			if nops > 10 {
				panic("nop")
			}
		}

		timeUsed -= tickBudget
	}
}
