package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

func main() {
	time.Sleep(1 * time.Second)

	err := initLeds()
	if err != nil {
		println(err)
		return
	}
	println("leds initialized")

	err = initUsb()
	if err != nil {
		println(err)
		return
	}
	println("usb initialized")

	err = initEncryption()
	if err != nil {
		println(err)
		return
	}
	println("encryption initialized")
	println("session key:", base64.StdEncoding.EncodeToString(sessionKey))

	addr, err := initBLE()
	if err != nil {
		println(err)
		return
	}
	println("ble initialized, address:", addr.String())

	input := make([]byte, 6)
	ledVal := false

	go func() {
		for {
			err = bleConnect()
			if err != nil {
				println(err)
			}
		}
	}()

	for {
		usbRead6(input)
		print(string(input))

		ledVal = !ledVal
		setLed(led, ledVal)
	}
}

const bleConnectionInterval = 7500 * time.Microsecond

func bleConnect() error {
	dev, err := connect()
	if err != nil {
		return err
	}
	defer dev.Disconnect()
	println("device connected")

	effectiveMTU, err := dev.ExchangeMTU(47)
	if err != nil {
		return err
	}
	println("effective mtu:", effectiveMTU)

	keyChar, dataChar, err := getCharacteristics(dev)
	if err != nil {
		return err
	}

	_, err = keyChar.WriteWithoutResponse(encryptionKeyCiphertext)
	if err != nil {
		return err
	}
	println("key updated")

	var i, j, fails uint8

	buf := make([]byte, 6)

	nonce := make([]byte, 12)
	var nonce1 uint64
	var nonce2 uint32

	for {
		start := time.Now()

		buf[0] = j
		j++

		if nonce1 < math.MaxUint64 {
			nonce1++
		} else {
			nonce2++
		}
		binary.LittleEndian.PutUint64(nonce[0:8], nonce1)
		binary.LittleEndian.PutUint32(nonce[8:12], nonce2)

		// TODO: preallocate buffer
		ciphertext := aesGcm.Seal(nil, nonce, buf, nil)
		payload := append(ciphertext, nonce...)

		_, err = dataChar.WriteWithoutResponse(payload)
		if err != nil {
			fails++
			if fails > 100 {
				setLed(led, true)
				return err
			}
		} else {
			fails = 0

			i = (i + 1) % 24
			if i == 0 {
				setLed(led, false)
			} else if i > 12 {
				setLed(led, true)
			}
		}

		fmt.Println(".")

		dt := time.Since(start)
		if dt < bleConnectionInterval {
			time.Sleep(bleConnectionInterval - dt)
		}
	}
}
