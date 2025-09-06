package go_list_like

// *******************
//  ListLike[T, IDX] *
// *******************

type ListLike[T any, IDX Integer] interface {
	SliceLike[T, IDX]
	// Ensure at least `n` empty capacity spaces exist to add new items without reallocating
	// the memory or perform any other expensive reorganization procedure
	//
	// If free space cannot be ensured and attempting to add `nMoreItems`
	// will definitely fail or cause undefined behaviour, `ok == false`
	TryEnsureFreeSlots(nMoreItems IDX) (ok bool)
	// Insert `n` new slots directly before existing index, shifting all existing items
	// at and after that index forward.
	//
	// Returns the first new slot and the last new slot, inclusive, but the first new slot might
	// not match the insert index, depending on the implementation behavior
	//
	// The implementation should assume that as long as `TryEnsureFreeSlots(count)` returns `true`,
	// calling this function with a valid insert idx should not fail
	InsertSlotsAssumeCapacity(idx IDX, count IDX) (firstNewSlot IDX, lastNewSlot IDX)
	// Append `n` new slots at the end of the list.
	//
	// Returns the first new slot and the last new slot, inclusive
	//
	// The implementation should assume that as long as `TryEnsureFreeSlots(count)` returns `true`,
	// calling this function with a valid insert idx should not fail
	AppendSlotsAssumeCapacity(count IDX) (firstNewSlot IDX, lastNewSlot IDX)
	// Remove all items between `firstRemoveIdx` and `lastRemovedIdx`, inclusive
	//
	// All items after `lastRemovedIdx` are shifted backward
	DeleteRange(firstRemovedIdx IDX, lastRemovedIdx IDX)
	// Reset list to an empty state. The list's capacity may or may not be retained.
	Clear()
	// Return the total number of values the slice/list can hold
	Cap() IDX
}

func TryEnsureFreeSlots[T any, IDX Integer, L ListLike[T, IDX]](list L, nMoreItems IDX) (ok bool) {
	ok = list.TryEnsureFreeSlots(nMoreItems)
	return
}
func EnsureFreeSlots[T any, IDX Integer, L ListLike[T, IDX]](list L, nMoreItems IDX) {
	list.TryEnsureFreeSlots(nMoreItems)
}
func Cap[T any, IDX Integer, L ListLike[T, IDX]](list L) IDX {
	return list.Cap()
}
func Clear[T any, IDX Integer, L ListLike[T, IDX]](list L) {
	list.Clear()
}
func AppendSlotsAssumeCapacity[T any, IDX Integer, L ListLike[T, IDX]](list L, count IDX) (firstNewSlot IDX, lastNewSlot IDX) {
	firstNewSlot, lastNewSlot = list.AppendSlotsAssumeCapacity(count)
	return
}
func AppendSlots[T any, IDX Integer, L ListLike[T, IDX]](list L, count IDX) (firstNewSlot IDX, lastNewSlot IDX) {
	list.TryEnsureFreeSlots(count)
	firstNewSlot, lastNewSlot = list.AppendSlotsAssumeCapacity(count)
	return
}
func TryAppendSlots[T any, IDX Integer, L ListLike[T, IDX]](list L, count IDX) (firstNewSlot IDX, lastNewSlot IDX, ok bool) {
	ok = list.TryEnsureFreeSlots(count)
	if !ok {
		return
	}
	firstNewSlot, lastNewSlot = list.AppendSlotsAssumeCapacity(count)
	return
}
func AppendVar[T any, IDX Integer, L ListLike[T, IDX]](list L, vals ...T) (firstAppendedIdx IDX, lastAppendedIdx IDX) {
	sVals := NewSliceAdapter(vals)
	firstAppendedIdx, lastAppendedIdx = Append(list, &sVals)
	return
}
func Append[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](list L, vals S) (firstAppendedIdx IDX1, lastAppendedIdx IDX1) {
	firstAppendedIdx, lastAppendedIdx = AppendSlots(list, IDX1(vals.Len()))
	CopyToRange(vals, list, firstAppendedIdx, lastAppendedIdx)
	return
}
func TryAppendVar[T any, IDX Integer, L ListLike[T, IDX]](list L, vals ...T) (firstAppendedIdx IDX, lastAppendedIdx IDX, ok bool) {
	sVals := NewSliceAdapter(vals)
	firstAppendedIdx, lastAppendedIdx, ok = TryAppend(list, &sVals)
	return
}
func TryAppend[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](list L, vals S) (firstAppendedIdx IDX1, lastAppendedIdx IDX1, ok bool) {
	firstAppendedIdx, lastAppendedIdx, ok = TryAppendSlots(list, IDX1(vals.Len()))
	if !ok {
		return
	}
	_, fs, fd, _, _ := CopyToRange(vals, list, firstAppendedIdx, lastAppendedIdx)
	ok = fs && fd
	return
}

