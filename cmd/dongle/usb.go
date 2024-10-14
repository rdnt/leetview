package main

import (
	"machine"
	"time"
)

var (
	uart = machine.Serial
	tx   = machine.UART_TX_PIN
	rx   = machine.UART_RX_PIN
)

func initUsb() error {
	err := uart.Configure(machine.UARTConfig{TX: tx, RX: rx})
	if err != nil {
		return err
	}

	return nil
}

func usbRead6(b []byte) {
	i := 0
	for {
		if uart.Buffered() > 0 {
			data, _ := uart.ReadByte()

			b[i] = data
			i++

			if i == 6 {
				return
			}
		}
		time.Sleep(1 * time.Microsecond)
	}
}
