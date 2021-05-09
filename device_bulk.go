package usbfs

import (
	"fmt"
	"time"
	"unsafe"
)

type bulkReq struct {
	Ep      uint32
	Len     uint32
	Timeout uint32
	Data    uintptr
}

func (d *Device) Bulk(ep uint32, dir Direction, data []byte, timeout time.Duration) error {
	found := false
	for itf, obj := range d.interfaces {
		for _, e := range obj.endpoints {
			if ep == uint32(e) {
				found = true
				break
			}
		}
		if found {
			if obj.isOpen {
				break
			}
			if err := d.claim(uint32(itf)); err != nil {
				return err
			}
			obj.isOpen = true
			break
		}
	}
	if !found {
		return fmt.Errorf("usbfs: device_build: endpoint not found: 0x%02x", ep)
	}

	var dataPointer uintptr
	if len(data) > 0 {
		dataPointer = uintptr(unsafe.Pointer(&data[0]))
	}

	req := &bulkReq{
		Ep:      ep | uint32(dir),
		Len:     uint32(len(data)),
		Timeout: uint32(timeout.Milliseconds()),
		Data:    dataPointer,
	}
	_, err := d.ioctl(usbdevfsBulk, uintptr(unsafe.Pointer(req)))
	return err
}
