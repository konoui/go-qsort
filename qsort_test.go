package qsort_test

import (
	"io"
	"math/rand"
	"reflect"
	"testing"

	"github.com/konoui/go-qsort"
	"github.com/konoui/go-qsort/cgo"
	"github.com/konoui/lipo/pkg/lipo/lmacho"
)

func cmpFunc(i, j lmacho.FatArchHeader) int {
	if i.Cpu == j.Cpu {
		return int((i.SubCpu & ^lmacho.MaskSubCpuType)) - int((j.SubCpu & ^lmacho.MaskSubCpuType))
	}

	if i.Cpu == lmacho.CpuTypeArm64 {
		return 1
	}
	if j.Cpu == lmacho.CpuTypeArm64 {
		return -1
	}

	return int(i.Align) - int(j.Align)
}

func Test_qsort(t *testing.T) {
	tests := []struct {
		name string
		num  int
	}{
		{
			name: "len <= 7",
			num:  4,
		},
		{
			name: "7 < len < 40",
			num:  9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// turn off debug output
			qsort.Out = io.Discard

			fatArches := shuffle(makeFatArches(t, tt.num))
			t.Logf("%d test data\n", len(fatArches))
			for _, in := range fatArches {
				want := make([]lmacho.FatArchHeader, len(in))
				got := make([]lmacho.FatArchHeader, len(in))
				copy(want, in)
				copy(got, in)
				cgo.Slice(want, func(a, b int) int {
					return cmpFunc(want[a], want[b])
				})

				qsort.Slice(got, cmpFunc)
				if !reflect.DeepEqual(want, got) {
					t.Errorf("\nwant: %v\ngot: %v", names(want), names(got))
				}
			}
		})
	}
}

func Test_qsortOver40(t *testing.T) {
	t.Run("over40-1", func(t *testing.T) {
		in := makeFatArches(t, 200)
		want := make([]lmacho.FatArchHeader, len(in))
		got := make([]lmacho.FatArchHeader, len(in))
		copy(want, in)
		copy(got, in)
		cgo.Slice(want, func(a, b int) int {
			return cmpFunc(want[a], want[b])
		})

		qsort.Slice(got, cmpFunc)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("\nwant: %v\ngot: %v", names(want), names(got))
		}
	})
}

func newFatArch(t *testing.T, arch string, align uint32) lmacho.FatArchHeader {
	cpu, sub, ok := lmacho.ToCpu(arch)
	if !ok {
		t.Fatalf("found no cpu %s\n", arch)
	}
	return lmacho.FatArchHeader{
		Cpu:    cpu,
		SubCpu: sub,
		Align:  align,
	}
}

var cpuNames = func() []string {
	ret := []string{}
	for _, v := range lmacho.CpuNames() {
		// apple lipo does not support them
		if v == "armv8m" || v == "arm64_32" {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

func names(values []lmacho.FatArchHeader) []string {
	ret := make([]string, len(values))
	for i, v := range values {
		ret[i] = lmacho.ToCpuString(v.Cpu, v.SubCpu)
	}
	return ret
}

func makeFatArches(t *testing.T, num int) []lmacho.FatArchHeader {
	arches := cpuNames()
	if num > len(arches) {
		for {
			if len(arches) > num {
				break
			}
			arches = append(arches, cpuNames()...)
		}
	}

	ret := make([]lmacho.FatArchHeader, num)
	for i := 0; i < num; i++ {
		ret[i] = newFatArch(t, arches[i], uint32(pickRandom(3)))
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

func pickRandom(i int) int {
	return rand.Intn(i)
}
