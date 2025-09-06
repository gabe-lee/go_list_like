package go_list_like

type SliceLike[T any, IDX Integer] interface {
	// Should return a constant boolean value describing whether
	// certain operations will peform better with linear operations
	// instead of binary-split operations
	//
	// If `LenBetween()`, `NthNextIdx()`, and `NthPrevIdx()` operate in O(N) time
	// instead of O(1), this should return true
	//
	// An example requiring `true` would be a linked list, where one must
	// traverse in linear time to find the true index n places after a given index,
	// or the number of items between two indexes
	//
	// Returning the correct value will allow some operations to use alternate,
	// more efficient algorithms
	PreferLinearOps() bool
	// Should return a constant boolean value describing whether consecutive indexes,
	// (eg. `0, 1, 2, 3, 4, 5`) are in their logical/proper order (not necessarily sorted)
	//
	// An example of this being true is the golang slice `[]T` and `SliceAdapter[T]`
	//
	// An example where this would be false is an implementation of a linked list
	//
	// This allows some algorithms to use more efficient paths
	ConsecutiveIndexesInOrder() bool
	// Should return a constant boolean value describing whether all indexes greater-than-or-equal-to
	// `0` AND less-than `slice.Len()` are valid
	//
	// An example of this being true is the golang slice `[]T` and `SliceAdapter[T]`
	//
	// An example where this would be false is an implementation of a linked list
	//
	// This allows some algorithms to use more efficient paths
	AllIndexesLessThanLenValid() bool
	// Returns whether the given index is valid for the slice
	IdxValid(idx IDX) bool
	// Returns whether the given index range is valid for the slice
	//
	// The following MUST be true:
	//   - `firstIdx` comes logically before OR is equal to `lastIdx`
	//   - all indexes including and between `firstIdx` and `lastIdx` are valid for the slice
	RangeValid(firstIdx IDX, lastIdx IDX) bool
	// Split an index range (roughly) in half, returning the index in the middle of the range
	//
	// Assumes `RangeValid(firstIdx, lastIdx) == true`, and if so,
	// the returned index MUST also be valid and MUST be between or equal to the first and/or last index
	//
	// The implementation should endeavor to return an index as close to the true middle index
	// as possible, but it is not required to as long as the returned index IS between or equal to
	// the first and/or last indexes. HOWEVER, some algorithms will have inconsitent performance
	// if the returned index is far from the true middle index
	SplitRange(firstIdx IDX, lastIdx IDX) (middleIdx IDX)
	// Get the value at the provided index
	Get(idx IDX) (val T)
	// Set the value at the provided index to the given value
	Set(idx IDX, val T)
	// Move the data located at `oldIdx` to `newIdx`, shifting all
	// values in between either up or down
	Move(oldIdx IDX, newIdx IDX)
	// Move the data from located between and including `firstIdx` and `lastIdx`,
	// to the position `newFirstIdx`, shifting the values at that location out of the way
	MoveRange(firstIdx IDX, lastIdx IDX, newFirstIdx IDX)
	// Return another SliceLike[T, I] that holds values in range [first, last] (inclusive)
	//
	// Analogous to slice[first:last+1]
	Slice(firstIdx IDX, lastIdx IDX) (slice SliceLike[T, IDX])
	// Return the first index in the slice.
	//
	// If the slice is empty, the index returned should
	// result in `IdxValid(idx) == false`
	FirstIdx() (idx IDX)
	// Return the last index in the slice.
	//
	// If the slice is empty, the index returned should
	// result in `IdxValid(idx) == false`
	LastIdx() (idx IDX)
	// Return the next index after the current index in the slice.
	//
	// If the given index is invalid or no next index exists,
	// the index returned should result in `IdxValid(idx) == false`
	NextIdx(thisIdx IDX) (nextIdx IDX)
	// Return the index `n` places after the current index in the slice.
	//
	// If the given index is invalid or no nth next index exists,
	// the index returned should result in `IdxValid(idx) == false`
	NthNextIdx(thisIdx IDX, n IDX) (nthNextIdx IDX)
	// Return the prev index before the current index in the slice.
	//
	// If the given index is invalid or no prev index exists,
	// the index returned should result in `IdxValid(idx) == false`
	PrevIdx(thisIdx IDX) (prevIdx IDX)
	// Return the index `n` places before the current index in the slice.
	//
	// If the given index is invalid or no nth previous index exists,
	// the index returned should result in `IdxValid(idx) == false`
	NthPrevIdx(thisIdx IDX, n IDX) (nthPrevIdx IDX)
	// Return the current number of values in the slice/list
	//
	// It is not guaranteed that all indexes less than `len` are valid for the slice
	Len() IDX
	// Return the number of items between (and including) `firstIdx` and `lastIdx`
	//
	// `slice.LenBetween(slice.FirstIdx(), slice.LastIdx())` MUST equal `slice.Len()`
	LenBetween(firstIdx IDX, lastIdx IDX) IDX
}

