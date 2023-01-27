package qsort

// https://github.com/apple-oss-distributions/Libc/blob/Libc-1534.40.2/stdlib/FreeBSD/qsort.c
func Slice[T any](x []T, cmpFunc func(i, j T) int) {
	qsort(x, 0, len(x), cmpFunc, depth(len(x)))
}

func qsort[T any](list []T, off int, num int, cmp func(T, T) int, depthLimit int) {
	utilPrintf("myqsort is called %d\n", depthLimit)
loop:
	swapCnt := 0
	utilPrintf("num %d\n", num)
	if depthLimit <= 0 {
		utilPrintf("switch to myheapsort\n")
		heapSort(list, off, num, cmp)
		return
	}
	depthLimit--

	if num <= 7 {
		utilPrintf("switch to isort\n")
		iSort(list, off, num, cmp, 0)
		dump(list, off, num)
		return
	}

	pl := off
	pm := off + num/2
	pn := off + num - 1

	if num > 40 {
		panic(">40")
	}
	utilPrintf("pl %s\n", cpu(list[pl]))
	utilPrintf("pm %s\n", cpu(list[pm]))
	utilPrintf("pn %s\n", cpu(list[pn]))

	pm = med3(list, pl, pm, pn, cmp)
	utilPrintf("med3 pm %s\n", cpu(list[pm]))

	swap(list, off, pm)
	pa := off + 1
	pb := pa
	pc := off + num - 1
	pd := pc
	for {
		for pb <= pc {
			ret := cmp(list[pb], list[off])
			if ret <= 0 {
				if ret == 0 {
					utilPrintf("pb <= pc cmp_result=0\n")
					swapCnt = 1
					swap(list, pa, pb)
					pa++
				}
				utilPrintf("pb <= pc bp++\n")
				pb++
			} else {
				break
			}
		}
		for pb <= pc {
			ret := cmp(list[pc], list[off])
			if ret >= 0 {
				if ret == 0 {
					utilPrintf("pb <= pc cmp_result=0\n")
					swapCnt = 1
					swap(list, pc, pd)
					pd--
				}
				utilPrintf("pb <= pc pc--\n")
				pc--
			} else {
				break
			}
		}

		if pb > pc {
			break
		}

		swap(list, pb, pc)
		swapCnt = 1
		pb++
		pc--
	}

	pn = off + num
	d1 := min(pa-off, pb-pa)
	vecswap(list, off, pb-d1, d1)
	d1 = min(pd-pc, pn-pd-1)
	vecswap(list, pb, pn-d1, d1)

	if swapCnt == 0 {
		r := 1 + num/4
		if !iSort(list, off, num, cmp, r) {
			utilPrintf("goto nevermind;\n")
			goto nevermind
		}
		utilPrintf("return swap_cnt=0 isort\n")
		return
	}

nevermind:
	d1 = pb - pa
	d2 := pd - pc
	if d1 <= d2 {
		if d1 > 1 {
			utilPrintf("d1 > 1 do qsort\n")
			qsort(list, off, d1, cmp, depthLimit)
		}
		if d2 > 1 {
			off = pn - d2
			num = d2
			utilPrintf("d2 > 1 goto loop\n")
			goto loop
		}
	} else {
		if d2 > 1 {
			utilPrintf("d2 > 1 do qsort\n")
			qsort(list, pn-d2, d2, cmp, depthLimit)
		}
		if d1 > 1 {
			num = d1
			utilPrintf("d1 > 1 goto loop\n")
			goto loop
		}
	}
}

func min(a int, b int) int {
	if a >= b {
		return b
	}
	return a
}

func vecswap[T any](list []T, a int, b int, n int) {
	for i := 0; i < n; i++ {
		swap(list, a, b)
		a++
		b++
	}
}

func swap[T any](x []T, i, j int) {
	x[i], x[j] = x[j], x[i]
}

func depth(n int) int {
	return 2 * (fls(n) - 1)
}

func med3[T any](list []T, a, b, c int, f func(i, j T) int) int {
	if f(list[a], list[b]) < 0 {
		if f(list[b], list[c]) < 0 {
			return b
		} else {
			if f(list[a], list[c]) < 0 {
				return c
			}
			return a
		}
	} else {
		if f(list[b], list[c]) > 0 {
			return b
		} else {
			if f(list[a], list[c]) < 0 {
				return a
			} else {
				return c
			}
		}
	}
}

func fls(v int) int {
	if v == 0 {
		return 0
	}

	idx := 1
	tmp := v
	for {
		tmp >>= 1
		if tmp == 0 {
			return idx
		}
		idx++
	}
}

func iSort[T any](list []T, off int, len int, cmp func(a, b T) int, swapLimit int) bool {
	swapCnt := 0
	for i := off + 1; i < off+len; i++ {
		for j := i; j > off && cmp(list[j-1], list[j]) > 0; j-- {
			swap(list, j, j-1)
			swapCnt++
			if swapLimit > 0 && swapCnt > swapLimit {
				return false
			}
		}
	}

	return true
}
