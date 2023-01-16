package common

import (
	"fmt"
	"math"
)

func AverageWithoutZero(windows []int64) float64 {
	var total int64
	var count float64
	for _, n := range windows {
		total += n
		if n != 0 {
			count++
		}
	}
	return float64(total) / count
}

func HashPower(power float64) string {
	if power < 1000 {
		return fmt.Sprintf("%vH", power)
	} else {
		return fmt.Sprintf("%.1fKH", power/1000.0)
	}
}

func Probability(p, diff, average float64) float64 {
	p = 1 - p
	expect := math.Floor(math.Log(p) / math.Log(1-1/diff))
	return expect / average
}
