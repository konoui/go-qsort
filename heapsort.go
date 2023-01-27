/*-
 * Copyright 2023 konoui
 * Copyright (c) 1991, 1993
 *	The Regents of the University of California.  All rights reserved.
 *
 * This code is derived from software contributed to Berkeley by
 * Ronnie Kon at Mindcraft Inc., Kevin Lew and Elmer Yglesias.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 4. Neither the name of the University nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

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
			if j < num && cmp(list[p], list[p+1]) < 0 {
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
			if j < num && cmp(list[p], list[p+1]) < 0 {
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
			if j == 1 || cmp(k, list[t]) < 0 {
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
