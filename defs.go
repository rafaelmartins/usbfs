package usbfs

const (
	usbdevfsClaimInterface   = 0x8004550f
	usbdevfsReleaseInterface = 0x80045510
)

type Direction byte

const (
	DirectionOut Direction = (iota << 7)
	DirectionIn
)
