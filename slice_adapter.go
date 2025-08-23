package go_list_like

import (
	"io"
)

type SliceAdapter[T any] struct {
	SlicePtr *[]T
}

func NewSliceAdapter[T any](slicePtr *[]T) SliceAdapter[T] {
	return SliceAdapter[T]{
		SlicePtr: slicePtr,
	}
}
func EmptySliceAdapter[T any](initCap int) SliceAdapter[T] {
	slice := make([]T, 0, initCap)
	return SliceAdapter[T]{
		SlicePtr: &slice,
	}
}
func (slice SliceAdapter[T]) WriteAt(src []T, off int64) (n int, err error) {
	if off > int64(slice.Len()) {
		return 0, io.EOF
	}
	end := int(off) + len(src)
	if end > slice.Len() {
		*slice.SlicePtr = (*slice.SlicePtr)[:off]
		*slice.SlicePtr = append(*slice.SlicePtr, src...)
		return len(src), nil
	}
	n = copy((*slice.SlicePtr)[off:], src)
	return n, nil
}

func (slice SliceAdapter[T]) ReadAt(dest []T, off int64) (n int, err error) {
	destLike := NewSliceAdapter(&dest)
	maxOff := min(int(off), slice.Len())
	subSrc := ((*slice.SlicePtr)[maxOff:])
	srcLike := NewSliceAdapter(&subSrc)
	n = min(destLike.Len(), srcLike.Len())
	Copy(destLike, 0, n, srcLike, 0, n)
	if n < destLike.Len() {
		err = io.EOF
	}
	return n, err
}

func (slice SliceAdapter[T]) Write(src []T) (n int, err error) {
	srcLike := NewSliceAdapter(&src)
	Append(slice, srcLike)
	return srcLike.Len(), nil
}

func (slice SliceAdapter[T]) Read(dest []T) (n int, err error) {
	destLike := NewSliceAdapter(&dest)
	n = min(destLike.Len(), slice.Len())
	Dequeue(slice, n, destLike)
	if n == 0 {
		err = io.EOF
	}
	return n, err
}

func (slice SliceAdapter[T]) Len() int {
	return len(*slice.SlicePtr)
}
func (slice SliceAdapter[T]) Cap() int {
	return cap(*slice.SlicePtr)
}
func (slice SliceAdapter[T]) Get(idx int) T {
	return (*slice.SlicePtr)[idx]
}
func (slice SliceAdapter[T]) GetPtr(idx int) *T {
	return &(*slice.SlicePtr)[idx]
}
func (slice SliceAdapter[T]) Set(idx int, val T) {
	(*slice.SlicePtr)[idx] = val
}
func (slice SliceAdapter[T]) ChangeLen(delta int) {
	*slice.SlicePtr = (*slice.SlicePtr)[:slice.Len()+delta]
}
func (slice SliceAdapter[T]) GrowCap(n int) {
	prevLen := len(*slice.SlicePtr)
	*slice.SlicePtr = (*slice.SlicePtr)[:cap(*slice.SlicePtr)]
	*slice.SlicePtr = append(*slice.SlicePtr, make([]T, n)...)
	*slice.SlicePtr = (*slice.SlicePtr)[:prevLen]
}
func (slice SliceAdapter[T]) OffsetStart(delta int) {
	*slice.SlicePtr = (*slice.SlicePtr)[delta:]
}

var _ MemQueueLike[byte] = SliceAdapter[byte]{}
var _ io.Reader = SliceAdapter[byte]{}
var _ io.Writer = SliceAdapter[byte]{}
var _ io.ReaderAt = SliceAdapter[byte]{}
var _ io.WriterAt = SliceAdapter[byte]{}
