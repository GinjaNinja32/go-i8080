package i8080

import (
	"fmt"
	"time"
)

const MHZ_2 = 500 * time.Nanosecond
const MHZ_3_125 = 320 * time.Nanosecond

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

type CPU struct {
	Memory    [65536]uint8
	Registers [8]uint8
	Flags     flags

	SP uint16
	PC uint16

	ClockTime time.Duration // nanoseconds per clock tick

	bdos
}

func New() (c *CPU) {
	c = &CPU{
		Flags:     F_BIT1_1,
		ClockTime: MHZ_2,
	}

	c.initBDOS()

	return
}

func (c *CPU) Output(s string) {
	//c.OutputStr += s
	fmt.Printf("%s", s)
}

const tickBudget = 10 * time.Millisecond

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
