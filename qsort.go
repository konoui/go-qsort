package qsort

// https://github.com/apple-oss-distributions/Libc/blob/Libc-1534.40.2/stdlib/FreeBSD/qsort.c
func Slice[T any](x []T, cmpFunc func(i, j T) int) {
	qsort(x, 0, len(x), cmpFunc, depth(len(x)))
}

func qsort[T any](list []T, off int, num int, cmp func(T, T) int, depthLimit int) {
loop:
	swapCnt := 0
	if depthLimit <= 0 {
		heapSort(list, off, num, cmp)
		return
	}
	depthLimit--

	if num <= 7 {
		iSort(list, off, num, cmp, 0)
		return
	}

	pl := off
	pm := off + num/2
	pn := off + num - 1

	if num > 40 {
		d := (num / 8)
		pl = med3(list, pl, pl+d, pl+2*d, cmp)
		pm = med3(list, pm-d, pm, pm+d, cmp)
		pn = med3(list, pn-2*d, pn-d, pn, cmp)
	}

	pm = med3(list, pl, pm, pn, cmp)

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
					swapCnt = 1
					swap(list, pa, pb)
					pa++
				}
				pb++
			} else {
				break
			}
		}
		for pb <= pc {
			ret := cmp(list[pc], list[off])
			if ret >= 0 {
				if ret == 0 {
					swapCnt = 1
					swap(list, pc, pd)
					pd--
				}
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
			goto nevermind
		}
		return
	}

nevermind:
	d1 = pb - pa
	d2 := pd - pc
	if d1 <= d2 {
		if d1 > 1 {
			qsort(list, off, d1, cmp, depthLimit)
		}
		if d2 > 1 {
			off = pn - d2
			num = d2
			goto loop
		}
	} else {
		if d2 > 1 {
			qsort(list, pn-d2, d2, cmp, depthLimit)
		}
		if d1 > 1 {
			num = d1
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
