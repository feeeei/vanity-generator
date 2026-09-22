package common

import (
	"fmt"
	"math"
	"time"
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
		return fmt.Sprintf("%.1fH", power)
	} else {
		return fmt.Sprintf("%.1fKH", power/1000.0)
	}
}

func Probability(p, diff, average float64) float64 {
	p = 1 - p
	expect := math.Floor(math.Log(p) / math.Log(1-1/diff))
	return expect / average
}

// FormatDuration 把时长格式化为 34d5h1m59s，有天数时小时和分钟固定输出；不足一天时与 time.Duration 的默认格式一致
func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	if d < 24*time.Hour {
		return d.String()
	}
	days := d / (24 * time.Hour)
	rest := d % (24 * time.Hour)
	return fmt.Sprintf("%dd%dh%dm%ds", days, rest/time.Hour, rest%time.Hour/time.Minute, rest%time.Minute/time.Second)
}
