package cgo

/*
#include <stdlib.h>
typedef int (*cmp_qsort_func_t)(const void* a, const void* b);
extern int  cmp_qsort(void* a, void* b);
*/
import "C"

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

var qsortComparable struct {
	base unsafe.Pointer
	num  int
	size int
	less func(a, b int) int
	sync.Mutex
}

//export cmp_qsort
func cmp_qsort(a, b unsafe.Pointer) C.int {
	var (
		base = uintptr(qsortComparable.base)
		size = uintptr(qsortComparable.size)
	)

	i := int((uintptr(a) - base) / size)
	j := int((uintptr(b) - base) / size)
	return C.int(qsortComparable.less(i, j))
}

func Slice(x any, less func(a, b int) int) {
	rv := reflect.ValueOf(x)
	if rv.Kind() != reflect.Slice {
		panic(fmt.Sprintf("non-slice passed: %T", x))
	}
	if rv.Len() == 0 {
		return
	}

	qsortComparable.Lock()
	defer qsortComparable.Unlock()

	qsortComparable.base = unsafe.Pointer(rv.Index(0).Addr().Pointer())
	qsortComparable.num = rv.Len()
	qsortComparable.size = int(rv.Type().Elem().Size())
	qsortComparable.less = less
	defer func() {
		qsortComparable.base = nil
		qsortComparable.num = 0
		qsortComparable.size = 0
		qsortComparable.less = nil
	}()

	C.qsort(
		qsortComparable.base,
		C.size_t(qsortComparable.num),
		C.size_t(qsortComparable.size),
		C.cmp_qsort_func_t(C.cmp_qsort),
	)
}
