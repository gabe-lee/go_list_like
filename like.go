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
	moveIdx := idx + count
	for moveIdx < listLen {
		oldptr := listLike.GetPtr(moveIdx)
		newptr := listLike.GetPtr(moveIdx - count)
		*newptr = *oldptr
		moveIdx += 1
	}
	listLike.AddLen(-count)
}

type SliceAdapter[T any] struct {
	SlicePtr *[]T
}

func New[T any](slicePtr *[]T) SliceAdapter[T] {
	return SliceAdapter[T]{
		SlicePtr: slicePtr,
	}
}

func (sa *SliceAdapter[T]) Len() int {
	return len(*sa.SlicePtr)
}
func (sa *SliceAdapter[T]) GetPtr(idx int) *T {
	return &(*sa.SlicePtr)[idx]
}
func (sa *SliceAdapter[T]) AddLen(delta int) {
	if delta < 0 {
		*sa.SlicePtr = (*sa.SlicePtr)[:sa.Len()+delta]
	} else if delta > 0 {
		*sa.SlicePtr = append(*sa.SlicePtr, make([]T, delta)...)
	}
}