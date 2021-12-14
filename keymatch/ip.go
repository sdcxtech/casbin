package keymatch

import (
	"fmt"
	"net"
)

// IPMatch determines whether IP address 'ip' matches the IP CIDR pattern.
//
// For example, "192.168.2.123" matches "192.168.2.0/24".
func IPMatch(ip string, ipCIDR string) (matched bool, err error) {
	objIP1 := net.ParseIP(ip)
	if objIP1 == nil {
		err = fmt.Errorf("argument 1 is not a valid IP address")
		return
	}

	_, cidr, err := net.ParseCIDR(ipCIDR)
	if err != nil {
		err = fmt.Errorf("argument 1 is not a valid IP CIDR address")
		return
	}

	matched = cidr.Contains(objIP1)
	return
}
