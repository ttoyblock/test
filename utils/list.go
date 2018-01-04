package utils

import (
	"sort"
)

// RemoveEleOrder remove element by index and keep order
func RemoveEleOrder(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

// RemoveEle remove by index and don't care the order
func RemoveEle(slice []int, i int) []int {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

type IntSlice []int

func (sli IntSlice) Len() int {
	return len(sli)
}

func (sli IntSlice) Less(i, j int) bool {
	return sli[i] < sli[j]
}

func (sli IntSlice) Swap(i, j int) {
	sli[i], sli[j] = sli[j], sli[i]
}

func ListSort(sli []int, desc bool) {
	var tmp IntSlice = sli
	if desc {
		sort.Sort(sort.Reverse(tmp))
	} else {
		sort.Sort(tmp)
	}
}
