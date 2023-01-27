package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/konoui/go-qsort"
	"github.com/konoui/go-qsort/cgo"
	"github.com/konoui/lipo/pkg/lipo/lmacho"
)

func main() {
	printCpus()
	inputs := []string{
		"i386", "x86_64", "x86_64h", "arm", "armv4t", "armv6", "armv7", "armv7f", "armv7s", "armv7k", "armv6m", "armv7m", "armv7em", "armv8m", "arm64", "arm64e", "arm64v8", "arm64_32",
	}
	arches := []lmacho.FatArchHeader{}
	cstruct := []string{}
	for _, i := range inputs {
		cpu, sub, ok := lmacho.ToCpu(i)
		if !ok {
			panic(i)
		}
		fa := lmacho.FatArchHeader{Cpu: cpu, SubCpu: sub, Align: 0}
		cstruct = append(cstruct, fmt.Sprintf(`{.name = "%s", .cputype = %d, .cpusubtype = %d, .align = %d}`, i, fa.Cpu, fa.SubCpu, fa.Align))
		arches = append(arches, fa)
	}
	fmt.Fprintln(os.Stderr, strings.Join(cstruct, ",\n"))

	qarches := make([]lmacho.FatArchHeader, len(arches))
	copy(qarches, arches)

	fmt.Printf("heapsort----------------------------------------------------------------------\n")
	qsort.HeapSort(arches, qsort.CmpFunc)
	for _, a := range arches {
		fmt.Println(lmacho.ToCpuString(a.Cpu, a.SubCpu))
	}

	fmt.Printf("qsort----------------------------------------------------------------------\n")
	sort(qarches, false)
	for _, a := range qarches {
		fmt.Println(lmacho.ToCpuString(a.Cpu, a.SubCpu))
	}
}

func sort(inputs []lmacho.FatArchHeader, cqsort bool) {
	if cqsort {
		cgo.Slice(inputs, func(a, b int) int {
			return qsort.CmpFunc(inputs[a], inputs[b])
		})
		return
	}
	qsort.Slice(inputs, qsort.CmpFunc)
}

func printCpus() {
	ret := []string{}
	for _, l := range lmacho.CpuNames() {
		ret = append(ret, `"`+l+`"`)
	}
	fmt.Fprintln(os.Stderr, strings.Join(ret, ",")+",")
}
