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
	// Get the value at the provided index
	Get(idx int) (val T)
	// Set the value at the provided index to the given value
	Set(idx int, val T)
	// Return the current number of values in the slice/list
	//
	// All values less than this length MUST be valid for Get() and Set()
	Len() int
}

func Len[T any](slice SliceLike[T]) int {
	return slice.Len()
}
func IdxInRange[T any](slice SliceLike[T], idx int) bool {
	return idx < slice.Len() && idx >= 0
}
func IdxUnderLen[T any](slice SliceLike[T], idx int) bool {
	return idx < slice.Len()
}
func IsEmpty[T any](slice SliceLike[T]) bool {
	return slice.Len() <= 0
}
func Get[T any](slice SliceLike[T], idx int) T {
	return slice.Get(idx)
}
func TryGet[T any](slice SliceLike[T], idx int) (val T, ok bool) {
	if !IdxInRange(slice, idx) {
		return val, false
	}
	return slice.Get(idx), true
}
func LastIdx[T any](slice SliceLike[T]) (lastIdx int) {
	return slice.Len() - 1
}
func GetLast[T any](slice SliceLike[T]) T {
	return slice.Get(LastIdx(slice))
}
func TryGetLast[T any](slice SliceLike[T]) (val T, ok bool) {
	if IsEmpty(slice) {
		return val, false
	}
	return slice.Get(LastIdx(slice)), true
}
func GetFirst[T any](slice SliceLike[T]) T {
	return slice.Get(0)
}
func TryGetFirst[T any](slice SliceLike[T]) (val T, ok bool) {
	if IsEmpty(slice) {
		return val, false
	}
	return slice.Get(0), true
}
func Set[T any](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, val)
}
func TrySet[T any](slice SliceLike[T], idx int, val T) (ok bool) {
	if !IdxInRange(slice, idx) {
		return false
	}
	slice.Set(idx, val)
	return true
}
func SetLast[T any](slice SliceLike[T], val T) {
	slice.Set(LastIdx(slice), val)
}
func TrySetLast[T any](slice SliceLike[T], val T) (ok bool) {
	if slice.Len() <= 0 {
		return false
	}
	slice.Set(LastIdx(slice), val)
	return true
}
func SetFirst[T any](slice SliceLike[T], val T) {
	slice.Set(0, val)
}
func TrySetFirst[T any](slice SliceLike[T], val T) (ok bool) {
	if slice.Len() <= 0 {
		return false
	}
	slice.Set(0, val)
	return true
}
func SetFrom[T any](dest SliceLike[T], destIdx int, source SliceLike[T], srcIdx int) {
	dest.Set(destIdx, source.Get(srcIdx))
}
func TrySetFrom[T any](dest SliceLike[T], destIdx int, source SliceLike[T], srcIdx int) (ok bool) {
	if !IdxInRange(dest, destIdx) || !IdxInRange(source, srcIdx) {
		return false
	}
	dest.Set(destIdx, source.Get(srcIdx))
	return true
}
func Swap[T any](slice SliceLike[T], idxA int, idxB int) {
	tmp := slice.Get(idxA)
	slice.Set(idxA, slice.Get(idxB))
	slice.Set(idxB, tmp)
}
func TrySwap[T any](slice SliceLike[T], idxA int, idxB int) (ok bool) {
	if !IdxInRange(slice, idxA) || !IdxInRange(slice, idxB) {
		return false
	}
	tmp := slice.Get(idxA)
	slice.Set(idxA, slice.Get(idxB))
	slice.Set(idxB, tmp)
	return true
}
func Move[T any](slice SliceLike[T], oldIdx int, newIdx int) {
	slice.Set(newIdx, slice.Get(oldIdx))
}
func TryMove[T any](slice SliceLike[T], oldIdx int, newIdx int) (ok bool) {
	if !IdxInRange(slice, oldIdx) || !IdxInRange(slice, newIdx) {
		return false
	}
	slice.Set(newIdx, slice.Get(oldIdx))
	return true
}
func Copy[T any](dest SliceLike[T], destStart, destLen int, source SliceLike[T], srcStart, srcLen int) (n int) {
	nn := min(destLen, srcLen)
	dIdx := destStart
	sIdx := srcStart
	n = 0
	for n < nn {
		dest.Set(dIdx, source.Get(sIdx))
		n += 1
		dIdx += 1
		sIdx += 1
	}
	return
}
func TryCopy[T any](dest SliceLike[T], destStart, destLen int, source SliceLike[T], srcStart, srcLen int) (n int, ok bool) {
	nn := min(destLen, srcLen)
	dIdx := destStart
	sIdx := srcStart
	n = 0
	for n < nn {
		if !IdxInRange(dest, dIdx) || !IdxInRange(source, sIdx) {
			return n, false
		}
		dest.Set(dIdx, source.Get(sIdx))
		n += 1
		dIdx += 1
		sIdx += 1
	}
	return n, true
}
func IsSorted[T any](slice SliceLike[T], greaterThan func(a T, b T) bool) bool {
	var i int = 0
	var ii int = 1
	for IdxInRange(slice, ii) {
		a := slice.Get(i)
		b := slice.Get(ii)
		if greaterThan(a, b) {
			return false
		}
		i = ii
		ii += 1
	}
	return true
}
func IsSortedImplicit[T cmp.Ordered](slice SliceLike[T]) bool {
	var i int = 0
	var ii int = 1
	for IdxInRange(slice, ii) {
		a := slice.Get(i)
		b := slice.Get(ii)
		if a > b {
			return false
		}
		i = ii
		ii += 1
	}
	return true
}
func Sort[T any](slice SliceLike[T], greaterThan func(a T, b T) bool) {
	if slice.Len() < 2 {
		return
	}
	var i int = 1
	var j int
	var jj int
	var elem T
	for IdxInRange(slice, i) {
		elem = Get(slice, i)
		j = i - 1
		jj = i
		for IdxInRange(slice, j) && greaterThan(slice.Get(j), elem) {
			Move(slice, j, jj)
			jj = j
			j -= 1
		}
		Set(slice, jj, elem)
		i += 1
	}
}

