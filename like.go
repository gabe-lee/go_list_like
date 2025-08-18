package go_slice_like

type SliceLike[T any] interface {
	GetPtr(idx int) *T
	Len() int
}

type ListLike[T any] interface {
	SliceLike[T]
	OffsetLen(delta int)
}

func Len[T any](sliceLike SliceLike[T]) int {
	return sliceLike.Len()
}
func Get[T any](sliceLike SliceLike[T], idx int) T {
	return *sliceLike.GetPtr(idx)
}
func GetPtr[T any](sliceLike SliceLike[T], idx int) *T {
	return sliceLike.GetPtr(idx)
}
func GetLast[T any](sliceLike SliceLike[T]) T {
	return *sliceLike.GetPtr(sliceLike.Len() - 1)
}
func GetLastPtr[T any](sliceLike SliceLike[T]) *T {
	return sliceLike.GetPtr(sliceLike.Len() - 1)
}
func Set[T any](sliceLike SliceLike[T], idx int, val T) {
	*sliceLike.GetPtr(idx) = val
}
func SetLast[T any](sliceLike SliceLike[T], val T) {
	*sliceLike.GetPtr(sliceLike.Len() - 1) = val
}
func Swap[T any](sliceLike SliceLike[T], idxA int, idxB int) {
	tmp := Get(sliceLike, idxA)
	*GetPtr(sliceLike, idxA) = Get(sliceLike, idxB)
	*GetPtr(sliceLike, idxB) = tmp
}
func Move[T any](sliceLike SliceLike[T], oldIdx int, newIdx int) {
	*GetPtr(sliceLike, newIdx) = Get(sliceLike, oldIdx)
}
func Copy[T any](dest SliceLike[T], destStart, destLen int, source SliceLike[T], srcStart, srcLen int) (n int) {
	n = min(destLen, srcLen)
	d := destStart
	s := srcStart
	for n > 0 {
		dPtr := dest.GetPtr(d)
		sPtr := source.GetPtr(s)
		*dPtr = *sPtr
		n -= 1
		d += 1
		s += 1
	}
	return
}

func Append[T any](listLike ListLike[T], vals ...T) {
	end := listLike.Len()
	listLike.OffsetLen(len(vals))
	for i, v := range vals {
		ptr := listLike.GetPtr(end + i)
		*ptr = v
	}
}
func Insert[T any](listLike ListLike[T], idx int, vals ...T) {
	moveIdx := listLike.Len() - 1
	moveLen := len(vals)
	listLike.OffsetLen(moveLen)
	for moveIdx >= idx {
		oldptr := listLike.GetPtr(moveIdx)
		newptr := listLike.GetPtr(moveIdx + moveLen)
		*newptr = *oldptr
		moveIdx -= 1
	}
	moveIdx += 1
	for i, v := range vals {
		ptr := listLike.GetPtr(idx + i)
		*ptr = v
	}
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
	listLike.OffsetLen(-count)
}
func Remove[T any](listLike ListLike[T], idx int, count int) []T {
	ret := make([]T, count)
	for i := range ret {
		ret[i] = Get(SliceLike[T](listLike), i+idx)
	}
	listLen := listLike.Len()
	moveIdx := idx + count
	for moveIdx < listLen {
		oldptr := listLike.GetPtr(moveIdx)
		newptr := listLike.GetPtr(moveIdx - count)
		*newptr = *oldptr
		moveIdx += 1
	}
	listLike.OffsetLen(-count)
	return ret
}

func Pop[T any](listLike ListLike[T]) T {
	ret := *listLike.GetPtr(listLike.Len() - 1)
	listLike.OffsetLen(-1)
	return ret
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
