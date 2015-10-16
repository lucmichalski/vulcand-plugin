package ipv4

import (
	"fmt"
	"sort"
)

// IPv4Segaments ========================================================
// Simple IPv4Segaments Struct, base on IPv4Segament
// Support sort.Interface
// ======================================================================
type IPv4Segaments []IPv4Segament

func NewIPv4SegamentsFromOther(other IPv4Segaments) IPv4Segaments {
	var res IPv4Segaments
	for _, v := range other {
		res = append(res, v)
	}
	return res
}

// Support sort.Interface
func (this IPv4Segaments) Len() int {
	return len(this)
}

func (this IPv4Segaments) Less(i, j int) bool {
	return this[i].Less(this[j])
}

func (this IPv4Segaments) Swap(i, j int) {
	for k := range this[i].begin {
		this[i].begin[k], this[j].begin[k] = this[j].begin[k], this[i].begin[k]
		this[i].end[k], this[j].end[k] = this[j].end[k], this[i].end[k]
	}
}

func (this IPv4Segaments) IsInclude(ip IPv4Addr) bool {
	var begin = 0
	var end = len(this) - 1
	fmt.Println("Is Include ", this)

	for {
		var mid = begin + (end-begin)/2

		if begin >= end {
			if begin > len(this)-1 {
				begin = len(this) - 1
			}
			if end < 0 {
				end = 0
			}
			return this[begin].IsInclude(ip) || this[end].IsInclude(ip)
		}
		if this[mid].IsInclude(ip) {
			return true
		} else {
			if ip.Greater(this[mid].end) {
				begin = mid + 1
			} else {
				end = mid - 1
			}
		}
	}
	return false
}

// Merge IPv4 Segament
// 1.1.1.1 - 1.1.1.5 and 1.1.1.3 - 1.1.1.7 will merge to 1.1.1.1 - 1.1.1.7
func IPv4SegamentsMerge(ipsms IPv4Segaments) IPv4Segaments {
	if ipsms.Len() <= 1 {
		return ipsms
	}

	sort.Sort(ipsms)
	var res IPv4Segaments
	var now = 0
	for i := range ipsms {
		if i == 0 {
			res = append(res, ipsms[0])
			continue
		}

		if res[now].IsInclude(ipsms[i].begin) {
			if res[now].end.Less(ipsms[i].end) {
				res[now].end = ipsms[i].end
			}
		} else {
			res = append(res, ipsms[i])
			now++
		}
	}
	return res
}