func IdxValid[T any, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX) bool {
	return slice.IdxValid(idx)
}
func AllIdxValid[T any, IDX Integer, S SliceLike[T, IDX]](slice S, idxs ...IDX) bool {
	for _, i := range idxs {
		if !slice.IdxValid(i) {
			return false
		}
	}
	return true
}
func RangeValid[T any, IDX Integer, S SliceLike[T, IDX]](slice S, firstIdx IDX, lastIdx IDX) bool {
	return slice.RangeValid(firstIdx, lastIdx)
}
func Len[T any, IDX Integer, S SliceLike[T, IDX]](slice S) IDX {
	return slice.Len()
}
func IsEmpty[T any, IDX Integer, S SliceLike[T, IDX]](slice S) bool {
	return slice.Len() <= 0
}
func TrySlice[T any, IDX Integer, S SliceLike[T, IDX]](slice S, firstIdx IDX, lastIdx IDX) (newSlice SliceLike[T, IDX], ok bool) {
	if !slice.RangeValid(firstIdx, lastIdx) {
		ok = false
		return
	}
	newSlice = slice.Slice(firstIdx, lastIdx)
	return
}
func Slice[T any, IDX Integer, S SliceLike[T, IDX]](slice S, firstIdx IDX, lastIdx IDX) (newSlice SliceLike[T, IDX]) {
	newSlice = slice.Slice(firstIdx, lastIdx)
	return
}
func Get[T any, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX) (val T) {
	val = slice.Get(idx)
	return
}
func TryGet[T any, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX) (val T, ok bool) {
	ok = slice.IdxValid(idx)
	if !ok {
		return
	}
	val = slice.Get(idx)
	return
}
func FirstIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S) (firstIdx IDX) {
	firstIdx = slice.FirstIdx()
	return
}
func TryFirstIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S) (firstIdx IDX, ok bool) {
	firstIdx = slice.FirstIdx()
	ok = slice.IdxValid(firstIdx)
	return
}
func LastIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S) (lastIdx IDX) {
	lastIdx = slice.LastIdx()
	return
}
func TryLastIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S) (lastIdx IDX, ok bool) {
	lastIdx = slice.LastIdx()
	ok = slice.IdxValid(lastIdx)
	return
}
func NextIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, thisIdx IDX) (nextIdx IDX) {
	nextIdx = slice.NextIdx(thisIdx)
	return
}
func TryNextIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, thisIdx IDX) (nextIdx IDX, ok bool) {
	ok = slice.IdxValid(thisIdx)
	if !ok {
		return
	}
	nextIdx = slice.NextIdx(thisIdx)
	ok = slice.IdxValid(nextIdx)
	return
}
func NthNextIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, thisIdx IDX, n IDX) (nextIdx IDX) {
	nextIdx = slice.NthNextIdx(thisIdx, n)
	return
}
func TryNthNextIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, thisIdx IDX, n IDX) (nextIdx IDX, ok bool) {
	ok = slice.IdxValid(thisIdx)
	if !ok {
		return
	}
	nextIdx = slice.NthNextIdx(thisIdx, n)
	ok = slice.IdxValid(nextIdx)
	return
}
func NthIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, n IDX) (nthIdx IDX) {
	if n == 0 {
		nthIdx = slice.FirstIdx()
		return
	}
	thisIdx := slice.FirstIdx()
	nthIdx = slice.NthNextIdx(thisIdx, n)
	return
}
func TryNthIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, n IDX) (nthIdx IDX, ok bool) {
	if n < 0 {
		return
	}
	if n == 0 {
		nthIdx = slice.FirstIdx()
		ok = slice.IdxValid(nthIdx)
		return
	}
	thisIdx := slice.FirstIdx()
	ok = slice.IdxValid(thisIdx)
	if !ok {
		return
	}
	nthIdx = slice.NthNextIdx(thisIdx, n)
	ok = slice.IdxValid(nthIdx)
	return
}
func PrevIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, thisIdx IDX) (prevIdx IDX) {
	prevIdx = slice.PrevIdx(thisIdx)
	return
}
func TryPrevIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, thisIdx IDX) (prevIdx IDX, ok bool) {
	ok = slice.IdxValid(thisIdx)
	if !ok {
		return
	}
	prevIdx = slice.PrevIdx(thisIdx)
	ok = slice.IdxValid(prevIdx)
	return
}
func NthPrevIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, thisIdx IDX, n IDX) (prevIdx IDX) {
	prevIdx = slice.NthPrevIdx(thisIdx, n)
	return
}
func TryNthPrevIdx[T any, IDX Integer, S SliceLike[T, IDX]](slice S, thisIdx IDX, n IDX) (prevIdx IDX, ok bool) {
	ok = slice.IdxValid(thisIdx)
	if !ok {
		return
	}
	prevIdx = slice.NthPrevIdx(thisIdx, n)
	ok = slice.IdxValid(prevIdx)
	return
}
func GetLast[T any, IDX Integer, S SliceLike[T, IDX]](slice S) (lastVal T) {
	lastIdx := slice.LastIdx()
	lastVal = slice.Get(lastIdx)
	return
}
func TryGetLast[T any, IDX Integer, S SliceLike[T, IDX]](slice S) (lastVal T, ok bool) {
	lastIdx := slice.LastIdx()
	ok = slice.IdxValid(lastIdx)
	if !ok {
		return
	}
	lastVal = slice.Get(lastIdx)
	return
}
func GetFirst[T any, IDX Integer, S SliceLike[T, IDX]](slice S) (firstVal T) {
	firstIdx := slice.FirstIdx()
	firstVal = slice.Get(firstIdx)
	return
}
func TryGetFirst[T any, IDX Integer, S SliceLike[T, IDX]](slice S) (firstVal T, ok bool) {
	firstIdx := slice.FirstIdx()
	ok = slice.IdxValid(firstIdx)
	if !ok {
		return
	}
	firstVal = slice.Get(firstIdx)
	return
}
func Set[T any, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, val)
}
func TrySet[T any, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (ok bool) {
	ok = slice.IdxValid(idx)
	if !ok {
		return
	}
	slice.Set(idx, val)
	return
}
func SetLast[T any, IDX Integer, S SliceLike[T, IDX]](slice S, val T) {
	lastIdx := slice.LastIdx()
	slice.Set(lastIdx, val)
}
func TrySetLast[T any, IDX Integer, S SliceLike[T, IDX]](slice S, val T) (ok bool) {
	lastIdx := slice.LastIdx()
	ok = slice.IdxValid(lastIdx)
	if !ok {
		return
	}
	slice.Set(lastIdx, val)
	return
}
func SetFirst[T any, IDX Integer, S SliceLike[T, IDX]](slice S, val T) {
	firstIdx := slice.FirstIdx()
	slice.Set(firstIdx, val)
}
func TrySetFirst[T any, IDX Integer, S SliceLike[T, IDX]](slice S, val T) (ok bool) {
	firstIdx := slice.FirstIdx()
	ok = slice.IdxValid(firstIdx)
	if !ok {
		return
	}
	slice.Set(firstIdx, val)
	return
}
func SetFrom[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S2, srcIdx IDX2, dest S1, destIdx IDX1) {
	val := dest.Get(destIdx)
	source.Set(srcIdx, val)
}
func TrySetFrom[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](dest S1, destIdx IDX1, source S2, srcIdx IDX2) (ok bool) {
	ok = dest.IdxValid(destIdx) && source.IdxValid(srcIdx)
	if !ok {
		return
	}
	val := dest.Get(destIdx)
	source.Set(srcIdx, val)
	return
}
func Move[T any, IDX Integer, S SliceLike[T, IDX]](slice S, oldIdx IDX, newIdx IDX) {
	slice.Move(oldIdx, newIdx)
}
func TryMove[T any, IDX Integer, S SliceLike[T, IDX]](slice S, oldIdx IDX, newIdx IDX) (ok bool) {
	ok = slice.IdxValid(oldIdx) && slice.IdxValid(newIdx)
	if !ok {
		return
	}
	slice.Move(oldIdx, newIdx)
	return
}
func MoveRange[T any, IDX Integer, S SliceLike[T, IDX]](slice S, firstIdx IDX, lastIdx IDX, newFirstIdx IDX) {
	slice.MoveRange(firstIdx, lastIdx, newFirstIdx)
}
func TryMoveRange[T any, IDX Integer, S SliceLike[T, IDX]](slice S, firstIdx IDX, lastIdx IDX, newFirstIdx IDX) (ok bool) {
	ok = slice.RangeValid(firstIdx, lastIdx) && slice.IdxValid(newFirstIdx) && slice.IdxValid(slice.NthNextIdx(newFirstIdx, slice.LenBetween(firstIdx, lastIdx)-1))
	if !ok {
		return
	}
	slice.MoveRange(firstIdx, lastIdx, newFirstIdx)
	return
}
func Swap[T any, IDX Integer, S SliceLike[T, IDX]](slice S, idxA IDX, idxB IDX) {
	oldB := slice.Get(idxB)
	oldA := slice.Get(idxA)
	slice.Set(idxA, oldB)
	slice.Set(idxB, oldA)
}
func TrySwap[T any, IDX Integer, S SliceLike[T, IDX]](slice S, idxA IDX, idxB IDX) (ok bool) {
	ok = slice.IdxValid(idxA) && slice.IdxValid(idxB)
	if !ok {
		return
	}
	Swap(slice, idxA, idxB)
	return
}
func Exchange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](sliceA S1, idxA IDX1, sliceB S2, idxB IDX2) {
	oldValA := sliceA.Get(idxA)
	oldValB := sliceB.Get(idxB)
	sliceA.Set(idxA, oldValB)
	sliceB.Set(idxB, oldValA)
}
func TryExchange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](sliceA S1, idxA IDX1, sliceB S2, idxB IDX2) (ok bool) {
	ok = sliceA.IdxValid(idxA) && sliceB.IdxValid(idxB)
	if !ok {
		return
	}
	Exchange(sliceA, idxA, sliceB, idxB)
	return
}
func Overwrite[T any, IDX Integer, S SliceLike[T, IDX]](slice S, oldIdx IDX, newIdx IDX) {
	val := slice.Get(oldIdx)
	slice.Set(newIdx, val)
}
func TryOverwrite[T any, IDX Integer, S SliceLike[T, IDX]](slice S, oldIdx IDX, newIdx IDX) (ok bool) {
	ok = slice.IdxValid(oldIdx) && slice.IdxValid(newIdx)
	if !ok {
		return
	}
	Overwrite(slice, oldIdx, newIdx)
	return
}
func Reverse[T any, IDX Integer, S SliceLike[T, IDX]](slice S) {
	left := slice.FirstIdx()
	right := slice.LastIdx()
	if left == right || !slice.IdxValid(left) || !slice.IdxValid(right) {
		return
	}
	for {
		Swap(slice, left, right)
		left = slice.NextIdx(left)
		if left == right {
			return
		}
		right = slice.PrevIdx(right)
		if left == right {
			return
		}
	}
}
func Fill[T any, IDX Integer, S SliceLike[T, IDX]](slice S, val T) {
	i := slice.FirstIdx()
	ok := slice.IdxValid(i)
	for ok {
		slice.Set(i, val)
		i = slice.NextIdx(i)
	}
}
func FillCount[T any, IDX Integer, S SliceLike[T, IDX]](slice S, count IDX, val T) (nFilled IDX, ok bool) {
	i := slice.FirstIdx()
	ok = slice.IdxValid(i)
	for ok {
		slice.Set(i, val)
		nFilled += 1
		i = slice.NextIdx(i)
	}
	ok = nFilled == count
	return
}
func FillCountFromPos[T any, IDX Integer, S SliceLike[T, IDX]](slice S, startIdx IDX, count IDX, val T) (nFilled IDX, ok bool) {
	i := startIdx
	ok = slice.IdxValid(i)
	for ok {
		slice.Set(i, val)
		nFilled += 1
		i = slice.NextIdx(i)
	}
	ok = nFilled == count
	return
}
func FillFromPos[T any, IDX Integer, S SliceLike[T, IDX]](slice S, startIdx IDX, val T) (nFilled IDX) {
	i := startIdx
	ok := slice.IdxValid(i)
	for ok {
		slice.Set(i, val)
		nFilled += 1
		i = slice.NextIdx(i)
	}
	return
}
func FillRange[T any, IDX Integer, S SliceLike[T, IDX]](slice S, firstIdx IDX, lastIdx IDX, val T) (nFilled IDX, fullRangeFilled bool, nextIdx IDX) {
	nextIdx = firstIdx
	ok := slice.IdxValid(nextIdx)
	for ok && !fullRangeFilled {
		slice.Set(nextIdx, val)
		nFilled += 1
		fullRangeFilled = nextIdx == lastIdx
		nextIdx = slice.NextIdx(nextIdx)
	}
	return
}
func Copy[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, dest S2) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d1 := dest.FirstIdx()
	d2 := dest.LastIdx()
	s1 := source.FirstIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, s1, s2, dest, d1, d2, 0, false)
	return
}
func CopyCount[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, dest S2, count IDX1) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d1 := dest.FirstIdx()
	d2 := dest.LastIdx()
	s1 := source.FirstIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, s1, s2, dest, d1, d2, count, true)
	return
}
func CopyFromPos[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, srcIdx IDX1, dest S2) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d1 := dest.FirstIdx()
	d2 := dest.LastIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, srcIdx, s2, dest, d1, d2, 0, false)
	return
}
func CopyCountFromPos[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, srcIdx IDX1, dest S2, count IDX1) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d1 := dest.FirstIdx()
	d2 := dest.LastIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, srcIdx, s2, dest, d1, d2, count, true)
	return
}
func CopyToPos[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, dest S2, destIdx IDX2) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d2 := dest.LastIdx()
	s1 := source.FirstIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, s1, s2, dest, destIdx, d2, 0, false)
	return
}
func CopyCountToPos[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, dest S2, destIdx IDX2, count IDX1) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d2 := dest.LastIdx()
	s1 := source.FirstIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, s1, s2, dest, destIdx, d2, count, true)
	return
}
func CopyFromPosToPos[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, srcIdx IDX1, dest S2, destIdx IDX2) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d2 := dest.LastIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, srcIdx, s2, dest, destIdx, d2, 0, false)
	return
}
func CopyCountFromPosToPos[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, srcIdx IDX1, dest S2, destIdx IDX2, count IDX1) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d2 := dest.LastIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, srcIdx, s2, dest, destIdx, d2, count, true)
	return
}
func CopyFromPosToRange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, srcIdx IDX1, dest S2, firstDestIdx IDX2, lastDestIdx IDX2) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, srcIdx, s2, dest, firstDestIdx, lastDestIdx, 0, false)
	return
}
func CopyCountFromPosToRange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, srcIdx IDX1, dest S2, firstDestIdx IDX2, lastDestIdx IDX2, count IDX1) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, srcIdx, s2, dest, firstDestIdx, lastDestIdx, count, true)
	return
}
func CopyFromRangeToPos[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, firstSrcIdx IDX1, lastSrcIdx IDX1, dest S2, destIdx IDX2) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d2 := dest.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, firstSrcIdx, lastSrcIdx, dest, destIdx, d2, 0, false)
	return
}
func CopyCountFromRangeToPos[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, firstSrcIdx IDX1, lastSrcIdx IDX1, dest S2, destIdx IDX2, count IDX1) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d2 := dest.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, firstSrcIdx, lastSrcIdx, dest, destIdx, d2, count, true)
	return
}
func CopyToRange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, dest S2, firstDestIdx IDX2, lastDestIdx IDX2) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	s1 := source.FirstIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, s1, s2, dest, firstDestIdx, lastDestIdx, 0, false)
	return
}
func CopyCountToRange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, dest S2, firstDestIdx IDX2, lastDestIdx IDX2, count IDX1) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	s1 := source.FirstIdx()
	s2 := source.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, s1, s2, dest, firstDestIdx, lastDestIdx, count, true)
	return
}
func CopyFromRange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, firstSourceIdx IDX1, lastSourceIdx IDX1, dest S2) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d1 := dest.FirstIdx()
	d2 := dest.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, firstSourceIdx, lastSourceIdx, dest, d1, d2, 0, false)
	return
}
func CopyCountFromRange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, firstSourceIdx IDX1, lastSourceIdx IDX1, dest S2, count IDX1) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	d1 := dest.FirstIdx()
	d2 := dest.LastIdx()
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, firstSourceIdx, lastSourceIdx, dest, d1, d2, count, true)
	return
}
func CopyFromRangeToRange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, firstSourceIdx IDX1, lastSourceIdx IDX1, dest S2, firstDestIdx IDX2, lastDestIdx IDX2) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, firstSourceIdx, lastSourceIdx, dest, firstDestIdx, lastDestIdx, 0, false)
	return
}
func CopyCountFromRangeToRange[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, firstSourceIdx IDX1, lastSourceIdx IDX1, dest S2, firstDestIdx IDX2, lastDestIdx IDX2, count IDX1) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	nCopied, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx = copyCountFromRangeToRange_internal(source, firstSourceIdx, lastSourceIdx, dest, firstDestIdx, lastDestIdx, count, true)
	return
}
func copyCountFromRangeToRange_internal[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX2]](source S1, firstSourceIdx IDX1, lastSourceIdx IDX1, dest S2, firstDestIdx IDX2, lastDestIdx IDX2, count IDX1, forceCount bool) (nCopied IDX1, fullSourceCopied bool, fullDestCopied bool, nextSourceIdx IDX1, nextDestIdx IDX2) {
	nextSourceIdx = firstSourceIdx
	nextDestIdx = firstDestIdx
	ok1 := source.IdxValid(firstSourceIdx)
	ok2 := dest.IdxValid(firstDestIdx)
	var val T
	for (!forceCount || count < nCopied) && ok1 && ok2 && !fullSourceCopied && !fullDestCopied {
		val = source.Get(nextSourceIdx)
		dest.Set(nextDestIdx, val)
		fullSourceCopied = nextSourceIdx == lastSourceIdx
		fullDestCopied = nextDestIdx == lastDestIdx
		nCopied += 1
		nextSourceIdx = source.NextIdx(nextSourceIdx)
		nextDestIdx = dest.NextIdx(nextDestIdx)
		ok1 = source.IdxValid(firstSourceIdx)
		ok2 = dest.IdxValid(firstDestIdx)
	}
	return
}
func Swizzle[T any, IDX1 Integer, IDX2 Integer, IDX3 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX3], SS SliceLike[S1, IDX1], IS SliceLike[IDX1, IDX2]](slices SS, selectors IS, dest S2) (nSwizzled IDX1, allSelectorsChosen bool, allDestFilled bool) {
	nSwizzled, allSelectorsChosen, allDestFilled = swizzleCount_internal(slices, selectors, dest, 0, false)
	return
}
func SwizzleCount[T any, IDX1 Integer, IDX2 Integer, IDX3 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX3], SS SliceLike[S1, IDX1], IS SliceLike[IDX1, IDX2]](slices SS, selectors IS, dest S2, count IDX1) (nSwizzled IDX1, allSelectorsChosen bool, allDestFilled bool) {
	nSwizzled, allSelectorsChosen, allDestFilled = swizzleCount_internal(slices, selectors, dest, count, true)
	return
}
func swizzleCount_internal[T any, IDX1 Integer, IDX2 Integer, IDX3 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[T, IDX3], SS SliceLike[S1, IDX1], IS SliceLike[IDX1, IDX2]](slices SS, selectors IS, dest S2, count IDX1, forceCount bool) (nSwizzled IDX1, allSelectorsChosen bool, allDestFilled bool) {
	idx := selectors.FirstIdx()
	moreSelectors := selectors.IdxValid(idx)
	var val T
	var sliceIdx IDX1
	var slice S1
	var selOk bool
	var sIdx IDX1
	var dIdx IDX3 = dest.FirstIdx()
	var moreDest = dest.IdxValid(dIdx)
	for (!forceCount || nSwizzled < count) && moreSelectors && moreDest {
		sliceIdx = selectors.Get(idx)
		selOk = slices.IdxValid(sliceIdx)
		if !selOk {
			break
		}
		slice = slices.Get(sliceIdx)
		sIdx = slice.FirstIdx()
		selOk = slice.IdxValid(sIdx)
		if !selOk {
			break
		}
		sIdx = slice.NthNextIdx(sIdx, nSwizzled)
		selOk = slice.IdxValid(sIdx)
		if !selOk {
			break
		}
		val = slice.Get(sIdx)
		dest.Set(dIdx, val)
		nSwizzled += 1
		idx = selectors.NextIdx(idx)
		moreSelectors = selectors.IdxValid(idx)
		dIdx = dest.NextIdx(dIdx)
		moreDest = dest.IdxValid(dIdx)
	}
	allDestFilled = !moreDest
	allSelectorsChosen = !moreSelectors
	return
}
func IsSorted[T any, IDX Integer, S SliceLike[T, IDX]](slice S, greaterThan func(a T, b T) (isGreaterThan bool)) (isSorted bool) {
	var i IDX
	var ii IDX
	var ok bool
	var a, b T
	i = slice.FirstIdx()
	ok = slice.IdxValid(i)
	if !ok {
		isSorted = true
		return
	}
	ii = slice.NextIdx(i)
	ok = slice.IdxValid(ii)
	if !ok {
		isSorted = true
		return
	}
	a = slice.Get(i)
	b = slice.Get(ii)
	for ok {
		if greaterThan(a, b) {
			isSorted = false
			return
		}
		i = ii
		ii = slice.NextIdx(ii)
		ok = slice.IdxValid(ii)
		if ok {
			a = b
			b = slice.Get(ii)
		}
	}
	return true
}
func IsSortedImplicit[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S) (isSorted bool) {
	isSorted = IsSorted(slice, GreaterThanImplicit)
	return
}
func InsertionSort[T any, IDX Integer, S SliceLike[T, IDX]](slice S, greaterThan func(a T, b T) (isGreaterThan bool)) {
	var ok bool
	var i, j, jj IDX
	var moveVal, testVal T
	i = slice.FirstIdx()
	ok = slice.IdxValid(i)
	if !ok {
		return
	}
	i = slice.NextIdx(i)
	ok = slice.IdxValid(i)
	if !ok {
		return
	}
	for ok {
		moveVal = slice.Get(i)
		j = slice.PrevIdx(i)
		ok = slice.IdxValid(j)
		if ok {
			jj = i
			testVal = slice.Get(j)
			for ok && greaterThan(testVal, moveVal) {
				Overwrite(slice, j, jj)
				jj = j
				j = slice.PrevIdx(j)
				ok = slice.IdxValid(j)
				if ok {
					testVal = slice.Get(j)
				}
			}
		}
		Set(slice, jj, moveVal)
		i = slice.NextIdx(i)
		ok = slice.IdxValid(i)
	}
}