func InsertSlotsAssumeCapacity[T any, IDX Integer, L ListLike[T, IDX]](list L, idx IDX, count IDX) (firstNewSlot IDX, lastNewSlot IDX) {
	firstNewSlot, lastNewSlot = list.InsertSlotsAssumeCapacity(idx, count)
	return
}

func InsertSlots[T any, IDX Integer, L ListLike[T, IDX]](list L, idx IDX, count IDX) (firstNewSlot IDX, lastNewSlot IDX) {
	list.TryEnsureFreeSlots(count)
	firstNewSlot, lastNewSlot = list.InsertSlotsAssumeCapacity(idx, count)
	return
}

func TryInsertSlots[T any, IDX Integer, L ListLike[T, IDX]](list L, idx IDX, count IDX) (firstNewSlot IDX, lastNewSlot IDX, ok bool) {
	ok = list.TryEnsureFreeSlots(count)
	if !ok {
		return
	}
	firstNewSlot, lastNewSlot = list.InsertSlotsAssumeCapacity(idx, count)
	return
}

func InsertVar[T any, IDX Integer, L ListLike[T, IDX]](list L, idx IDX, vals ...T) (firstInsertedIdx IDX, lastInsertedIdx IDX) {
	sVals := NewSliceAdapter(vals)
	firstInsertedIdx, lastInsertedIdx = Insert(list, idx, &sVals)
	return
}
func Insert[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](list L, idx IDX1, vals S) (firstInsertedIdx IDX1, lastInsertedIdx IDX1) {
	firstInsertedIdx, lastInsertedIdx = InsertSlots(list, idx, IDX1(vals.Len()))
	CopyToRange(vals, list, firstInsertedIdx, lastInsertedIdx)
	return
}

func TryInsertV[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](list L, idx IDX1, vals ...T) (firstInsertedIdx IDX1, lastInsertedIdx IDX1, ok bool) {
	sVals := NewSliceAdapter(vals)
	firstInsertedIdx, lastInsertedIdx, ok = TryInsert(list, idx, &sVals)
	return
}
func TryInsert[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](list L, idx IDX1, vals S) (firstInsertedIdx IDX1, lastInsertedIdx IDX1, ok bool) {
	firstInsertedIdx, lastInsertedIdx, ok = TryInsertSlots(list, idx, IDX1(vals.Len()))
	if !ok {
		return
	}
	_, fs, fd, _, _ := CopyToRange(vals, list, firstInsertedIdx, lastInsertedIdx)
	ok = fs && fd
	return
}
func DeleteRange[T any, IDX Integer, L ListLike[T, IDX]](list L, firstDeletedIdx IDX, lastDeletedIdx IDX) {
	list.DeleteRange(firstDeletedIdx, lastDeletedIdx)
}
func TryDeleteRange[T any, IDX Integer, L ListLike[T, IDX]](list L, firstDeletedIdx IDX, lastDeletedIdx IDX) (ok bool) {
	ok = list.RangeValid(firstDeletedIdx, lastDeletedIdx)
	if !ok {
		return
	}
	list.DeleteRange(firstDeletedIdx, lastDeletedIdx)
	return
}

