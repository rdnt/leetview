package main

import (
	"fmt"
	"time"
)

func main() {
	err := initLeds()
	if err != nil {
		println(err.Error())
		return
	}

	err = initUsb()
	if err != nil {
		println(err.Error())
		return
	}

	setLed(led, true)
	time.Sleep(1 * time.Second)
	setLed(led, false)

	input := make([]byte, 6)
	ledVal := false

	for {
		usbRead6(input)
		fmt.Print(string(input))

		ledVal = !ledVal
		setLed(led, ledVal)
	}
}
