package usbfs

import (
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
	// TODO: validate endpoints against d.interfaces

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
