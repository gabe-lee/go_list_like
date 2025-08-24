package go_list_like

import (
	"io"
)

type SliceAdapterIndirect[T any] struct {
	SlicePtr *[]T
}

func NewSliceAdapterIndirect[T any](slicePtr *[]T) SliceAdapterIndirect[T] {
	return SliceAdapterIndirect[T]{
		SlicePtr: slicePtr,
	}
}
func EmptySliceAdapterIndirect[T any](initCap int) SliceAdapterIndirect[T] {
	slice := make([]T, 0, initCap)
	return SliceAdapterIndirect[T]{
		SlicePtr: &slice,
	}
}
func (slice SliceAdapterIndirect[T]) WriteAt(src []T, off int64) (n int, err error) {
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

func (slice SliceAdapterIndirect[T]) ReadAt(dest []T, off int64) (n int, err error) {
	destLike := NewSliceAdapterIndirect(&dest)
	maxOff := min(int(off), slice.Len())
	subSrc := ((*slice.SlicePtr)[maxOff:])
	srcLike := NewSliceAdapterIndirect(&subSrc)
	n = min(destLike.Len(), srcLike.Len())
	Copy(destLike, 0, n, srcLike, 0, n)
	if n < destLike.Len() {
		err = io.EOF
	}
	return n, err
}

func (slice SliceAdapterIndirect[T]) Write(src []T) (n int, err error) {
	srcLike := NewSliceAdapterIndirect(&src)
	Append(slice, srcLike)
	return srcLike.Len(), nil
}

func (slice SliceAdapterIndirect[T]) Read(dest []T) (n int, err error) {
	destLike := NewSliceAdapterIndirect(&dest)
	n = min(destLike.Len(), slice.Len())
	Dequeue(slice, n, destLike)
	if n == 0 {
		err = io.EOF
	}
	return n, err
}

func (slice SliceAdapterIndirect[T]) Len() int {
	return len(*slice.SlicePtr)
}
func (slice SliceAdapterIndirect[T]) Cap() int {
	return cap(*slice.SlicePtr)
}
func (slice SliceAdapterIndirect[T]) Get(idx int) T {
	return (*slice.SlicePtr)[idx]
}
func (slice SliceAdapterIndirect[T]) GetPtr(idx int) *T {
	return &(*slice.SlicePtr)[idx]
}
func (slice SliceAdapterIndirect[T]) Set(idx int, val T) {
	(*slice.SlicePtr)[idx] = val
}
func (slice SliceAdapterIndirect[T]) ChangeLen(delta int) {
	*slice.SlicePtr = (*slice.SlicePtr)[:slice.Len()+delta]
}
func (slice SliceAdapterIndirect[T]) GrowCap(n int) {
	prevLen := len(*slice.SlicePtr)
	*slice.SlicePtr = (*slice.SlicePtr)[:cap(*slice.SlicePtr)]
	*slice.SlicePtr = append(*slice.SlicePtr, make([]T, n)...)
	*slice.SlicePtr = (*slice.SlicePtr)[:prevLen]
}
func (slice SliceAdapterIndirect[T]) OffsetStart(delta int) {
	*slice.SlicePtr = (*slice.SlicePtr)[delta:]
}

func (slice SliceAdapterIndirect[T]) Slice(start int, end int) SliceLike[T] {
	newSlice := (*slice.SlicePtr)[start:end]
	return SliceAdapterIndirect[T]{
		SlicePtr: &newSlice,
	}
}

func (slice SliceAdapterIndirect[T]) GoSlice() []T {
	return *slice.SlicePtr
}

var _ MemQueueLike[byte] = SliceAdapterIndirect[byte]{}
var _ io.Reader = SliceAdapterIndirect[byte]{}
var _ io.Writer = SliceAdapterIndirect[byte]{}
var _ io.ReaderAt = SliceAdapterIndirect[byte]{}
var _ io.WriterAt = SliceAdapterIndirect[byte]{}

type SliceAdapter[T any] struct {
	data []T
}

func NewSliceAdapter[T any](slice []T) SliceAdapter[T] {
	return SliceAdapter[T]{
		data: slice,
	}
}
func EmptySliceAdapter[T any](initCap int) SliceAdapter[T] {
	slice := make([]T, 0, initCap)
	return SliceAdapter[T]{
		data: slice,
	}
}
func (slice *SliceAdapter[T]) WriteAt(src []T, off int64) (n int, err error) {
	if off > int64(slice.Len()) {
		return 0, io.EOF
	}
	end := int(off) + len(src)
	if end > slice.Len() {
		slice.data = slice.data[:off]
		slice.data = append(slice.data, src...)
		return len(src), nil
	}
	n = copy(slice.data[off:], src)
	return n, nil
}

func (slice *SliceAdapter[T]) ReadAt(dest []T, off int64) (n int, err error) {
	destLike := NewSliceAdapter(dest)
	maxOff := min(int(off), slice.Len())
	subSrc := (slice.data[maxOff:])
	srcLike := NewSliceAdapter(subSrc)
	n = min(destLike.Len(), srcLike.Len())
	Copy(&destLike, 0, n, &srcLike, 0, n)
	if n < destLike.Len() {
		err = io.EOF
	}
	return n, err
}

func (slice *SliceAdapter[T]) Write(src []T) (n int, err error) {
	srcLike := NewSliceAdapter(src)
	Append(slice, &srcLike)
	return srcLike.Len(), nil
}

func (slice *SliceAdapter[T]) Read(dest []T) (n int, err error) {
	destLike := NewSliceAdapter(dest)
	n = min(destLike.Len(), slice.Len())
	Dequeue(slice, n, &destLike)
	if n == 0 {
		err = io.EOF
	}
	return n, err
}

func (slice *SliceAdapter[T]) Len() int {
	return len(slice.data)
}
func (slice *SliceAdapter[T]) Cap() int {
	return cap(slice.data)
}
func (slice *SliceAdapter[T]) Get(idx int) T {
	return slice.data[idx]
}
func (slice *SliceAdapter[T]) GetPtr(idx int) *T {
	return &slice.data[idx]
}
func (slice *SliceAdapter[T]) Set(idx int, val T) {
	slice.data[idx] = val
}
func (slice *SliceAdapter[T]) ChangeLen(delta int) {
	slice.data = slice.data[:slice.Len()+delta]
}
func (slice *SliceAdapter[T]) GrowCap(n int) {
	prevLen := len(slice.data)
	slice.data = slice.data[:cap(slice.data)]
	slice.data = append(slice.data, make([]T, n)...)
	slice.data = slice.data[:prevLen]
}
func (slice *SliceAdapter[T]) OffsetStart(delta int) {
	slice.data = slice.data[delta:]
}

func (slice *SliceAdapter[T]) Slice(start int, end int) SliceLike[T] {
	newSlice := slice.data[start:end]
	return &SliceAdapter[T]{
		data: newSlice,
	}
}

func (slice *SliceAdapter[T]) GoSlice() []T {
	return slice.data
}

var _ MemQueueLike[byte] = (*SliceAdapter[byte])(nil)
var _ io.Reader = (*SliceAdapter[byte])(nil)
var _ io.Writer = (*SliceAdapter[byte])(nil)
var _ io.ReaderAt = (*SliceAdapter[byte])(nil)
var _ io.WriterAt = (*SliceAdapter[byte])(nil)
