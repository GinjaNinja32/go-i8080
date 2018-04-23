package i8080

import (
	"fmt"
	"time"
)

// Inverse speed constants; Speed2MHz per instruction is a speed of 2MHz, etc
const (
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

	bdos
}

// New creates a new emulated Intel 8080 CPU
func New() (c *CPU) {
	c = &CPU{
		Flags:     FlagBit1,
		ClockTime: Speed2Mhz,
	}

	c.initBDOS()

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
	ticker := time.NewTicker(tickBudget)

	defer func() {
		f := recover()
		if f != nil && f != "hlt" {
			panic(f)
		}
	}()
	defer ticker.Stop()

	var timeUsed time.Duration

	for {
		<-ticker.C // wait for next tick

		for timeUsed < tickBudget {
			op := c.Memory[c.PC]
			c.PC++
			cyclesThisOp := ops[op](op, c)
			cycles += cyclesThisOp
			timeUsed += time.Duration(cyclesThisOp) * c.ClockTime
		}

		timeUsed -= tickBudget
	}
}