func Delete[T any, IDX Integer, L ListLike[T, IDX]](list L, idx IDX, count IDX) {
	lastIdx := list.NthNextIdx(idx, count-1)
	DeleteRange(list, idx, lastIdx)
}
func TryDelete[T any, IDX Integer, L ListLike[T, IDX]](list L, idx IDX, count IDX) (ok bool) {
	ok = list.IdxValid(idx)
	if !ok {
		return
	}
	lastIdx := list.NthNextIdx(idx, count-1)
	ok = list.IdxValid(lastIdx)
	if !ok {
		return
	}
	ok = TryDeleteRange(list, idx, lastIdx)
	return
}
func RemoveRange[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](list L, firstRemovedIdx IDX1, lastRemovedIdx IDX1, dest S) {
	CopyFromRange(list, firstRemovedIdx, lastRemovedIdx, dest)
	list.DeleteRange(firstRemovedIdx, lastRemovedIdx)
}
func TryRemoveRange[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](list L, firstRemovedIdx IDX1, lastRemovedIdx IDX1, dest S) (ok bool) {
	_, fs, fd, _, _ := CopyFromRange(list, firstRemovedIdx, lastRemovedIdx, dest)
	ok = fs && fd
	if !ok {
		return
	}
	ok = TryDeleteRange(list, firstRemovedIdx, lastRemovedIdx)
	return
}
func Remove[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](list L, idx IDX1, count IDX1, dest S) {
	lastIdx := list.NthNextIdx(idx, count-1)
	RemoveRange(list, idx, lastIdx, dest)
}
func TryRemove[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](list L, idx IDX1, count IDX1, dest S) (ok bool) {
	ok = list.IdxValid(idx)
	if !ok {
		return
	}
	lastIdx := list.NthNextIdx(idx, count-1)
	ok = list.IdxValid(lastIdx)
	if !ok {
		return
	}
	ok = TryRemoveRange(list, idx, lastIdx, dest)
	return
}
func ReplaceRange[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](dest L, destStart IDX1, destEnd IDX1, source S) {
	_, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx := CopyToRange(source, dest, destStart, destEnd)
	if fullSourceCopied && !fullDestCopied {
		dest.DeleteRange(nextDestIdx, destEnd)
	} else if fullDestCopied && !fullSourceCopied {
		lastSrc := source.LastIdx()
		newLen := IDX1(source.LenBetween(nextSourceIdx, lastSrc))
		dest.TryEnsureFreeSlots(newLen)
		lastDestIdx := dest.LastIdx()
		var firstNew, lastNew IDX1
		if lastDestIdx == destEnd {
			firstNew, lastNew = dest.AppendSlotsAssumeCapacity(newLen)
		} else {
			destEnd = dest.NextIdx(destEnd)
			firstNew, lastNew = dest.InsertSlotsAssumeCapacity(destEnd, newLen)
		}
		CopyFromRangeToRange(source, nextSourceIdx, lastSrc, dest, firstNew, lastNew)
	}
}
func TryReplaceRange[T any, IDX1 Integer, IDX2 Integer, L ListLike[T, IDX1], S SliceLike[T, IDX2]](dest L, destStart IDX1, destEnd IDX1, source S) (ok bool) {
	ok = dest.RangeValid(destStart, destEnd)
	if !ok {
		return
	}
	_, fullSourceCopied, fullDestCopied, nextSourceIdx, nextDestIdx := CopyToRange(source, dest, destStart, destEnd)
	if fullSourceCopied && !fullDestCopied {
		ok = TryDeleteRange(dest, nextDestIdx, destEnd)
	} else if fullDestCopied && !fullSourceCopied {
		lastSrc := source.LastIdx()
		newLen := IDX1(source.LenBetween(nextSourceIdx, lastSrc))
		ok = dest.TryEnsureFreeSlots(newLen)
		if !ok {
			return
		}
		lastDestIdx := dest.LastIdx()
		var firstNew, lastNew IDX1
		if lastDestIdx == destEnd {
			firstNew, lastNew = dest.AppendSlotsAssumeCapacity(newLen)
		} else {
			destEnd = dest.NextIdx(destEnd)
			firstNew, lastNew = dest.InsertSlotsAssumeCapacity(destEnd, newLen)
		}
		_, fs, fd, _, _ := CopyFromRangeToRange(source, nextSourceIdx, lastSrc, dest, firstNew, lastNew)
		ok = fs && fd
	} else {
		ok = true
	}
	return
}

func Push[T any, IDX Integer, L ListLike[T, IDX]](list L, val T) {
	AppendVar(list, val)
}
func TryPush[T any, IDX Integer, L ListLike[T, IDX]](list L, val T) (ok bool) {
	_, _, ok = TryAppendVar(list, val)
	return
}
func PushGetIdx[T any, IDX Integer, L ListLike[T, IDX]](list L, val T) (idx IDX) {
	idx, _ = AppendVar(list, val)
	return
}
func TryPushGetIdx[T any, IDX Integer, L ListLike[T, IDX]](list L, val T) (idx IDX, ok bool) {
	idx, _, ok = TryAppendVar(list, val)
	return
}

func Pop[T any, IDX Integer, L ListLike[T, IDX]](list L) (val T) {
	last := LastIdx(list)
	val = Get(list, last)
	list.DeleteRange(last, last)
	return
}

func TryPop[T any, IDX Integer, L ListLike[T, IDX]](list L) (val T, ok bool) {
	last := LastIdx(list)
	ok = list.IdxValid(last)
	if !ok {
		return
	}
	val = Get(list, last)
	ok = TryDeleteRange(list, last, last)
	return
}