func InsertionSortImplicit[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S) {
	InsertionSort(slice, GreaterThanImplicit)
}

func DoActionOnItemsUntilFalse[T any, IDX Integer, S SliceLike[T, IDX]](slice S, action func(slice S, idx IDX, item T) (shouldContinue bool)) (stopIdx IDX, actionCount IDX, stoppedAtEnd bool) {
	var ok bool
	var val T
	stopIdx = slice.FirstIdx()
	ok = slice.IdxValid(stopIdx)
	shouldContinue := true
	for ok && shouldContinue {
		val = slice.Get(stopIdx)
		shouldContinue = action(slice, stopIdx, val)
		actionCount += 1
		stopIdx = slice.NextIdx(stopIdx)
		ok = slice.IdxValid(stopIdx)
	}
	stoppedAtEnd = !ok
	return
}
func DoActionOnItemsUntilFalseWithUserdata[T any, IDX Integer, U any, S SliceLike[T, IDX]](slice S, action func(slice S, idx IDX, item T, userdata *U) (shouldContinue bool), userdata *U) (stopIdx IDX, actionCount IDX, stoppedAtEnd bool) {
	var ok bool
	var val T
	stopIdx = slice.FirstIdx()
	ok = slice.IdxValid(stopIdx)
	shouldContinue := true
	for ok && shouldContinue {
		val = slice.Get(stopIdx)
		shouldContinue = action(slice, stopIdx, val, userdata)
		actionCount += 1
		stopIdx = slice.NextIdx(stopIdx)
		ok = slice.IdxValid(stopIdx)
	}
	stoppedAtEnd = !ok
	return
}