func SortImplicit[T cmp.Ordered](slice SliceLike[T]) {
	if slice.Len() < 2 {
		return
	}
	var i int = 1
	var j int
	var jj int
	var elem T
	for IdxInRange(slice, i) {
		elem = Get(slice, i)
		j = i - 1
		jj = i
		for IdxInRange(slice, j) && (slice.Get(j) > elem) {
			Move(slice, j, jj)
			jj = j
			j -= 1
		}
		Set(slice, jj, elem)
		i += 1
	}
}

func DoActionOnItemsUntilFalse[T any](slice SliceLike[T], action func(slice SliceLike[T], idx int, item T) (shouldContinue bool)) (stopIdx int) {
	idx := 0
	shouldContinue := true
	for IdxInRange(slice, idx) && shouldContinue {
		val := slice.Get(idx)
		shouldContinue = action(slice, idx, val)
		idx += 1
	}
	return idx
}
func DoActionOnItemsUntilFalseWithUserdata[T any, U any](slice SliceLike[T], action func(slice SliceLike[T], idx int, item T, userdata *U) (shouldContinue bool), userdata *U) (stopIdx int) {
	idx := 0
	shouldContinue := true
	for IdxInRange(slice, idx) && shouldContinue {
		val := slice.Get(idx)
		shouldContinue = action(slice, idx, val, userdata)
		idx += 1
	}
	return idx
}

