package usbfs

import (
	"os"
	"path/filepath"
	"strings"
)

type DeviceFilterFunc func(*Device) bool

func List(f DeviceFilterFunc) ([]*Device, error) {
	devices := []*Device{}

	if err := filepath.Walk("/sys/bus/usb/devices", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink == 0 || strings.Contains(info.Name(), ":") {
			return nil
		}

		dev, err := newDevice(path)
		if err != nil {
			return err
		}

		if f == nil || f(dev) {
			devices = append(devices, dev)
		}

		return nil

	}); err != nil {
		return nil, err
	}

	return devices, nil
}

func First(f DeviceFilterFunc) (*Device, error) {
	devs, err := List(f)
	if err != nil {
		return nil, err
	}

	if len(devs) == 0 {
		return nil, nil
	}

	return devs[0], nil
}
