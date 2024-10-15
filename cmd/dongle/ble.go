package main

import (
	"time"

	"github.com/google/uuid"
	"tinygo.org/x/bluetooth"
)

var deviceAddress string

var serviceUUID = bluetooth.NewUUID(uuid.MustParse("1afd9a71-465c-4b37-847a-67d4644dbc30"))
var keyCharacteristicUUID = bluetooth.NewUUID(uuid.MustParse("210d0f78-6ca7-4bca-b9ed-dfa70bfba9f8"))
var dataCharacteristicUUID = bluetooth.NewUUID(uuid.MustParse("313b5804-5030-484c-ab33-e11d4938c45f"))

func initBLE() (bluetooth.MACAddress, error) {
	err := bluetooth.DefaultAdapter.Configure(bluetooth.Config{
		Gap: bluetooth.GapConfig{
			EventLength: 3,
		},
		Gatt: bluetooth.GattConfig{
			AttMtu: 47,
		},
	})
	if err != nil {
		return bluetooth.MACAddress{}, err
	}

	err = bluetooth.DefaultAdapter.Enable()
	if err != nil {
		return bluetooth.MACAddress{}, err
	}

	addr, err := bluetooth.DefaultAdapter.Address()
	if err != nil {
		return bluetooth.MACAddress{}, err
	}

	return addr, nil
}

func connect() (bluetooth.Device, error) {
	ch := make(chan bluetooth.ScanResult, 1)
	err := bluetooth.DefaultAdapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.Address.String() == deviceAddress {
			_ = adapter.StopScan()
			ch <- result
		}
	})
	if err != nil {
		return bluetooth.Device{}, err
	}

	var dev bluetooth.Device
	select {
	case result := <-ch:
		println("device found")
		dev, err = bluetooth.DefaultAdapter.Connect(result.Address, bluetooth.ConnectionParams{
			ConnectionTimeout: bluetooth.NewDuration(4 * time.Second),
			MinInterval:       bluetooth.NewDuration(7500 * time.Microsecond),
			MaxInterval:       bluetooth.NewDuration(7500 * time.Microsecond),
		})
		if err != nil {
			return bluetooth.Device{}, err
		}
	}

	return dev, nil
}

func getCharacteristics(dev bluetooth.Device) (bluetooth.DeviceCharacteristic, bluetooth.DeviceCharacteristic, error) {
	svcs, err := dev.DiscoverServices([]bluetooth.UUID{serviceUUID})
	if err != nil || len(svcs) == 0 {
		return bluetooth.DeviceCharacteristic{}, bluetooth.DeviceCharacteristic{}, err
	}

	chars, err := svcs[0].DiscoverCharacteristics([]bluetooth.UUID{
		keyCharacteristicUUID,
		dataCharacteristicUUID,
	})
	if err != nil || len(chars) == 0 {
		return bluetooth.DeviceCharacteristic{}, bluetooth.DeviceCharacteristic{}, err
	}

	return chars[0], chars[1], nil
}