func DoActionOnAllItems[T any](slice SliceLike[T], action func(slice SliceLike[T], idx int, item T)) {
	idx := 0
	for IdxInRange(slice, idx) {
		val := slice.Get(idx)
		action(slice, idx, val)
		idx += 1
	}
}
func DoActionOnAllItemsWithUserdata[T any, U any](slice SliceLike[T], action func(slice SliceLike[T], idx int, item T, userdata *U), userdata *U) {
	idx := 0
	for IdxInRange(slice, idx) {
		val := slice.Get(idx)
		action(slice, idx, val, userdata)
		idx += 1
	}
}
func FilterIndexes[T any, I Index](slice SliceLike[T], selectFunc func(slice SliceLike[T], idx I, item T) bool, outputList ListLike[I]) {
	Clear(outputList)
	idx := 0
	for IdxInRange(slice, 0) {
		val := slice.Get(idx)
		sel := selectFunc(slice, I(idx), val)
		if sel {
			AppendV(outputList, I(idx))
		}
		idx += 1
	}
}
func FilterIndexesWithUserdata[T any, I Index, U any](slice SliceLike[T], selectFunc func(slice SliceLike[T], idx I, item T, userdata *U) bool, outputList ListLike[I], userdata *U) {
	Clear(outputList)
	idx := 0
	for IdxInRange(slice, 0) {
		val := slice.Get(idx)
		sel := selectFunc(slice, I(idx), val, userdata)
		if sel {
			AppendV(outputList, I(idx))
		}
		idx += 1
	}
}
func MapValues[T any, TT any](slice SliceLike[T], mapFunc func(slice SliceLike[T], idx int, item T) TT, outputList ListLike[TT]) {
	Clear(outputList)
	GrowCapIfNeeded(outputList, slice.Len())
	idx := 0
	for IdxInRange(slice, idx) {
		val := slice.Get(idx)
		AppendV(outputList, mapFunc(slice, idx, val))
		idx += 1
	}
}
func MapValuesWithUserdata[T any, TT any, U any](slice SliceLike[T], mapFunc func(slice SliceLike[T], idx int, item T, userdata *U) TT, outputList ListLike[TT], userdata *U) {
	Clear(outputList)
	GrowCapIfNeeded(outputList, slice.Len())
	idx := 0
	for IdxInRange(slice, idx) {
		val := slice.Get(idx)
		AppendV(outputList, mapFunc(slice, idx, val, userdata))
		idx += 1
	}
}
func Accumulate[T any, TT any](slice SliceLike[T], initialAccumulation TT, accumulate func(slice SliceLike[T], idx int, item T, currentAccumulation TT) (newAccumulation TT)) (finalAccumulation TT) {
	idx := 0
	for IdxInRange(slice, idx) {
		val := slice.Get(idx)
		initialAccumulation = accumulate(slice, idx, val, initialAccumulation)
		idx += 1
	}
	return initialAccumulation
}
func AccumulateWithUserdata[T any, TT any, U any](slice SliceLike[T], initialAccumulation TT, accumulate func(slice SliceLike[T], idx int, item T, currentAccumulation TT, userdata *U) (newAccumulation TT), userdata *U) (finalAccumulation TT) {
	idx := 0
	for IdxInRange(slice, idx) {
		val := slice.Get(idx)
		initialAccumulation = accumulate(slice, idx, val, initialAccumulation, userdata)
		idx += 1
	}
	return initialAccumulation
}

// *************
// ListLike[T] *
// *************

type ListLike[T any] interface {
	SliceLike[T]
	// Increase the total number of elements the list can hold by `n` elements,
	// possibly reallocating/moving the data if needed
	GrowCap(n int)
	// Increase or decrease the length of the slice/list by `delta` elements,
	// assuming capacity already exists
	ChangeLen(delta int)
	// Return the total number of values the slice/list can hold
	Cap() int
}

