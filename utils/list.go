package utils

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
