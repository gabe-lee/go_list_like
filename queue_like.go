package go_list_like

type QueueLike[T any, IDX Integer] interface {
	SliceLike[T, IDX]
	// Increment the start location (index/pointer/etc.) of this queue by
	// `n` positions. The new 'first' item in the queue should be the item
	// previously located at index `delta`
	IncrementStart(n IDX)
}

func DequeueOverwriteList[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], L ListLike[T, IDX2]](queue Q, count IDX1, destList L) {
	PeekOverwriteList(queue, count, destList)
	Discard(queue, count)
}
func TryDequeueOverwriteList[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], L ListLike[T, IDX2]](queue Q, count IDX1, destList L) (ok bool) {
	ok = TryPeekOverwriteList(queue, count, destList)
	if !ok {
		return
	}
	ok = TryDiscard(queue, count)
	return
}
func DequeueAppendToList[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], L ListLike[T, IDX2]](queue Q, count IDX1, destList L) {
	PeekAppendToList(queue, count, destList)
	Discard(queue, count)
}
func TryDequeueAppendToList[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], L ListLike[T, IDX2]](queue Q, count IDX1, destList L) (ok bool) {
	ok = TryPeekAppendToList(queue, count, destList)
	if !ok {
		return
	}
	ok = TryDiscard(queue, count)
	return
}
func DequeueToSlice[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], S SliceLike[T, IDX2]](queue Q, count IDX1, destSlice S) (nCopied IDX1) {
	nCopied = PeekToSlice(queue, count, destSlice)
	Discard(queue, count)
	return
}
func TryDequeueToSlice[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], S SliceLike[T, IDX2]](queue Q, count IDX1, destSlice S) (nCopied IDX1, ok bool) {
	nCopied, ok = TryPeekToSlice(queue, count, destSlice)
	if !ok {
		return
	}
	Discard(queue, count)
	return
}
func PeekToSlice[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], S SliceLike[T, IDX2]](queue Q, count IDX1, destSlice S) (nCopied IDX1) {
	nCopied, _, _, _, _ = CopyCount(queue, destSlice, count)
	return
}
func PeekOverwriteList[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], L ListLike[T, IDX2]](queue Q, count IDX1, destList L) {
	destList.Clear()
	destList.TryEnsureFreeSlots(IDX2(count))
	d1, d2 := destList.AppendSlotsAssumeCapacity(IDX2(count))
	CopyCountToRange(queue, destList, d1, d2, count)
}
func PeekAppendToList[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], L ListLike[T, IDX2]](queue Q, count IDX1, destList L) {
	destList.TryEnsureFreeSlots(IDX2(count))
	d1, d2 := destList.AppendSlotsAssumeCapacity(IDX2(count))
	CopyCountToRange(queue, destList, d1, d2, count)
}
func TryPeekToSlice[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], S SliceLike[T, IDX2]](queue Q, count IDX1, destSlice S) (nCopied IDX1, ok bool) {
	nCopied, _, _, _, _ = CopyCount(queue, destSlice, count)
	ok = nCopied == count
	return
}
func TryPeekOverwriteList[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], L ListLike[T, IDX2]](queue Q, count IDX1, destList L) (ok bool) {
	destList.Clear()
	ok = destList.TryEnsureFreeSlots(IDX2(count))
	if !ok {
		return
	}
	d1, d2 := destList.AppendSlotsAssumeCapacity(IDX2(count))
	_, fs, fd, _, _ := CopyCountToRange(queue, destList, d1, d2, count)
	ok = fs && fd
	return
}
func TryPeekAppendToList[T any, IDX1 Integer, IDX2 Integer, Q QueueLike[T, IDX1], L ListLike[T, IDX2]](queue Q, count IDX1, destList L) (ok bool) {
	ok = destList.TryEnsureFreeSlots(IDX2(count))
	if !ok {
		return
	}
	d1, d2 := destList.AppendSlotsAssumeCapacity(IDX2(count))
	_, fs, fd, _, _ := CopyCountToRange(queue, destList, d1, d2, count)
	ok = fs && fd
	return
}
func Discard[T any, IDX Integer, Q QueueLike[T, IDX]](queue Q, count IDX) {
	queue.IncrementStart(count)
}
func TryDiscard[T any, IDX Integer, Q QueueLike[T, IDX]](queue Q, count IDX) (ok bool) {
	ok = queue.Len() >= count
	if !ok {
		return
	}
	queue.IncrementStart(count)
	return
}
func DiscardAll[T any, IDX Integer, Q QueueLike[T, IDX]](queue Q) {
	Discard(queue, queue.Len())
}
