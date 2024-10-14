package main

import (
	"machine"
)

var led uint8

func initLeds() error {
	var err error
	led, err = initLed(machine.LED)
	if err != nil {
		println("init led failed", err)
		return err
	}

	err = machine.PWM0.Configure(machine.PWMConfig{
		Period: 1e6,
	})
	if err != nil {
		println("pwm configure failed", err)
		return err
	}
	machine.PWM0.Set(led, 0)

	return nil
}

func initLed(led machine.Pin) (uint8, error) {
	ledCh, err := machine.PWM0.Channel(led)
	if err != nil {
		println("get pwm channel failed", err)
		return 0, err
	}
	machine.PWM0.SetInverting(ledCh, false)
	machine.PWM0.Set(ledCh, 0)

	return ledCh, nil
}

func setLed(led uint8, enable bool) {
	if enable {
		machine.PWM0.Set(led, machine.PWM0.Top()/1000)
	} else {
		machine.PWM0.Set(led, 0)
	}
}
