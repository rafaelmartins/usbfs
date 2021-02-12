package usbfs

import (
	"time"
	"unsafe"

	"github.com/rafaelmartins/usbfs/internal/endianness"
)

type RequestType byte

const (
	RequestTypeStandard RequestType = (iota << 5)
	RequestTypeClass
	RequestTypeVendor
)

type RequestRecipient byte

const (
	RequestRecipientDevice RequestRecipient = (iota << 0)
	RequestRecipientInterface
	RequestRecipientEndpoint
	RequestRecipientOther
)

type ctrlReq struct {
	ReqType uint8
	Req     uint8
	Value   uint16
	Index   uint16
	Len     uint16
	Timeout uint32
	Data    uintptr
}

func (d *Device) Control(typ RequestType, rcpt RequestRecipient, dir Direction, req byte, val uint16, idx uint16, data []byte, timeout time.Duration) error {
	var dataPointer uintptr
	if len(data) > 0 {
		dataPointer = uintptr(unsafe.Pointer(&data[0]))
	}

	creq := &ctrlReq{
		ReqType: byte(typ) | byte(rcpt) | byte(dir),
		Req:     req,
		Value:   endianness.ToLittleEndian16(val),
		Index:   endianness.ToLittleEndian16(idx),
		Len:     uint16(len(data)),
		Timeout: uint32(timeout.Milliseconds()),
		Data:    dataPointer,
	}
	_, err := d.ioctl(usbdevfsControl, uintptr(unsafe.Pointer(creq)))
	return err
}
