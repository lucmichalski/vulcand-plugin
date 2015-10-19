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
	Begin IPv4Addr
	End   IPv4Addr
}

func splitWithoutSpace(str string, flag string) []string {
		res := strings.Split(str, flag)
	for i := range res {
		res[i] = strings.TrimSpace(res[i])
	}
	return res
}

// Less compare begin ip at first, if equal then compare end ip
func (this IPv4Segament) Less(other IPv4Segament) bool {
	if this.Begin.Equal(other.Begin) {
		return this.End.Less(other.End)
	}
	return this.Begin.Less(other.Begin)
}

// Greater compare begin ip at first, if equal then compare end ip
func (this IPv4Segament) Greater(other IPv4Segament) bool {
	if this.Begin.Equal(other.Begin) {
		return this.End.Greater(other.End)
	}
	return this.Begin.Greater(other.Begin)
}

func (this IPv4Segament) Equal(other IPv4Segament) bool {
	return this.Begin.Equal(other.Begin) && this.End.Equal(other.End)
}

func (this IPv4Segament) String() string {
	return fmt.Sprint(this.Begin, "-", this.End)
}

func NewIPv4SegamentFromIPv4Addr(begin, end IPv4Addr) (IPv4Segament, error) {
	if end.Less(begin) {
		return IPv4Segament{}, fmt.Errorf("Begin IP must Less (or equal) End IP")
	}
	return IPv4Segament{begin, end}, nil
}

func NewIPv4SegamentFromOther(other IPv4Segament) IPv4Segament {
	return IPv4Segament{
		other.Begin, 
		other.End,
		}
}

func IPv4IntToMask(mint int) (IPv4Addr, error) {
	if mint < 0 || mint > 32 {
		return IPv4Addr{}, fmt.Errorf("mask must between 0 and 32: ", mint)
	}
	
	var now int
	var res IPv4Addr
	for now = 0; now < (mint / 8); now++ {
		res[now] = 255
	}
	for mint %= 8; now < 4; now++ {
		for ; mint > 0; mint-- {
			res[now] += (1 << (uint8(8 - mint)))
		}
	}
	return res, nil
}

func NewIPv4SegamentFromIPandMask(ip, mask IPv4Addr) (IPv4Segament, error) {
	begin := ip
	end := ip
	for i := 0; i < 4; i++ {
		begin[i] &= mask[i]
		end[i] ^= (^mask[i])
	}
	return IPv4Segament{
		Begin: begin,
		End: end,
	}, nil
}

func NewIPv4SegamentFromString(str string) (IPv4Segament, error) {
	switch {
	case strings.Contains(str, "/"):
		strs := splitWithoutSpace(str, "/")
		if len(strs) != 2 {
			return IPv4Segament{}, fmt.Errorf("Unsupport IP segament format: ", str)
		}
		tmpIP, err := NewIPv4AddrFromString(strs[0])
		if err != nil {
			return IPv4Segament{}, err
		}
		tmpMint, err := strconv.Atoi(strs[1])
		if err != nil {
			return IPv4Segament{}, fmt.Errorf("Unsupport IP segament format: ", str)
		}
		tmpMask, err := IPv4IntToMask(tmpMint)
		if err != nil {
			return IPv4Segament{}, err
		}
		return NewIPv4SegamentFromIPandMask(tmpIP, tmpMask)

	case strings.Contains(str, "-"):
		strs := splitWithoutSpace(str, "-")
		if len(strs) != 2 {
			return IPv4Segament{}, fmt.Errorf("Unsupport IP segament format: ", str)
		}
		tmpBegin, err := NewIPv4AddrFromString(strs[0])
		if err != nil {
			return IPv4Segament{}, err
		}
		tmpEnd, err := NewIPv4AddrFromString(strs[1])
		if err != nil {
			return IPv4Segament{}, err
		}
		return NewIPv4SegamentFromIPv4Addr(tmpBegin, tmpEnd)
	default:
		tmpIP, err := NewIPv4AddrFromString(str)
		if err != nil {
			return IPv4Segament{}, err
		}
		return NewIPv4SegamentFromIPv4Addr(tmpIP, tmpIP)
	}
}

func (this IPv4Segament) IsInclude(ip IPv4Addr) bool {
	return !(ip.Less(this.Begin) || ip.Greater(this.End))
}
