package ipv4

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"github.com/YunxiangHuang/vulcand-plugin/utils"
)

// Include IPv4, IPv4Segament, IPv4Segaments

// IPv4 ================================================
// Simple IPv4 Struct
// Support Constructs, String
// =====================================================

type IPv4Addr [4]uint8

func (this IPv4Addr) Equal(other IPv4Addr) bool {
	for i := range this {
		if this[i] != other[i] {
			return false
		}
	}
	return true
}

func (this IPv4Addr) Greater(other IPv4Addr) bool {
	for i := range this {
		if this[i] != other[i] {
			return this[i] > other[i]
		}
	}
	// this Equal other, So NOT Greater
	return false
}

func (this IPv4Addr) Less(other IPv4Addr) bool {
	for i := range this {
		if this[i] != other[i] {
			return this[i] < other[i]
		}
	}
	// this Equal other, So NOT Less
	return false
}

func (this IPv4Addr) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", this[0], this[1], this[2], this[3])
}

func NewIPv4AddrFromOther(other IPv4Addr) (IPv4Addr, error) {
	return IPv4Addr{other[0], other[1], other[2], other[3]}, nil
}

func NewIPv4AddrFromString(str string) (IPv4Addr, error) {
	var res IPv4Addr
	tmp := utils.SplitWithoutSpace(str, ".")
	for i := range tmp {
		tmp[i] = strings.TrimSpace(tmp[i])
	}
	if strings.Contains(str, "/") || len(tmp) != 4 {
		return IPv4Addr{}, fmt.Errorf("Not Support IP format: ", str)
	}
	for i := range tmp {
		ipslice, err := strconv.Atoi(tmp[i])
		if err != nil || ipslice > 255 || ipslice < 0 {
			return IPv4Addr{}, fmt.Errorf("Not a leagal IP: ", str)
		}
		res[i] = uint8(ipslice)
	}
	return res, nil
}

func (this IPv4Addr) ToNetIP() net.IP {
	return net.IPv4(this[0], this[1], this[2], this[3])
}
