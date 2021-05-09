package usbfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/unix"
)

type interf struct {
	endpoints []byte
	isOpen    bool
}

type Device struct {
	isOpen     bool
	sysfspath  string
	devpath    string
	fd         int
	interfaces map[byte]*interf

	busnum       uint16
	devnum       uint16
	bcdDevice    uint16
	idVendor     uint16
	idProduct    uint16
	manufacturer string
	product      string
	serial       string
}

func newDevice(path string) (*Device, error) {
	dev := &Device{
		sysfspath:  path,
		interfaces: map[byte]*interf{},
	}

	itfs, err := filepath.Glob(filepath.Join(path, "?-*:?.*"))
	if err != nil {
		return nil, err
	}

	for _, itf := range itfs {
		bInterfaceNumber, err := sysfsReadAsHexByte(itf, "bInterfaceNumber")
		if err != nil {
			return nil, err
		}

		eps, err := filepath.Glob(filepath.Join(itf, "ep_*"))
		if err != nil {
			return nil, err
		}

		endpoints := []byte{}
		for _, ep := range eps {
			bEndpointAddress, err := sysfsReadAsHexByte(ep, "bEndpointAddress")
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, bEndpointAddress)
		}

		dev.interfaces[bInterfaceNumber] = &interf{
			endpoints: endpoints,
			isOpen:    false,
		}
	}

	dev.busnum, err = sysfsReadAsUint16(path, "busnum")
	if err != nil {
		return nil, err
	}

	dev.devnum, err = sysfsReadAsUint16(path, "devnum")
	if err != nil {
		return nil, err
	}

	dev.devpath = fmt.Sprintf("/dev/bus/usb/%03d/%03d", dev.busnum, dev.devnum)

	return dev, nil
}

func (d *Device) BcdDevice() (uint16, error) {
	if d.bcdDevice != 0 {
		return d.bcdDevice, nil
	}
	bcdDevice, err := sysfsReadAsHexUint16(d.sysfspath, "bcdDevice")
	if err == nil {
		d.bcdDevice = bcdDevice
	}
	return bcdDevice, err
}

func (d *Device) IdVendor() (uint16, error) {
	if d.idVendor != 0 {
		return d.idVendor, nil
	}
	idVendor, err := sysfsReadAsHexUint16(d.sysfspath, "idVendor")
	if err == nil {
		d.idVendor = idVendor
	}
	return idVendor, err
}

func (d *Device) IdProduct() (uint16, error) {
	if d.idProduct != 0 {
		return d.idVendor, nil
	}
	idProduct, err := sysfsReadAsHexUint16(d.sysfspath, "idProduct")
	if err == nil {
		d.idProduct = idProduct
	}
	return idProduct, err
}

func (d *Device) Manufacturer() (string, error) {
	if d.manufacturer != "" {
		return d.manufacturer, nil
	}
	manufacturer, err := sysfsReadAsString(d.sysfspath, "manufacturer")
	if err == nil {
		d.manufacturer = manufacturer
	}
	if os.IsNotExist(err) {
		return "", nil
	}
	return manufacturer, err
}

func (d *Device) Product() (string, error) {
	if d.product != "" {
		return d.product, nil
	}
	product, err := sysfsReadAsString(d.sysfspath, "product")
	if err == nil {
		d.product = product
	}
	if os.IsNotExist(err) {
		return "", nil
	}
	return product, err
}

func (d *Device) Serial() (string, error) {
	if d.serial != "" {
		return d.serial, nil
	}
	serial, err := sysfsReadAsString(d.sysfspath, "serial")
	if err == nil {
		d.serial = serial
	}
	if os.IsNotExist(err) {
		return "", nil
	}
	return serial, err
}

func (d *Device) Open() error {
	if d.isOpen {
		return errors.New("usbfs: device: already open")
	}

	fd, err := unix.Open(d.devpath, unix.O_RDWR, 0600)
	if err != nil {
		return err
	}

	d.fd = fd
	d.isOpen = true

	return nil
}

func (d *Device) Close() error {
	if !d.isOpen {
		return errors.New("usbfs: device: not open")
	}

	// closing the fd releases the interfaces
	if err := unix.Close(d.fd); err != nil {
		return err
	}

	d.isOpen = false
	return nil
}

func (d *Device) ioctl(request uint32, data uintptr) (int, error) {
	if !d.isOpen {
		return 0, errors.New("usbfs: device: not open")
	}

	rv, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(d.fd), uintptr(request), data)
	if errno != 0 {
		return 0, fmt.Errorf("usbfs: device: ioctl: 0x%x: %s", request, errno)
	}
	return int(rv), nil
}

func (d *Device) claim(itf uint32) error {
	_, err := d.ioctl(usbdevfsClaimInterface, uintptr(unsafe.Pointer(&itf)))
	return err
}
