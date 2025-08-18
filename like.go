package go_slice_like

import "unsafe"

type SliceLike[T any] interface {
	GetPtr(idx int) *T
	Len() int
}

type ListLike[T any] interface {
	SliceLike[T]
	AddLen(delta int)
}

func Get[T any](sliceLike SliceLike[T], idx int) T {
	return *sliceLike.GetPtr(idx)
}
func GetPtr[T any](sliceLike SliceLike[T], idx int) *T {
	return sliceLike.GetPtr(idx)
}
func Set[T any](sliceLike SliceLike[T], idx int, val T) {
	*sliceLike.GetPtr(idx) = val
}
func Len[T any](sliceLike SliceLike[T]) int {
	return sliceLike.Len()
}

func Append[T any](listLike ListLike[T], vals ...T) {
	start := listLike.Len()
	listLike.AddLen(len(vals))
	ptr := listLike.GetPtr(start)
	slice := unsafe.Slice(ptr, len(vals))
	copy(slice, vals)
}
func Insert[T any](listLike ListLike[T], idx int, vals ...T) {
	moveIdx := listLike.Len() - 1
	moveLen := len(vals)
	listLike.AddLen(moveLen)
	for moveIdx >= idx {
		oldptr := listLike.GetPtr(moveIdx)
		newptr := listLike.GetPtr(moveIdx + moveLen)
		*newptr = *oldptr
		moveIdx -= 1
	}
	moveIdx += 1
	ptr := listLike.GetPtr(moveIdx)
	slice := unsafe.Slice(ptr, len(vals))
	copy(slice, vals)
}
func Delete[T any](listLike ListLike[T], idx int, count int) {
	listLen := listLike.Len()
	for idx < listLen {
		oldptr := listLike.GetPtr(idx + count)
		newptr := listLike.GetPtr(idx)
		*newptr = *oldptr
		idx += 1
	}
	listLike.AddLen(-count)
}

type SliceAdapter[T any] struct {
	Slice []T
}

func New[T any](slice []T) SliceAdapter[T] {
	return SliceAdapter[T]{
		Slice: slice,
	}
}

func (sa *SliceAdapter[T]) Len() int {
	return len(sa.Slice)
}
func (sa *SliceAdapter[T]) GetPtr(idx int) *T {
	return &sa.Slice[idx]
}
func (sa *SliceAdapter[T]) AddLen(delta int) {
	sa.Slice = append(sa.Slice, make([]T, delta)...)
}
