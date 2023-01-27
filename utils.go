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

func dump[T any](list []T, off int, num int) {
	fmt.Fprintf(Out, "dump-----\n")
	for i := off; i < len(list); i++ {
		v := any(list[i]).(lmacho.FatArchHeader)
		fmt.Fprintf(Out, "%v\n", lmacho.ToCpuString(v.Cpu, v.SubCpu))
	}
	fmt.Fprintf(Out, "dump-----\n")
}

func cast[T any](v T) lmacho.FatArchHeader {
	return any(v).(lmacho.FatArchHeader)
}

func cpu[T any](v T) string {
	f := cast(v)
	return lmacho.ToCpuString(f.Cpu, f.SubCpu)
}
