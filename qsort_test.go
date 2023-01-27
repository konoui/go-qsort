package qsort

import (
	"reflect"
	"testing"
)

func Test_Slice(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		noShuffle bool
		heap      bool
	}{
		{
			name:  "len <= 7",
			input: makeInt(0, 7),
		},
		{
			name:  "7 <= len < 40",
			input: makeInt(0, 9),
		},
		// {
		// 	name:      "40 < len",
		// 	input:     makeInt(0, 43),
		// 	noShuffle: true,
		// },
		{
			name:  "heap: len <= 7",
			input: makeInt(0, 7),
			heap:  true,
		},
		{
			name:  "heap: 7 <= len < 40",
			input: makeInt(0, 9),
			heap:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := make([]int, len(tt.input))
			copy(want, tt.input)

			inputs := [][]int{tt.input}
			if !tt.noShuffle {
				inputs = shuffle(tt.input)
			}
			for _, values := range inputs {
				if tt.heap {
					HeapSort(values, func(a, b int) int { return a - b })
				} else {
					Slice(values, func(a, b int) int { return a - b })
				}
				if !reflect.DeepEqual(want, values) {
					t.Errorf("want: %v, got: %v", want, values)
				}
			}
		})
	}
}

func makeInt(start, num int) []int {
	ret := make([]int, 0, num)
	for i := start; i < start+num; i++ {
		ret = append(ret, i)
	}
	return ret
}

func shuffle[T any](dataSet []T) [][]T {
	patterns := [][]T{}
	permutation(dataSet, func(ds []T) {
		ptn := []T{}
		ptn = append(ptn, ds...)
		patterns = append(patterns, ptn)
	})
	return patterns
}

func permutation[T any](a []T, f func([]T)) {
	perm(a, f, 0)
}

func perm[T any](a []T, f func([]T), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}
