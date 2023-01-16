package sliding_window

import "vanity-generator/common"

type SlidingWindow struct {
	windows []int64
	index   int

	cacheIndex   int
	cacheAverage float64
}

func NewSlidingWindow(size int) *SlidingWindow {
	return &SlidingWindow{
		windows: make([]int64, size),
	}
}

func (s *SlidingWindow) Add(count int64) {
	s.windows[s.index] = count
	s.index = (s.index + 1) % len(s.windows)
}

func (s *SlidingWindow) Average() float64 {
	// single flight
	if s.index == s.cacheIndex {
		return s.cacheAverage
	}

	average := common.AverageWithoutZero(s.windows)
	s.cacheIndex, s.cacheAverage = s.index, average
	return average
}
