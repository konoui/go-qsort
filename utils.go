package qsort

import (
	"fmt"
	"io"
	"os"

	"github.com/konoui/lipo/pkg/lipo/lmacho"
)

var (
	Out io.Writer = os.Stdout
)

func utilPrintf(format string, a ...any) {
	fmt.Fprintf(Out, format, a...)
}

func utilDump[T any](list []T, off int, num int) {
	fmt.Fprintf(Out, "dump-----\n")
	for i := off; i < len(list); i++ {
		v := any(list[i]).(lmacho.FatArchHeader)
		fmt.Fprintf(Out, "%v\n", lmacho.ToCpuString(v.Cpu, v.SubCpu))
	}
	fmt.Fprintf(Out, "dump-----\n")
}

func cast(v any) lmacho.FatArchHeader {
	return v.(lmacho.FatArchHeader)
}

func cpu(v any) string {
	f := cast(v)
	return lmacho.ToCpuString(f.Cpu, f.SubCpu)
}

func CmpFunc(i, j lmacho.FatArchHeader) int {
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
