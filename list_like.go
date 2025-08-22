package go_list_like

import (
	"cmp"
)

type Index interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~int | ~int8 | ~int16 | ~int32 | ~int64
}

// **************
// SliceLike[T] *
// **************

type SliceLike[T any] interface {
	// Return a pointer to the value at index `idx`
	GetPtr(idx int) *T
	// Return the current number of values in the slice/list
	//
	// All indexes less than this value should be valid for `GetPtr(idx)`
	Len() int
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
func IsSorted[T any](sliceLike SliceLike[T], greaterThan func(a *T, b *T) bool) bool {
	var n int = sliceLike.Len()
	var i int = 1
	for i < n {
		a := sliceLike.GetPtr(i - 1)
		b := sliceLike.GetPtr(i)
		if greaterThan(a, b) {
			return false
		}
	}
	return true
}
func IsSortedImplicit[T cmp.Ordered](sliceLike SliceLike[T]) bool {
	var n int = sliceLike.Len()
	var i int = 1
	for i < n {
		a := sliceLike.GetPtr(i - 1)
		b := sliceLike.GetPtr(i)
		if *a > *b {
			return false
		}
	}
	return true
}
func Sort[T any](sliceLike SliceLike[T], greaterThan func(a *T, b T) bool) {
	var n int = sliceLike.Len()
	var i int = 1
	var j int
	var jj int
	var elem T
	for i < n {
		elem = Get(sliceLike, i)
		j = i - 1
		jj = i
		for j >= 0 && greaterThan(sliceLike.GetPtr(j), elem) {
			Move(sliceLike, j, jj)
			j -= 1
			jj -= 1
		}
		Set(sliceLike, jj, elem)
		i += 1
	}
}
func SortImplicit[T cmp.Ordered](sliceLike SliceLike[T]) {
	var n int = sliceLike.Len()
	var i int = 1
	var j int
	var jj int
	var elem T
	for i < n {
		elem = Get(sliceLike, i)
		j = i - 1
		jj = i
		for j >= 0 && Get(sliceLike, j) > elem {
			Move(sliceLike, j, jj)
			j -= 1
			jj -= 1
		}
		Set(sliceLike, jj, elem)
		i += 1
	}
}

// *************
// ListLike[T] *
// *************

type ListLike[T any] interface {
	SliceLike[T]
	// Increase or decrease the length of the slice/list by `delta` elements,
	// possibly reallocating/resizing the data if needed
	OffsetLen(delta int)
	// Return the total number of values the slice/list can hold
	Cap() int
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
func Clear[T any](listLike ListLike[T]) {
	length := listLike.Len()
	listLike.OffsetLen(-length)
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
func DeleteSparse[T any, I Index](listLike ListLike[T], deleteIndexSlice SliceLike[I], sortDeleteIndexes bool) {
	if sortDeleteIndexes {
		SortImplicit(deleteIndexSlice)
	}
	insertIdx := Get(deleteIndexSlice, 0)
	removeIdx := insertIdx
	deleteIdxIdx := 0
	deleteIdx := insertIdx
	for deleteIdxIdx < deleteIndexSlice.Len() {
		if removeIdx == deleteIdx {
			removeIdx += 1
			deleteIdxIdx += 1
			if deleteIdxIdx < deleteIndexSlice.Len() {
				deleteIdx = Get(deleteIndexSlice, deleteIdxIdx)
			}
		} else {
			Move(listLike, int(removeIdx), int(insertIdx))
			insertIdx += 1
			removeIdx += 1
		}
	}
	for removeIdx < I(listLike.Len()) {
		Move(listLike, int(removeIdx), int(insertIdx))
		insertIdx += 1
		removeIdx += 1
	}
	listLike.OffsetLen(-deleteIndexSlice.Len())
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

// *******************
// FwdTraversable[T] *
// *******************

type FwdTraversable[T any] interface {
	SliceLike[T]
	// Return the first index in the slice/list
	FirstIdx() (firstIdx int, hasFirst bool)
	// Return the next idx after this one, and whether the next idx is valid/exists
	NextIdx(thisIdx int) (nextIdx int, hasNext bool)
}

func DoActionOnItemsUntilFalse[T any](slice FwdTraversable[T], action func(slice FwdTraversable[T], idx int, item *T) (shouldContinue bool)) (prevIdx int, stopIdx int, stoppedAtEnd bool) {
	idx, hasNext := slice.FirstIdx()
	prevIdx = idx
	shouldContinue := true
	for hasNext && shouldContinue {
		ptr := slice.GetPtr(idx)
		shouldContinue = action(slice, idx, ptr)
		prevIdx = idx
		idx, hasNext = slice.NextIdx(idx)
	}
	return prevIdx, idx, !hasNext
}
func DoActionOnItemsUntilFalseWithUserdata[T any, U any](slice FwdTraversable[T], action func(slice FwdTraversable[T], idx int, item *T, userdata *U) (shouldContinue bool), userdata *U) (prevIdx int, stopIdx int, stoppedAtEnd bool) {
	idx, hasNext := slice.FirstIdx()
	prevIdx = idx
	shouldContinue := true
	for hasNext && shouldContinue {
		ptr := slice.GetPtr(idx)
		shouldContinue = action(slice, idx, ptr, userdata)
		prevIdx = idx
		idx, hasNext = slice.NextIdx(idx)
	}
	return prevIdx, idx, !hasNext
}

func DoActionOnAllItems[T any](slice FwdTraversable[T], action func(slice FwdTraversable[T], idx int, item *T)) {
	idx, hasNext := slice.FirstIdx()
	for hasNext {
		ptr := slice.GetPtr(idx)
		action(slice, idx, ptr)
		idx, hasNext = slice.NextIdx(idx)
	}
}
func DoActionOnAllItemsWithUserdata[T any, U any](slice FwdTraversable[T], action func(slice FwdTraversable[T], idx int, item *T, userdata *U), userdata *U) {
	idx, hasNext := slice.FirstIdx()
	for hasNext {
		ptr := slice.GetPtr(idx)
		action(slice, idx, ptr, userdata)
		idx, hasNext = slice.NextIdx(idx)
	}
}
func FilterIndexes[T any, I Index](slice FwdTraversable[T], selectFunc func(slice FwdTraversable[T], idx I, item *T) bool, outputList ListLike[I]) {
	Clear(outputList)
	idx, hasNext := slice.FirstIdx()
	for hasNext {
		ptr := slice.GetPtr(idx)
		sel := selectFunc(slice, I(idx), ptr)
		if sel {
			Append(outputList, I(idx))
		}
		idx, hasNext = slice.NextIdx(idx)
	}
}
func FilterIndexesWithUserdata[T any, I Index, U any](slice FwdTraversable[T], selectFunc func(slice FwdTraversable[T], idx I, item *T, userdata *U) bool, outputList ListLike[I], userdata *U) {
	Clear(outputList)
	idx, hasNext := slice.FirstIdx()
	for hasNext {
		ptr := slice.GetPtr(idx)
		sel := selectFunc(slice, I(idx), ptr, userdata)
		if sel {
			Append(outputList, I(idx))
		}
		idx, hasNext = slice.NextIdx(idx)
	}
}
func MapValues[T any, TT any](slice FwdTraversable[T], mapFunc func(slice FwdTraversable[T], idx int, item *T) TT, outputList ListLike[TT]) {
	Clear(outputList)
	GrowCapIfNeeded(outputList, slice.Len())
	idx, hasNext := slice.FirstIdx()
	for hasNext {
		ptr := slice.GetPtr(idx)
		Append(outputList, mapFunc(slice, idx, ptr))
		idx, hasNext = slice.NextIdx(idx)
	}
}
func MapValuesWithUserdata[T any, TT any, U any](slice FwdTraversable[T], mapFunc func(slice FwdTraversable[T], idx int, item *T, userdata *U) TT, outputList ListLike[TT], userdata *U) {
	Clear(outputList)
	GrowCapIfNeeded(outputList, slice.Len())
	idx, hasNext := slice.FirstIdx()
	for hasNext {
		ptr := slice.GetPtr(idx)
		Append(outputList, mapFunc(slice, idx, ptr, userdata))
		idx, hasNext = slice.NextIdx(idx)
	}
}
func Accumulate[T any, TT any](slice FwdTraversable[T], initialAccumulation TT, accumulate func(slice FwdTraversable[T], idx int, item *T, currentAccumulation TT) (newAccumulation TT)) (finalAccumulation TT) {
	idx, hasNext := slice.FirstIdx()
	for hasNext {
		ptr := slice.GetPtr(idx)
		initialAccumulation = accumulate(slice, idx, ptr, initialAccumulation)
		idx, hasNext = slice.NextIdx(idx)
	}
	return initialAccumulation
}
func AccumulateWithUserdata[T any, TT any, U any](slice FwdTraversable[T], initialAccumulation TT, accumulate func(slice FwdTraversable[T], idx int, item *T, currentAccumulation TT, userdata *U) (newAccumulation TT), userdata *U) (finalAccumulation TT) {
	idx, hasNext := slice.FirstIdx()
	for hasNext {
		ptr := slice.GetPtr(idx)
		initialAccumulation = accumulate(slice, idx, ptr, initialAccumulation, userdata)
		idx, hasNext = slice.NextIdx(idx)
	}
	return initialAccumulation
}

type FwdLinkedListLike[T any] interface {
	FwdTraversable[T]
	// Set the next idx after this one on the type located at this idx
	SetNextIdx(thisIdx int, nextIdx int)
}

type RevTraversable[T any] interface {
	SliceLike[T]
	// Return the last index in the slice/list
	LastIdx() (lastIdx int, hasLast bool)
	// Return the prev idx before this one, and whether the prev idx is valid/exists
	PrevIdx(thisIdx int) (prevIdx int, hasPrev bool)
}

type RevLinkedListLike[T any] interface {
	RevTraversable[T]
	// Set the prev idx before this one on the type located at this idx
	SetPrevIdx(thisIdx int, prevIdx int)
}