func DoActionOnAllItems[T any, IDX Integer, S SliceLike[T, IDX]](slice S, action func(slice S, idx IDX, item T)) {
	var ok bool
	var idx IDX
	var val T
	idx = slice.FirstIdx()
	ok = slice.IdxValid(idx)
	for ok {
		val = slice.Get(idx)
		action(slice, idx, val)
		idx = slice.NextIdx(idx)
		ok = slice.IdxValid(idx)
	}
}
func DoActionOnAllItemsWithUserdata[T any, IDX Integer, U any, S SliceLike[T, IDX]](slice S, action func(slice S, idx IDX, item T, userdata *U), userdata *U) {
	var ok bool
	var idx IDX
	var val T
	idx = slice.FirstIdx()
	ok = slice.IdxValid(idx)
	for ok {
		val = slice.Get(idx)
		action(slice, idx, val, userdata)
		idx = slice.NextIdx(idx)
		ok = slice.IdxValid(idx)
	}
}

func DoActionOnItemsInRange[T any, IDX Integer, S SliceLike[T, IDX]](slice S, firstIdx, lastIdx IDX, action func(slice S, idx IDX, item T)) {
	var ok bool
	var idx IDX
	var val T
	idx = firstIdx
	ok = slice.IdxValid(idx)
	for ok {
		val = slice.Get(idx)
		action(slice, idx, val)
		ok = idx != lastIdx
		idx = slice.NextIdx(idx)
		ok = ok && slice.IdxValid(idx)
	}
}
func DoActionOnItemsInRangeWithUserdata[T any, IDX Integer, S SliceLike[T, IDX], U any](slice S, firstIdx, lastIdx IDX, action func(slice S, idx IDX, item T, userdata *U), userdata *U) {
	var ok bool
	var idx IDX
	var val T
	idx = firstIdx
	ok = slice.IdxValid(idx)
	for ok {
		val = slice.Get(idx)
		action(slice, idx, val, userdata)
		ok = idx != lastIdx
		idx = slice.NextIdx(idx)
		ok = ok && slice.IdxValid(idx)
	}
}

