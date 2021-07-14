package hal

import (
	"fmt"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

var validGPIOPins = []int{
	rpi.GPIO2,
	rpi.GPIO3,
	rpi.GPIO5,
	rpi.GPIO6,
	rpi.GPIO7,
	rpi.GPIO8,
	rpi.GPIO9,
	rpi.GPIO10,
	rpi.GPIO11,
	rpi.GPIO12,
	rpi.GPIO13,
	rpi.GPIO14,
	rpi.GPIO15,
	rpi.GPIO16,
	rpi.GPIO17,
	rpi.GPIO18,
	rpi.GPIO19,
	rpi.GPIO20,
	rpi.GPIO21,
	rpi.GPIO22,
	rpi.GPIO23,
	rpi.GPIO24,
	rpi.GPIO25,
	rpi.GPIO26,
	rpi.GPIO27,
}

type Pin struct {
	number    int
	name      string
	lastState bool

	line *gpiod.Line
}

func (p *Pin) Name() string {
	return p.name
}
func (p *Pin) Number() int {
	return p.number
}

func (p *Pin) Close() error {
	return p.line.Close()
}

func (p *Pin) Read() (bool, error) {
	if err := p.line.Reconfigure(gpiod.AsInput); err != nil {
		return false, fmt.Errorf("can't read input from channel %d: %v", p.number, err)
	}
	v, err := p.line.Value()
	return v == 1, err
}

func (p *Pin) Write(state bool) error {
	if err := p.line.Reconfigure(gpiod.AsOutput()); err != nil {
		return fmt.Errorf("can't set output on channel %d: %v", p.number, err)
	}
	var value int
	if state {
		value = 1
	}
	p.lastState = state
	return p.line.SetValue(value)
}

func (p *Pin) LastState() bool {
	return p.lastState
}