func ChangeLen[T any](list ListLike[T], delta int) {
	list.ChangeLen(delta)
}
func GrowLen[T any](list ListLike[T], grow int) {
	list.ChangeLen(grow)
}
func ShrinkLen[T any](list ListLike[T], shrink int) {
	list.ChangeLen(-shrink)
}
func GrowCap[T any](list ListLike[T], n int) {
	list.GrowCap(n)
}
func GrowCapIfNeeded[T any](list ListLike[T], nMoreItems int) {
	space := list.Cap() - list.Len()
	if space >= nMoreItems {
		return
	}
	need := nMoreItems - space
	list.GrowCap(need)
}
func Cap[T any](list ListLike[T]) int {
	return list.Cap()
}
func Clear[T any](list ListLike[T]) {
	length := list.Len()
	list.ChangeLen(-length)
}
func AppendV[T any](list ListLike[T], vals ...T) {
	sVals := NewSliceAdapter(&vals)
	Append(list, sVals)
}
func Append[T any](list ListLike[T], vals SliceLike[T]) {
	start := list.Len()
	n := vals.Len()
	GrowCapIfNeeded(list, n)
	GrowLen(list, n)
	Copy(list, start, n, vals, 0, n)
}
func InsertV[T any](list ListLike[T], idx int, vals ...T) {
	sVals := NewSliceAdapter(&vals)
	Insert(list, idx, sVals)
}
func Insert[T any](list ListLike[T], idx int, vals SliceLike[T]) {
	removeIdx := list.Len() - 1
	moveLen := vals.Len()
	GrowCapIfNeeded(list, moveLen)
	GrowLen(list, moveLen)
	insertIdx := removeIdx + moveLen
	for removeIdx >= idx {
		Move(list, removeIdx, insertIdx)
		removeIdx -= 1
		insertIdx -= 1
	}
	Copy(list, idx, moveLen, vals, 0, moveLen)
}
func Delete[T any](list ListLike[T], idx int, count int) {
	listLen := list.Len()
	removeIdx := idx + count
	insertIdx := idx
	for removeIdx < listLen {
		Move(list, removeIdx, insertIdx)
		removeIdx += 1
		insertIdx += 1
	}
	ShrinkLen(list, count)
}
func DeleteSparse[T any, I Index](list ListLike[T], deleteIndexSlice SliceLike[I], sortDeleteIndexes bool) {
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
			Move(list, int(removeIdx), int(insertIdx))
			insertIdx += 1
			removeIdx += 1
		}
	}
	for removeIdx < I(list.Len()) {
		Move(list, int(removeIdx), int(insertIdx))
		insertIdx += 1
		removeIdx += 1
	}
	ShrinkLen(list, deleteIndexSlice.Len())
}
func Remove[T any](list ListLike[T], idx int, count int, outputList ListLike[T]) {
	Clear(outputList)
	GrowCapIfNeeded(outputList, count)
	GrowLen(outputList, count)
	Copy(outputList, 0, count, list, idx, count)
	Delete(list, idx, count)
}
func RemoveSparse[T any, I Index](list ListLike[T], removeIndexSlice SliceLike[I], sortRemoveIndexes bool, outputList ListLike[T]) {
	Clear(outputList)
	removeLen := removeIndexSlice.Len()
	GrowCapIfNeeded(outputList, removeLen)
	if sortRemoveIndexes {
		SortImplicit(removeIndexSlice)
	}
	insertIdx := Get(removeIndexSlice, 0)
	removeIdx := insertIdx
	deleteIdxIdx := 0
	deleteIdx := insertIdx
	for deleteIdxIdx < removeLen {
		if removeIdx == deleteIdx {
			AppendV(outputList, Get(list, int(removeIdx)))
			removeIdx += 1
			deleteIdxIdx += 1
			if deleteIdxIdx < removeLen {
				deleteIdx = Get(removeIndexSlice, deleteIdxIdx)
			}
		} else {
			Move(list, int(removeIdx), int(insertIdx))
			insertIdx += 1
			removeIdx += 1
		}
	}
	for removeIdx < I(list.Len()) {
		Move(list, int(removeIdx), int(insertIdx))
		insertIdx += 1
		removeIdx += 1
	}
	ShrinkLen(list, removeLen)
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
		dest.ChangeLen(-delta)
		Copy(dest, destStart, destLen-delta, source, srcStart, srcLen)
	} else {
		delta = srcLen - destLen
		moveUpIdx := dest.Len() - 1
		moveUpEnd := destStart + destLen - 1
		dest.ChangeLen(delta)
		for moveUpIdx > moveUpEnd {
			Move(dest, moveUpIdx, moveUpIdx+delta)
			moveUpIdx -= 1
		}
		Copy(dest, destStart, destLen+delta, source, srcStart, srcLen)
	}
	return delta
}

func Pop[T any](list ListLike[T]) T {
	last := LastIdx(list)
	val := Get(list, last)
	ShrinkLen(list, 1)
	return val
}

func TryPop[T any](list ListLike[T]) (val T, ok bool) {
	if IsEmpty(list) {
		return val, false
	}
	last := LastIdx(list)
	val = Get(list, last)
	ShrinkLen(list, 1)
	return val, true
}

// **************
// QueueLike[T] *
// **************

type QueueLike[T any] interface {
	ListLike[T]
	// Offset the start location (index/pointer/etc.) of this queue by
	// the given delta. The new 'first' item in the queue should be the item
	// previously located at `queue.GetPtr(0+delta)`.
	OffsetStart(delta int)
}

func Dequeue[T any](queueLike QueueLike[T], count int, outputList ListLike[T]) {
	Clear(outputList)
	GrowLen(outputList, count)
	Copy(outputList, 0, count, queueLike, 0, count)
	queueLike.OffsetStart(count)
}

// *****************
// MemSliceLike[T] *
// *****************

type MemSliceLike[T any] interface {
	SliceLike[T]
	// Get the a pointer to the value at the provided index
	GetPtr(idx int) *T
}

func GetPtr[T any](memSliceLike MemSliceLike[T], idx int) *T {
	return memSliceLike.GetPtr(idx)
}
func TryGetPtr[T any](memSliceLike MemSliceLike[T], idx int) (ptr *T, ok bool) {
	if !IdxInRange(memSliceLike, idx) {
		return ptr, false
	}
	return memSliceLike.GetPtr(idx), true
}

type MemListLike[T any] interface {
	MemSliceLike[T]
	ListLike[T]
}

type MemQueueLike[T any] interface {
	MemSliceLike[T]
	QueueLike[T]
}
