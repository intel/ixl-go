<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# mem

```go
import "github.com/intel/ixl-go/util/mem"
```

Package mem declares functions for allocating memory that is aligned to specific byte boundaries.

## Index

- [func Alloc32Align[T any]() *T](<#func-alloc32align>)
- [func Alloc64Align[T any]() *T](<#func-alloc64align>)
- [func Alloc64ByteAligned(size uintptr) []byte](<#func-alloc64bytealigned>)


## func Alloc32Align

```go
func Alloc32Align[T any]() *T
```

Alloc32Align returns a pointer to a value of type T that is aligned to 32 bytes. The returned pointer may point to additional unused memory to ensure alignment.

## func Alloc64Align

```go
func Alloc64Align[T any]() *T
```

Alloc64Align returns a pointer to a value of type T that is aligned to 64 bytes. The returned pointer may point to additional unused memory to ensure alignment.

## func Alloc64ByteAligned

```go
func Alloc64ByteAligned(size uintptr) []byte
```

Alloc64ByteAligned returns a byte slice of the specified size that is aligned to 64 bytes. The returned slice may have additional unused capacity to ensure alignment.



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
