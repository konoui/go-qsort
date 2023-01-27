#!/bin/bash
set -ie

sed -i.bak -E '/^\s+utilPrintf+/d' qsort.go
sed -i.bak -E '/^\s+utilDump+/d' qsort.go
! grep -i import qsort.go
! grep -i util qsort.go
! grep -i fmt qsort.go
! grep -i import heapsort.go
! grep -i util heapsort.go
! grep -i fmt heapsort.go
echo success!
