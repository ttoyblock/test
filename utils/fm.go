package utils

import "math"

// RoundTool - 四舍五入
// params:
// val - 需要被 round 的值
// roundOn - 进位方式 例如：0.5 表示四舍五入
// places - 从 val 第 0 位开始，保留的位置 例如: (0.55555, 0.5, 2) 得到 0.560000
func RoundTool(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
