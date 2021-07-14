package hal

import (
	"fmt"
	"sort"

	"github.com/warthog618/gpiod"

	"github.com/reef-pi/hal"

	"github.com/wraul/rpi/pwm"
)

type Driver struct {
	meta      hal.Metadata
	pins      map[int]*Pin
	channels  map[int]*Channel
	pwmDriver pwm.Driver
	chip      *gpiod.Chip
}

func (d *Driver) Metadata() hal.Metadata {
	return d.meta
}

func (d *Driver) Close() error {
	for _, p := range d.pins {
		if err := p.Close(); err != nil {
			return fmt.Errorf("can't close hal driver due to channel %s: %v", p.Name(), err)
		}
	}

	if err := d.chip.Close(); err != nil {
		return fmt.Errorf("can't close hal driver due to chip: %v", err)
	}

	return nil
}

func (d *Driver) Pins(cap hal.Capability) ([]hal.Pin, error) {
	var pins []hal.Pin
	switch cap {
	case hal.DigitalInput, hal.DigitalOutput:
		for _, p := range d.GPIOPins() {
			pins = append(pins, p)
		}
		return pins, nil
	case hal.PWM:
		for _, p := range d.PWMChannels() {
			pins = append(pins, p)
		}
		return pins, nil
	default:
		return nil, fmt.Errorf("Unsupported capability:%s", cap.String())
	}
}

func (d *Driver) GPIOPins() []*Pin {
	var pins []*Pin
	for _, p := range d.pins {
		pins = append(pins, p)
	}
	sort.Slice(pins, func(i, j int) bool { return pins[i].Name() < pins[j].Name() })
	return pins
}

func (d *Driver) GPIOPin(p int) (*Pin, error) {
	pin, ok := d.pins[p]
	if !ok {
		return nil, fmt.Errorf("pin %d unknown", p)
	}
	return pin, nil
}

func (d *Driver) DigitalInputPins() []*Pin {
	return d.GPIOPins()
}

func (d *Driver) DigitalInputPin(p int) (*Pin, error) {
	return d.GPIOPin(p)
}

func (d *Driver) DigitalOutputPins() []*Pin {
	return d.GPIOPins()
}

func (d *Driver) DigitalOutputPin(p int) (*Pin, error) {
	return d.GPIOPin(p)
}

func (d *Driver) PWMChannels() []*Channel {
	var chs []*Channel
	for _, ch := range d.channels {
		chs = append(chs, ch)
	}
	sort.Slice(chs, func(i, j int) bool { return chs[i].Name() < chs[j].Name() })
	return chs
}

func (d *Driver) PWMChannel(p int) (*Channel, error) {
	ch, ok := d.channels[p]
	if !ok {
		return nil, fmt.Errorf("unknown pwm channel %d", p)
	}
	return ch, nil
}
