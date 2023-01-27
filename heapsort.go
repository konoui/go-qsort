package qsort

func HeapSort[T any](list []T, cmp func(T, T) int) {
	heapSort(list, 0, len(list), cmp)
}

// https://github.com/apple-oss-distributions/Libc/blob/Libc-1534.81.1/stdlib/FreeBSD/heapsort.c
func heapSort[T any](list []T, off, num int, cmp func(T, T) int) {
	if num <= 1 {
		return
	}

	base := off - 1
	for l := (num / 2) + 1; l > 0; l-- {
		i := l
		j := i * 2
		for ; double(i, &j) <= num; i = j {
			p := base + j
			//fmt.Printf("loop1: %d %d %d\n", l, i, j)
			if j < num && cmp(list[p], list[p+1]) < 0 {
				//fmt.Printf("loop1-swap: %d %d %d\n", l, i, j)
				p++
				j++
			}
			t := base + i
			if cmp(list[p], list[t]) <= 0 {
				break
			}
			swap(list, t, p)
		}
	}

	var k T
	for num > 1 {
		k = list[base+num] // _copy
		_copy(list, base+num, base+1)
		num--

		i := 1
		j := i * 2
		for ; double(i, &j) <= num; i = j {
			p := base + j
			//fmt.Printf("loop2: %d %d %d\n", num, i, j)
			if j < num && cmp(list[p], list[p+1]) < 0 {
				//fmt.Printf("loop2-swap: %d %d %d\n", num, i, j)
				p++
				j++
			}
			t := base + i
			_copy(list, t, p)
		}

		for {
			j := i
			i = j / 2
			p := base + j
			t := base + i
			//fmt.Printf("loop3: %d %d %d\n", num, i, j)
			if j == 1 || cmp(k, list[t]) < 0 {
				//fmt.Printf("loop3-swap: %d %d %d\n", num, i, j)
				list[p] = k // _copy
				break
			}
			_copy(list, p, t)
		}
	}
}

func _copy[T any](list []T, i, j int) {
	list[i] = list[j]
}

func double(i int, j *int) int {
	*j = i * 2
	return *j
}
