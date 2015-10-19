package ipv4

import "sort"

type IPv4Segaments []IPv4Segament

func (this IPv4Segaments) Merge() bool {
	sort.Sort(this)
	var res IPv4Segaments
	var now = 0
	for i := range this {
		if i == 0 {
			res = append(res, this[i])
		} else {
			if res[now].IsInclude(this[i].Begin) {
				if res[now].End.Greater(this[i].End) {
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
			if this[mid].End.Less(ip) {
				low = mid + 1
			} else {
				end = mid - 1
			}
		}
	}
	return false
}
