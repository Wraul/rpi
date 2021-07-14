package hal

import (
	"fmt"
	"log"

	"github.com/wraul/rpi/pwm"
)

type Channel struct {
	pin       int
	name      string
	driver    pwm.Driver
	frequency int
	v         float64
}

func (ch Channel) Set(value float64) error {
	if ch.frequency <= 0 {
		log.Printf("warning: RPI PWM frequency is 0, defaulting to 150")
		ch.frequency = 150
	}
	if value < 0 || value > 100 {
		return fmt.Errorf("value must be 0-100, got %f", value)
	}

	exported, err := ch.driver.IsExported(ch.pin)
	if err != nil {
		return err
	}
	if !exported {
		if err := ch.driver.Export(ch.pin); err != nil {
			return err
		}
	}
	if err := ch.driver.Frequency(ch.pin, ch.frequency); err != nil {
		return err
	}
	if err := ch.driver.DutyCycle(ch.pin, value); err != nil {
		return err
	}
	if err := ch.driver.Enable(ch.pin); err != nil {
		return err
	}
	ch.v = value
	return nil
}

func (ch *Channel) Close() error { return nil }
func (ch *Channel) LastState() bool {
	return ch.v == 100
}

func (ch *Channel) Write(b bool) error {
	var v float64
	if b == true {
		v = 100
	}
	return ch.Set(v)
}

func (ch *Channel) Name() string {
	return ch.name
}

func (ch *Channel) Number() int {
	return ch.pin
}
