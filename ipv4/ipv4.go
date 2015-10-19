package ipv4

import (
	"fmt"
	"net"
)

// Based on net.IP
type IPv4Addr net.IP

func (this IPv4Addr) Equal(other IPv4Addr) bool {
	return net.IP(this).Equal(net.IP(other))
}

func (this IPv4Addr) Greater(other IPv4Addr) bool {
	for i := range this {
		if this[i] != other[i] {
			return this[i] > other[i]
		}
	}
	return false
}

func (this IPv4Addr) Less(other IPv4Addr) bool {
	for i := range this {
		if this[i] != other[i] {
			return this[i] < other[i]
		}
	}
	return false
}

func (this IPv4Addr) String() string {
	return net.IP(this).String()
}

func NewIPv4AddrFromOther(other IPv4Addr) IPv4Addr {
	return IPv4Addr(net.IPv4(other[0], other[1], other[2], other[3]))
}

func NewIPv4AddrFromString(str string) (IPv4Addr, error) {
	res := net.ParseIP(str).To4()
	if res == nil {
		return nil, fmt.Errorf("IP illegal: ", str)
	}
	return IPv4Addr(res), nil
}
