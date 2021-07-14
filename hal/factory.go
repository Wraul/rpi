package hal

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/warthog618/gpiod"

	"github.com/reef-pi/hal"

	"github.com/wraul/rpi/pwm"
)

type Factory struct {
	meta       hal.Metadata
	parameters []hal.ConfigParameter
}

type pinFactory func(int) (Pin, error)

var factory *Factory
var once sync.Once

// RpiFactory provides the factory to get RPI Driver parameters and RPI Drivers
func RpiFactory() *Factory {
	once.Do(func() {
		factory = &Factory{
			meta: hal.Metadata{
				Name:         "rpi",
				Description:  "hardware peripherals and GPIO channels on the base raspberry pi hardware",
				Capabilities: []hal.Capability{hal.DigitalInput, hal.DigitalOutput, hal.PWM},
			},
			parameters: []hal.ConfigParameter{
				{
					Name:    "Frequency",
					Type:    hal.Integer,
					Order:   0,
					Default: "200",
				},
				{
					Name:    "GPIO Device",
					Type:    hal.String,
					Order:   1,
					Default: "gpiochip0",
				},
			},
		}
	})
	return factory
}

func (f *Factory) GetParameters() []hal.ConfigParameter {
	return f.parameters
}

func (f *Factory) ValidateParameters(parameters map[string]interface{}) (bool, map[string][]string) {

	var failures = make(map[string][]string)

	var v interface{}
	var ok bool

	if v, ok = parameters["Frequency"]; ok {
		_, ok := hal.ConvertToInt(v)
		if !ok {
			failure := fmt.Sprint("Frequency is not a number. ", v, " was received.")
			failures["Frequency"] = append(failures["Frequency"], failure)
		}
	} else {
		failure := fmt.Sprint("Frequency is required parameter, but was not received.")
		failures["Frequency"] = append(failures["Frequency"], failure)
	}

	if v, ok = parameters["GPIO Device"]; ok {
		path := path.Join("/dev/", v.(string))
		_, err := os.Stat(path)
		if err != nil {
			failure := fmt.Sprintf("Invalid GPIO Device %s. %v", path, err)
			failures["GPIO Device"] = append(failures["GPIO Device"], failure)
		}
	} else {
		failure := fmt.Sprint("GPIO Device is a required parameter, but was not received.")
		failures["GPIO Device"] = append(failures["GPIO Device"], failure)
	}

	return len(failures) == 0, failures
}

func (f *Factory) Metadata() hal.Metadata {
	return f.meta
}

func (f *Factory) NewDriver(parameters map[string]interface{}, hardwareResources interface{}) (hal.Driver, error) {
	if valid, failures := f.ValidateParameters(parameters); !valid {
		return nil, errors.New(hal.ToErrorString(failures))
	}

	frequency, _ := hal.ConvertToInt(parameters["Frequency"])
	gpioDev, _ := parameters["GPIO Device"]

	gpioChip, err := gpiod.NewChip(gpioDev.(string))
	if err != nil {
		return nil, fmt.Errorf("can't create GPIO chip: %v", err)
	}

	pwmDriver := pwm.New()

	driver := &Driver{
		pins:     make(map[int]*Pin),
		channels: make(map[int]*Channel),
		meta:     f.meta,
		chip:     gpioChip,
	}

	for _, i := range validGPIOPins {
		l, err := driver.chip.RequestLine(i)

		if err != nil {
			return nil, fmt.Errorf("can't build hal pin %d: %v", i, err)
		}
		name := fmt.Sprintf("GP%d", i)
		driver.pins[i] = &Pin{
			name:   name,
			number: i,
			line:   l,
		}
	}

	for _, p := range []int{0, 1} {
		ch := &Channel{
			pin:       p,
			driver:    pwmDriver,
			frequency: frequency,
			name:      fmt.Sprintf("%d", p),
		}
		driver.channels[p] = ch
	}
	return driver, nil
}
