package common

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{4 * time.Second, "4s"},
		{90 * time.Minute, "1h30m0s"},
		{23*time.Hour + 59*time.Minute + 59*time.Second, "23h59m59s"},
		{24 * time.Hour, "1d0h0m0s"},
		{86*24*time.Hour + 39*time.Minute + 21*time.Second, "86d0h39m21s"},
		{821*time.Hour + time.Minute + 59*time.Second, "34d5h1m59s"},
		{2727*time.Hour + 24*time.Minute + 47*time.Second, "113d15h24m47s"},
		{36*time.Hour + 500*time.Millisecond, "1d12h0m1s"},
	}
	for _, tt := range tests {
		if got := FormatDuration(tt.d); got != tt.want {
			t.Errorf("FormatDuration(%v) = %s, want %s", tt.d, got, tt.want)
		}
	}
}

func TestHashPower(t *testing.T) {
	tests := []struct {
		power float64
		want  string
	}{
		{823.3333333, "823.3H"},
		{1000, "1.0KH"},
		{515712.4, "515.7KH"},
	}
	for _, tt := range tests {
		if got := HashPower(tt.power); got != tt.want {
			t.Errorf("HashPower(%v) = %s, want %s", tt.power, got, tt.want)
		}
	}
}
