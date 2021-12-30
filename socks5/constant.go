package socks5

// Version socks5 version
const Version = 5

// method
const (
	MethodNone = 0 // NO AUTHENTICATION REQUIRED
)

// request cmd
const (
	CmdConnect      = 0x01
	CmdUDPAssociate = 0x03
)

// address type
const (
	AddressTypeIPV4   = 0x01
	AddressTypeDomain = 0x03
	AddressTypeIPV6   = 0x04
)
