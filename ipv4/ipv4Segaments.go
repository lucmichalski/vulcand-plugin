package ipv4

import "sort"

// IPv4Segaments ========================================================
// Simple IPv4Segaments Struct, base on IPv4Segament
// Support sort.Interface
// ======================================================================
type IPv4Segaments []IPv4Segament

// Merge IPv4 Segament
// 1.1.1.1 - 1.1.1.5 and 1.1.1.3 - 1.1.1.7 will merge to 1.1.1.1 - 1.1.1.7
func (this IPv4Segaments) Merge() bool {
	if this.Len() <= 1 {
		return true
	}
	sort.Sort(this)
	var res IPv4Segaments
	var now = 0
	for i := range this {
		if i == 0 {
			res = append(res, this[i])
		} else {
			if res[now].IsInclude(this[i].Begin) {
				if res[now].End.Less(this[i].End) {
					res[now].End = this[i].End
				}
			} else {
				res = append(res, this[i])
				now++
			}
		}
	}
	this = res
	return true
}

// Support sort.Interface
func (this IPv4Segaments) Len() int {
	return len(this)
}

func (this IPv4Segaments) Less(i, j int) bool {
	return (this[i].Less(this[j]))
}

func (this IPv4Segaments) Swap(i, j int) {
	for k := range this[i].Begin {
		this[i].Begin[k], this[j].Begin[k] = this[j].Begin[k], this[i].Begin[k]
		this[i].End[k], this[j].End[k] = this[j].End[k], this[i].End[k]
	}
}

func (this IPv4Segaments) IsInclude(ip IPv4Addr) bool {
	var low = 0
	var end = this.Len() - 1
	for {
		var mid = low + (end-low)/2
		if low >= end {
			// Keep low and end are in safe range
			if low >= this.Len() {
				low = this.Len() - 1
			}
			if end < 0 {
				end = 0
			}
			return this[low].IsInclude(ip) || this[end].IsInclude(ip)
		}
		
		if this[mid].IsInclude(ip) {
			return true
		} else {
			if ip.Greater(this[mid].End) {
				low = mid + 1
			} else {
				end = mid - 1
			}
		}
	}
	return false
}
