package main

import (
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type testReorder struct {
	Current    int
	NumSamples int
	Series     []float32
	Expect     []float32
}

func sliceCompare(t *testing.T, a []float32, expect []float32) {
	if len(a) != len(expect) {
		t.Errorf("series should be same length, got %v; expected %v", a, expect)
		return
	}
	for i := 0; i < len(a); i++ {
		if a[i] != expect[i] {
			t.Errorf("got %v; expected %v", a, expect)
			return
		}
	}
}
func TestReorder(t *testing.T) {
	ths := []testReorder{
		testReorder{
			Current:    0,
			NumSamples: 0,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{},
		},
		testReorder{
			Current:    0,
			NumSamples: 1,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{1},
		},
		testReorder{
			Current:    0,
			NumSamples: 100,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{1},
		},
		testReorder{
			Current:    1,
			NumSamples: 1,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{2},
		},
		testReorder{
			Current:    1,
			NumSamples: 2,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{1, 2},
		},
		testReorder{
			Current:    1,
			NumSamples: 3,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{1, 2},
		},
		testReorder{
			Current:    4,
			NumSamples: 3,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{3, 4, 5},
		},
		testReorder{
			Current:    5,
			NumSamples: 3,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{4, 5, 1},
		},
		testReorder{
			Current:    6,
			NumSamples: 3,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{5, 1, 2},
		},
		testReorder{
			Current:    6,
			NumSamples: 5,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{3, 4, 5, 1, 2},
		},
		testReorder{
			Current:    7,
			NumSamples: 3,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{1, 2, 3},
		},
		testReorder{
			Current:    7,
			NumSamples: 5,
			Series:     []float32{1, 2, 3, 4, 5},
			Expect:     []float32{4, 5, 1, 2, 3},
		},
	}
	for _, th := range ths {
		res := reorderSeries(th.Series, th.Current, th.NumSamples)
		sliceCompare(t, res, th.Expect)

	}
}
