#!/bin/bash
set -ie

sed -i.bak -E '/^\s+utilPrintf+/d' qsort.go
! grep -i import qsort.go
! grep -i import heapsort.go

echo success!
