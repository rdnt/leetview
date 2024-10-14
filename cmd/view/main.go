package main

import (
	"time"
)

func main() {
	err := initLeds()
	if err != nil {
		println(err.Error())
		return
	}

	setLed(led, true)
	time.Sleep(1 * time.Second)
	setLed(led, false)

	ledVal := false

	for {
		time.Sleep(1 * time.Second)
		ledVal = !ledVal
		setLed(led, ledVal)
	}
}
