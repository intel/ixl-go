#!/bin/bash
# Copyright (c) 2023, Intel Corporation.
# SPDX-License-Identifier: BSD-3-Clause

# Notice: c2goasm should call nasm to finish the job here. 
set -ex
ARGS="-masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -O3 -S"

clang $ARGS -msse -msse2 -o sse2.s ./codes.c 
c2goasm -a -f ./sse2.s ../sse2/codes_amd64.s

clang $ARGS -msse -msse2 -mavx2 -o avx2.s ./codes.c 
c2goasm -a -f ./avx2.s ../avx2/codes_amd64.s

clang $ARGS -msse -msse2 -mavx2 -mavx512f -o avx512.s ./codes.c 
c2goasm -a -f ./avx512.s ../avx512/codes_amd64.s

rm ./avx2.s ./avx512.s ./sse2.s