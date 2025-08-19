package go_list_like

type SliceLike[T any] interface {
	GetPtr(idx int) *T
	Len() int
}

type ListLike[T any] interface {
	SliceLike[T]
	OffsetLen(delta int)
	Cap() int
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
func SetFrom[T any](dest SliceLike[T], destIdx int, source SliceLike[T], srcIdx int) {
	*dest.GetPtr(destIdx) = *source.GetPtr(srcIdx)
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

func Cap[T any](listLike ListLike[T]) int {
	return listLike.Cap()
}
func GrowLen[T any](listLike ListLike[T], grow int) {
	listLike.OffsetLen(grow)
}
func ShrinkLen[T any](listLike ListLike[T], shrink int) {
	listLike.OffsetLen(-shrink)
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
func Replace[T any](dest ListLike[T], destStart, destLen int, source SliceLike[T], srcStart, srcLen int) (delta int) {
	if destLen == srcLen {
		Copy(dest, destStart, destLen, source, srcStart, srcLen)
		return 0
	}
	if destLen > srcLen {
		delta = destLen - srcLen
		moveDownIdx := destStart + destLen
		for moveDownIdx < dest.Len() {
			Move(dest, moveDownIdx, moveDownIdx-delta)
			moveDownIdx += 1
		}
		dest.OffsetLen(-delta)
		Copy(dest, destStart, destLen-delta, source, srcStart, srcLen)
	} else {
		delta = srcLen - destLen
		moveUpIdx := dest.Len() - 1
		moveUpEnd := destStart + destLen - 1
		dest.OffsetLen(delta)
		for moveUpIdx > moveUpEnd {
			Move(dest, moveUpIdx, moveUpIdx+delta)
			moveUpIdx -= 1
		}
		Copy(dest, destStart, destLen+delta, source, srcStart, srcLen)
	}
	return delta
}

func Pop[T any](listLike ListLike[T]) T {
	ret := *listLike.GetPtr(listLike.Len() - 1)
	listLike.OffsetLen(-1)
	return ret
}

func GrowCapIfNeeded[T any](listLike ListLike[T], nMoreItems int) {
	space := listLike.Cap() - listLike.Len()
	if space >= nMoreItems {
		return
	}
	listLike.OffsetLen(nMoreItems)
	listLike.OffsetLen(-nMoreItems)
}

type SliceAdapter[T any] struct {
	SlicePtr *[]T
}

func New[T any](slicePtr *[]T) SliceAdapter[T] {
	return SliceAdapter[T]{
		SlicePtr: slicePtr,
	}
}

func (sa SliceAdapter[T]) Len() int {
	return len(*sa.SlicePtr)
}
func (sa SliceAdapter[T]) Cap() int {
	return cap(*sa.SlicePtr)
}
func (sa SliceAdapter[T]) GetPtr(idx int) *T {
	return &(*sa.SlicePtr)[idx]
}
func (sa SliceAdapter[T]) OffsetLen(delta int) {
	if delta < 0 {
		*sa.SlicePtr = (*sa.SlicePtr)[:sa.Len()+delta]
	} else if delta > 0 {
		*sa.SlicePtr = append(*sa.SlicePtr, make([]T, delta)...)
	}
}

var _ ListLike[byte] = SliceAdapter[byte]{}
