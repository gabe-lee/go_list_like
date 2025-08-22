package go_list_like

type SliceAdapter[T any] struct {
	SlicePtr *[]T
}

func NewSliceAdapter[T any](slicePtr *[]T) SliceAdapter[T] {
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
func (sa SliceAdapter[T]) FirstIdx() (firstIdx int, hasFirst bool) {
	return 0, sa.Len() > 0
}
func (sa SliceAdapter[T]) NextIdx(thisIdx int) (nextIdx int, hasNext bool) {
	nextIdx = thisIdx + 1
	hasNext = nextIdx < sa.Len()
	return
}
func (sa SliceAdapter[T]) LastIdx() (lastIdx int, hasLast bool) {
	return sa.Len() - 1, sa.Len() > 0
}
func (sa SliceAdapter[T]) PrevIdx(thisIdx int) (prevIdx int, hasPrev bool) {
	prevIdx = thisIdx - 1
	hasPrev = prevIdx >= 0
	return
}

var _ ListLike[byte] = SliceAdapter[byte]{}
var _ FwdTraversable[byte] = SliceAdapter[byte]{}
var _ RevTraversable[byte] = SliceAdapter[byte]{}
