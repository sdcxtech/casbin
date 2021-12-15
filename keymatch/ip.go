package keymatch

import (
	"errors"
	"fmt"
	"net"
)

var ErrInvalidIP = errors.New("invalid textual representation of an IP address")

// IPMatch determines whether IP address 'ip' matches the IP CIDR pattern.
//
// For example, "192.168.2.123" matches "192.168.2.0/24".
func IPMatch(ip string, ipCIDR string) (matched bool, err error) {
	objIP1 := net.ParseIP(ip)
	if objIP1 == nil {
		err = fmt.Errorf("argument 1: %w", ErrInvalidIP)

		return
	}

	_, cidr, err := net.ParseCIDR(ipCIDR)
	if err != nil {
		err = fmt.Errorf("argument 2: %w", err)

		return
	}

	matched = cidr.Contains(objIP1)

	return
}