func DoActionOnItemsInRangeUntilFalse[T any, IDX Integer, S SliceLike[T, IDX]](slice S, firstIdx, lastIdx IDX, action func(slice S, idx IDX, item T) bool) {
	var ok bool
	var idx IDX
	var val T
	shouldContinue := true
	idx = firstIdx
	ok = slice.IdxValid(idx)
	for ok && shouldContinue {
		val = slice.Get(idx)
		shouldContinue = action(slice, idx, val)
		ok = idx != lastIdx
		idx = slice.NextIdx(idx)
		ok = ok && slice.IdxValid(idx)
	}
}
func DoActionOnItemsInRangeWithUntilFalseUserdata[T any, IDX Integer, S SliceLike[T, IDX], U any](slice S, firstIdx, lastIdx IDX, action func(slice S, idx IDX, item T, userdata *U) bool, userdata *U) {
	var ok bool
	var idx IDX
	var val T
	shouldContinue := true
	idx = firstIdx
	ok = slice.IdxValid(idx)
	for ok && shouldContinue {
		val = slice.Get(idx)
		shouldContinue = action(slice, idx, val, userdata)
		ok = idx != lastIdx
		idx = slice.NextIdx(idx)
		ok = ok && slice.IdxValid(idx)
	}
}

func FilterIndexes[T any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[IDX1, IDX2]](source S1, selectFunc func(slice S1, idx IDX1, item T) bool, dest S2) (allIindexesChecked bool) {
	var ok1, ok2, sel bool
	var idx1 IDX1
	var idx2 IDX2
	var val T
	idx1 = source.FirstIdx()
	idx2 = dest.FirstIdx()
	ok1 = source.IdxValid(idx1)
	ok2 = dest.IdxValid(idx2)
	for ok1 && ok2 {
		val = source.Get(idx1)
		sel = selectFunc(source, idx1, val)
		if sel {
			dest.Set(idx2, idx1)
			idx2 = dest.NextIdx(idx2)
			ok2 = dest.IdxValid(idx2)
		}
		idx1 = source.NextIdx(idx1)
		ok1 = source.IdxValid(idx1)
	}
	allIindexesChecked = !ok1
	return
}
func FilterIndexesWithUserdata[T any, IDX1 Integer, IDX2 Integer, U any, S1 SliceLike[T, IDX1], S2 SliceLike[IDX1, IDX2]](source S1, selectFunc func(slice S1, idx IDX1, item T, userdata *U) bool, dest S2, userdata *U) (allIindexesChecked bool) {
	var ok1, ok2, sel bool
	var idx1 IDX1
	var idx2 IDX2
	var val T
	idx1 = source.FirstIdx()
	idx2 = dest.FirstIdx()
	ok1 = source.IdxValid(idx1)
	ok2 = dest.IdxValid(idx2)
	for ok1 && ok2 {
		val = source.Get(idx1)
		sel = selectFunc(source, idx1, val, userdata)
		if sel {
			dest.Set(idx2, idx1)
			idx2 = dest.NextIdx(idx2)
			ok2 = dest.IdxValid(idx2)
		}
		idx1 = source.NextIdx(idx1)
		ok1 = source.IdxValid(idx1)
	}
	allIindexesChecked = !ok1
	return
}
func MapValues[T any, TT any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[TT, IDX2]](source S1, mapFunc func(slice S1, idx IDX1, item T) (newItem TT), dest S2) (allValuesMapped bool) {
	var ok1, ok2 bool
	var idx1 IDX1
	var idx2 IDX2
	var val1 T
	var val2 TT
	idx1 = source.FirstIdx()
	idx2 = dest.FirstIdx()
	ok1 = source.IdxValid(idx1)
	ok2 = dest.IdxValid(idx2)
	for ok1 && ok2 {
		val1 = source.Get(idx1)
		val2 = mapFunc(source, idx1, val1)
		dest.Set(idx2, val2)
		idx1 = source.NextIdx(idx1)
		idx2 = dest.NextIdx(idx2)
		ok1 = source.IdxValid(idx1)
		ok2 = dest.IdxValid(idx2)
	}
	allValuesMapped = !ok1
	return
}
func MapValuesWithUserdata[T any, TT any, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[TT, IDX2], U any](source S1, mapFunc func(slice S1, idx IDX1, item T, userdata *U) (newItem TT), dest S2, userdata *U) (allValuesMapped bool) {
	var ok1, ok2 bool
	var idx1 IDX1
	var idx2 IDX2
	var val1 T
	var val2 TT
	idx1 = source.FirstIdx()
	idx2 = dest.FirstIdx()
	ok1 = source.IdxValid(idx1)
	ok2 = dest.IdxValid(idx2)
	for ok1 && ok2 {
		val1 = source.Get(idx1)
		val2 = mapFunc(source, idx1, val1, userdata)
		dest.Set(idx2, val2)
		idx1 = source.NextIdx(idx1)
		idx2 = dest.NextIdx(idx2)
		ok1 = source.IdxValid(idx1)
		ok2 = dest.IdxValid(idx2)
	}
	allValuesMapped = !ok1
	return
}
func Accumulate[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, initialAccumulation TT, accumulate func(slice S, idx IDX, item T, currentAccumulation TT) (newAccumulation TT)) (finalAccumulation TT) {
	var ok bool
	var idx IDX
	var val T
	finalAccumulation = initialAccumulation
	idx = slice.FirstIdx()
	ok = slice.IdxValid(idx)
	for ok {
		val = slice.Get(idx)
		finalAccumulation = accumulate(slice, idx, val, finalAccumulation)
		idx = slice.NextIdx(idx)
		ok = slice.IdxValid(idx)
	}
	return
}
func AccumulateWithUserdata[T any, TT any, IDX Integer, U any, S SliceLike[T, IDX]](slice S, initialAccumulation TT, accumulate func(slice S, idx IDX, item T, currentAccumulation TT, userdata *U) (newAccumulation TT), userdata *U) (finalAccumulation TT) {
	var ok bool
	var idx IDX
	var val T
	finalAccumulation = initialAccumulation
	idx = slice.FirstIdx()
	ok = slice.IdxValid(idx)
	for ok {
		val = slice.Get(idx)
		finalAccumulation = accumulate(slice, idx, val, finalAccumulation, userdata)
		idx = slice.NextIdx(idx)
		ok = slice.IdxValid(idx)
	}
	return
}
