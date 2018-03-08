package utils

import (
	"math/rand"
	"sort"
	"time"
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

func SortIntsList(sli []int, desc bool) {
	var tmp IntSlice = sli
	if desc {
		sort.Sort(sort.Reverse(tmp))
	} else {
		sort.Sort(tmp)
	}
}

// ListReverse reverse list
func ListReverse(yourSlice []int) {
	if len := len(yourSlice); len > 0 {
		lastIdx := len - 1
		for i := 0; i < len/2; i++ {
			yourSlice[i], yourSlice[lastIdx-i] = yourSlice[lastIdx-i], yourSlice[i]
			// fmt.Println(yourSlice)
		}
	}
}

func ShuffleInts(yourSlice []int) {
	t := time.Now()
	seed := int(t.Nanosecond()) // no shuffling without this line

	ShuffleIntsWithSeed(yourSlice, seed)
}

func ShuffleIntsWithSeed(yourSlice []int, seed int) {
	rand.Seed(int64(seed)) // no shuffling without this line

	for i := len(yourSlice) - 1; i > 0; i-- {
		j := rand.Intn(i)
		yourSlice[i], yourSlice[j] = yourSlice[j], yourSlice[i]
	}
}
