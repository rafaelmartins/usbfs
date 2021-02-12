package usbfs

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

func sysfsReadAsString(dir string, entry string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Join(dir, entry))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func sysfsReadAsUint(dir string, entry string, base int, bitSize int) (uint64, error) {
	v, err := sysfsReadAsString(dir, entry)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(v, base, bitSize)
}

func sysfsReadAsByte(dir string, entry string) (byte, error) {
	v, err := sysfsReadAsUint(dir, entry, 10, 8)
	return byte(v), err
}

func sysfsReadAsUint16(dir string, entry string) (uint16, error) {
	v, err := sysfsReadAsUint(dir, entry, 10, 16)
	return uint16(v), err
}

func sysfsReadAsHexByte(dir string, entry string) (byte, error) {
	v, err := sysfsReadAsUint(dir, entry, 16, 8)
	return byte(v), err
}

func sysfsReadAsHexUint16(dir string, entry string) (uint16, error) {
	v, err := sysfsReadAsUint(dir, entry, 16, 16)
	return uint16(v), err
}
