package ipv4

import (
	"fmt"
	"strconv"
	"strings"
)

// ipv4 segament =====================================
// Simple IPv4Segament struct, base on IPv4Addr
// Support constructs, String, sort.Interface
// ===================================================
type IPv4Segament struct {
	begin IPv4Addr
	end   IPv4Addr
}

func NewIPv4SegamentFromOther(other IPv4Segament) (IPv4Segament, error) {
	return IPv4Segament{other.begin, other.end}, nil
}

func NewIPv4SegamentFromIPv4Addr(begin, end IPv4Addr) (IPv4Segament, error) {
	if end.Less(begin) {
		return IPv4Segament{}, fmt.Errorf("begin ip must Less(or equal) end ip")
	}
	return IPv4Segament{begin, end}, nil
}

func (this IPv4Segament) String() string {
	return fmt.Sprint("begin:", this.begin, ", end:", this.end)
}

func getIPv4MaskFromInt(mask int) (IPv4Addr, error) {
	if mask < 0 || mask > 32 {
		return IPv4Addr{}, fmt.Errorf("Mask must between 0 and 32: %v", mask)
	}
	var res IPv4Addr
	var now int
	for now = 0; now < (mask / 8); now++ {
		res[now] = 255
	}
	for mask %= 8; now < 4; now++ {
		switch mask {
		case 7:
			res[now] += (1 << 1)
			fallthrough
		case 6:
			res[now] += (1 << 2)
			fallthrough
		case 5:
			res[now] += (1 << 3)
			fallthrough
		case 4:
			res[now] += (1 << 4)
			fallthrough
		case 3:
			res[now] += (1 << 5)
			fallthrough
		case 2:
			res[now] += (1 << 6)
			fallthrough
		case 1:
			res[now] += (1 << 7)
			mask = 0
		case 0:
			res[now] = 0
		}
	}
	return res, nil
}

func toolGetIPv4SegamentBeginAddr(ip, mask IPv4Addr) (IPv4Addr, error) {
	res, err := NewIPv4AddrFromOther(ip)

	fmt.Println("GetIpv4Begin ", ip, mask)

	if err != nil {
		return res, err
	}
	for i := range mask {
		res[i] &= mask[i]
	}
	return res, nil
}

func NewIPv4SegamentFromString(str string) (IPv4Segament, error) {
	var res IPv4Segament

	switch {
	// format like this 1.1.1.1/24
	case strings.Contains(str, "/"):
		ipAndmask := strings.Split(str, "/")
		for i := range ipAndmask {
			ipAndmask[i] = strings.TrimSpace(ipAndmask[i])
		}
		ip, err := NewIPv4AddrFromString(ipAndmask[0])
		if err != nil {
			return IPv4Segament{}, err
		}
		maskInt, err := strconv.Atoi(ipAndmask[1])
		if err != nil {
			return IPv4Segament{}, fmt.Errorf("Not Support Mask Format: ", ipAndmask[1])
		}
		mask, err := getIPv4MaskFromInt(maskInt)
		if err != nil {
			return IPv4Segament{}, err
		}
		res.begin, err = toolGetIPv4SegamentBeginAddr(ip, mask)
		if err != nil {
			return res, err
		}
		for i := range mask {
			res.end[i] = res.begin[i] | (^(mask[i]))
		}

		return res, nil
	case strings.Contains(str, "-"):
		// TODO: Will Supported 1.1.1.1-2.2.2.2 format
		ipBeginAndEnd := strings.Split(str, "-")
		if len(ipBeginAndEnd) != 2 {
			return res, fmt.Errorf("Unsupoorted format: ", str)
		}
		tbegin, err := NewIPv4AddrFromString(ipBeginAndEnd[0])
		if err != nil {
			return res, err
		}
		tend, err := NewIPv4AddrFromString(ipBeginAndEnd[1])
		if err != nil {
			return res, err
		}
		return NewIPv4SegamentFromIPv4Addr(tbegin, tend)
	default:
		// Just one IPv4 like this: 1.1.1.1
		var err error
		res.begin, err = NewIPv4AddrFromString(str)
		if err != nil {
			return res, err
		}
		res.end = res.begin
		return res, nil
	}

}

// Less compare begin ip at first, if equal then compare end ip
func (this IPv4Segament) Less(other IPv4Segament) bool {
	if this.begin.Equal(other.begin) {
		return this.end.Less(other.end)
	}
	return this.begin.Less(other.end)
}

// Greater compare begin ip at first, if equal then compare end ip
func (this IPv4Segament) Greater(other IPv4Segament) bool {
	if this.begin.Equal(other.begin) {
		return this.end.Greater(other.end)
	}
	return this.begin.Greater(other.begin)
}

func (this IPv4Segament) IsInclude(ip IPv4Addr) bool {
	if ip.Less(this.begin) || ip.Greater(this.end) {
		return false
	}
	return true
}
