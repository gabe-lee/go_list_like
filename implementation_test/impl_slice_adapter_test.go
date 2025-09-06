package implementation_test

import (
	"slices"
	"testing"

	LL "github.com/gabe-lee/go_list_like"
)

func newSliceAdapterPtr(t *testing.T, slice []byte) *LL.SliceAdapter[byte] {
	ss := slices.Clone(slice)
	sa := LL.NewSliceAdapter(ss)
	return &sa
}
func newSliceAdapterIndr(t *testing.T, slice []byte) LL.SliceAdapterIndirect[byte] {
	ss := slices.Clone(slice)
	return LL.NewSliceAdapterIndirect(&ss)
}

func Fuzz_SliceAdapter_(f *testing.F) {
	InitImplementationFuzz(f)
	PerformListImplementationFuzz(f, "SliceAdapter[byte]", newSliceAdapterPtr, func(t *testing.T, sa *LL.SliceAdapter[byte]) {})
}

func Fuzz_SliceAdapterIndirect_(f *testing.F) {
	InitImplementationFuzz(f)
	PerformListImplementationFuzz(f, "SliceAdapterIndirect[byte]", newSliceAdapterIndr, func(t *testing.T, sa LL.SliceAdapterIndirect[byte]) {})
}
