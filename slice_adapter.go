package go_list_like

import (
	"io"
	"slices"
)

type SliceAdapterIndirect[T any] struct {
	data *[]T
}

func NewSliceAdapterIndirect[T any](slicePtr *[]T) SliceAdapterIndirect[T] {
	return SliceAdapterIndirect[T]{
		data: slicePtr,
	}
}
func EmptySliceAdapterIndirect[T any](initCap int) SliceAdapterIndirect[T] {
	slice := make([]T, 0, initCap)
	return SliceAdapterIndirect[T]{
		data: &slice,
	}
}
func (slice SliceAdapterIndirect[T]) PreferLinearOps() bool {
	return false
}

func (slice SliceAdapterIndirect[T]) ConsecutiveIndexesInOrder() bool {
	return true
}
func (slice SliceAdapterIndirect[T]) AllIndexesLessThanLenValid() bool {
	return true
}

// Returns whether the given index is valid for the slice
func (slice SliceAdapterIndirect[T]) IdxValid(idx int) bool {
	return idx >= 0 && idx < len(*slice.data)
}

// Returns whether the given index range is valid for the slice
func (slice SliceAdapterIndirect[T]) RangeValid(firstIdx int, lastIdx int) bool {
	return firstIdx >= 0 && firstIdx <= lastIdx && lastIdx < len(*slice.data)
}

// Split an index range in half, returning the index in the middle of the range
//
// Assumes `RangeValid(firstIdx, lastIdx) == true`
func (slice SliceAdapterIndirect[T]) SplitRange(firstIdx int, lastIdx int) (middleIdx int) {
	return firstIdx + (((lastIdx + 1) - firstIdx) >> 1)
}

// Get the value at the provided index
func (slice SliceAdapterIndirect[T]) Get(idx int) (val T) {
	return (*slice.data)[idx]
}

// Set the value at the provided index to the given value
func (slice SliceAdapterIndirect[T]) Set(idx int, val T) {
	(*slice.data)[idx] = val
}

// Move the data located at `oldIdx` to `newIdx`, shifting all
// values in between either up or down
func (slice SliceAdapterIndirect[T]) Move(oldIdx int, newIdx int) {
	sa := NewSliceAdapter(*slice.data)
	sa.Move(oldIdx, newIdx)
}

// Remove all data contained in range `firstIdx` to `lastIdx` (inclusive),
// and re-insert it at the `newFirstIdx` position
func (slice SliceAdapterIndirect[T]) MoveRange(firstIdx int, lastIdx int, newFirstIdx int) {
	sa := NewSliceAdapter(*slice.data)
	sa.MoveRange(firstIdx, lastIdx, newFirstIdx)
}

// Return another SliceLike[T, I] that holds values in range [first, last]
//
// Analogous to slice[first:last+1]
func (slice SliceAdapterIndirect[T]) Slice(firstIdx int, lastIdx int) (newSlice SliceLike[T, int]) {
	s := *slice.data
	ns := s[firstIdx : lastIdx+1]
	return SliceAdapterIndirect[T]{
		data: &ns,
	}
}

// Return the first index in the slice.
//
// If the slice is empty, the index returned should
// result in `IdxValid(idx) == false`
func (slice SliceAdapterIndirect[T]) FirstIdx() (idx int) {
	return 0
}

// Return the last index in the slice.
//
// If the slice is empty, the index returned should
// result in `IdxValid(idx) == false`
func (slice SliceAdapterIndirect[T]) LastIdx() (idx int) {
	return len(*slice.data) - 1
}

// Return the next index after the current index in the slice.
//
// If the given index is invalid or no next index exists,
// the index returned should result in `IdxValid(idx) == false`
func (slice SliceAdapterIndirect[T]) NextIdx(thisIdx int) (nextIdx int) {
	return thisIdx + 1
}

// Return the index `n` places after the current index in the slice.
//
// If the given index is invalid or no nth next index exists,
// the index returned should result in `IdxValid(idx) == false`
func (slice SliceAdapterIndirect[T]) NthNextIdx(thisIdx int, n int) (nthNextIdx int) {
	return thisIdx + n
}

// Return the prev index before the current index in the slice.
//
// If the given index is invalid or no prev index exists,
// the index returned should result in `IdxValid(idx) == false`
func (slice SliceAdapterIndirect[T]) PrevIdx(thisIdx int) (prevIdx int) {
	return thisIdx - 1
}

// Return the index `n` places before the current index in the slice.
//
// If the given index is invalid or no nth previous index exists,
// the index returned should result in `IdxValid(idx) == false`
func (slice SliceAdapterIndirect[T]) NthPrevIdx(thisIdx int, n int) (nthPrevIdx int) {
	return thisIdx - n

}

// Return the current number of values in the slice/list
//
// It is not guaranteed that all indexes less than `len` are valid for the slice
func (slice SliceAdapterIndirect[T]) Len() int {
	return len(*slice.data)
}

// Return the number of items between (and including) `firstIdx` and `lastIdx`
func (slice SliceAdapterIndirect[T]) LenBetween(firstIdx int, lastIdx int) int {
	return min((lastIdx-firstIdx)+1, len(*slice.data))
}

// Ensure at least `n` empty capacity spaces exist to add new items without reallocating
// the memory or perform any other expensive reorganization procedure
//
// If free space cannot be ensured and attempting to add `nMoreItems`
// will definitely fail or cause undefined behaviour, `ok == false`
func (slice SliceAdapterIndirect[T]) TryEnsureFreeSlots(nMoreItems int) (ok bool) {
	*slice.data = slices.Grow(*slice.data, nMoreItems)
	ok = true
	return
}

// Insert `n` new slots at existing index, shifting all existing items
// after them forward. Returns the first new slot and the last new slot, inclusive.
//
// The implementation should assume no checks need need to be made to ensure free space exists,
// and calling this should not perform any reallocation
//
// May or may not work if attempting to insert items at an invalid indexIDX
func (slice SliceAdapterIndirect[T]) InsertSlotsAssumeCapacity(idx int, count int) (firstNewSlot int, lastNewSlot int) {
	sa := NewSliceAdapter(*slice.data)
	firstNewSlot, lastNewSlot = sa.InsertSlotsAssumeCapacity(idx, count)
	*slice.data = sa.data
	return
}

// Append `n` new slots at the end of the list.
//
// Returns the first new slot and the last new slot, inclusive.
//
// The implementation should assume no checks need need to be made to ensure free space exists,
// and calling this should not perform any reallocation
func (slice SliceAdapterIndirect[T]) AppendSlotsAssumeCapacity(count int) (firstNewSlot int, lastNewSlot int) {
	sa := NewSliceAdapter(*slice.data)
	firstNewSlot, lastNewSlot = sa.AppendSlotsAssumeCapacity(count)
	*slice.data = sa.data
	return
}

// Remove all items between `firstRemoveIdx` and `lastRemovedIdx`, inclusive
//
// All items after `lastRemovedIdx` are shifted backward
func (slice SliceAdapterIndirect[T]) DeleteRange(firstRemovedIdx int, lastRemovedIdx int) {
	sa := NewSliceAdapter(*slice.data)
	sa.DeleteRange(firstRemovedIdx, lastRemovedIdx)
	*slice.data = sa.data
}

// Reset list to an empty state. The list's capacity may or may not be retained.
func (slice SliceAdapterIndirect[T]) Clear() {
	*slice.data = (*slice.data)[:0]
}

// Return the total number of values the slice/list can hold
func (slice SliceAdapterIndirect[T]) Cap() int {
	return cap(*slice.data)
}

// GetPtr implements MemQueueLike.
func (slice SliceAdapterIndirect[T]) GetPtr(idx int) *T {
	return &(*slice.data)[idx]
}

// OffsetStart implements MemQueueLike.
func (slice SliceAdapterIndirect[T]) IncrementStart(n int) {
	*slice.data = (*slice.data)[n:]
}

func (slice SliceAdapterIndirect[T]) WriteAt(src []T, off int64) (n int, err error) {
	if off < 0 || off > int64(len(*slice.data)) {
		return 0, io.EOF
	}
	end := int(off) + len(src)
	if end >= slice.Len() {
		(*slice.data) = (*slice.data)[:off]
		(*slice.data) = append((*slice.data), src...)
		n = len(src)
	} else {
		n = copy((*slice.data)[off:], src)
	}
	return n, nil
}

func (slice SliceAdapterIndirect[T]) ReadAt(dest []T, off int64) (n int, err error) {
	n = copy(dest, (*slice.data)[off:])
	if n < len(dest) {
		err = io.EOF
	}
	return n, err
}

func (slice SliceAdapterIndirect[T]) Write(src []T) (n int, err error) {
	(*slice.data) = append((*slice.data), src...)
	return len(src), nil
}

func (slice SliceAdapterIndirect[T]) Read(dest []T) (n int, err error) {
	n = copy(dest, (*slice.data))
	if n == 0 {
		err = io.EOF
	}
	return n, err
}

func (slice SliceAdapterIndirect[T]) GoSlice() []T {
	return *slice.data
}

var _ MemQueueLike[byte, int] = SliceAdapterIndirect[byte]{}
var _ MemListLike[byte, int] = SliceAdapterIndirect[byte]{}
var _ GoSliceLike[byte] = SliceAdapterIndirect[byte]{}
var _ io.Reader = SliceAdapterIndirect[byte]{}
var _ io.Writer = SliceAdapterIndirect[byte]{}
var _ io.ReaderAt = SliceAdapterIndirect[byte]{}
var _ io.WriterAt = SliceAdapterIndirect[byte]{}

type SliceAdapter[T any] struct {
	data []T
}

func NewSliceAdapter[T any](slice []T) SliceAdapter[T] {
	return SliceAdapter[T]{
		data: slice,
	}
}
func EmptySliceAdapter[T any](initCap int) SliceAdapter[T] {
	slice := make([]T, 0, initCap)
	return SliceAdapter[T]{
		data: slice,
	}
}

func (slice *SliceAdapter[T]) PreferLinearOps() bool {
	return false
}

func (slice *SliceAdapter[T]) ConsecutiveIndexesInOrder() bool {
	return true
}
func (slice *SliceAdapter[T]) AllIndexesLessThanLenValid() bool {
	return true
}

// Returns whether the given index is valid for the slice
func (slice *SliceAdapter[T]) IdxValid(idx int) bool {
	return idx >= 0 && idx < len(slice.data)
}

// Returns whether the given index range is valid for the slice
func (slice *SliceAdapter[T]) RangeValid(firstIdx int, lastIdx int) bool {
	return firstIdx >= 0 && firstIdx <= lastIdx && lastIdx < len(slice.data)
}

// Split an index range in half, returning the index in the middle of the range
//
// Assumes `RangeValid(firstIdx, lastIdx) == true`
func (slice *SliceAdapter[T]) SplitRange(firstIdx int, lastIdx int) (middleIdx int) {
	return firstIdx + (((lastIdx + 1) - firstIdx) >> 1)
}

// Get the value at the provided index
func (slice *SliceAdapter[T]) Get(idx int) (val T) {
	return slice.data[idx]
}

// Set the value at the provided index to the given value
func (slice *SliceAdapter[T]) Set(idx int, val T) {
	slice.data[idx] = val
}

// Move the data located at `oldIdx` to `newIdx`, shifting all
// values in between either up or down
func (slice *SliceAdapter[T]) Move(oldIdx int, newIdx int) {
	val := slice.Get(oldIdx)
	if newIdx < oldIdx {
		prevIdx := oldIdx - 1
		for oldIdx > newIdx {
			slice.data[oldIdx] = slice.data[prevIdx]
			oldIdx = prevIdx
			prevIdx -= 1
		}
	} else {
		nextIdx := oldIdx + 1
		for oldIdx < newIdx {
			slice.data[oldIdx] = slice.data[nextIdx]
			oldIdx = nextIdx
			nextIdx += 1
		}
	}
	slice.data[newIdx] = val
}

// Remove all data contained in range `firstIdx` to `lastIdx` (inclusive),
// and re-insert it at the `newFirstIdx` position
func (slice *SliceAdapter[T]) MoveRange(firstIdx int, lastIdx int, newFirstIdx int) {
	lenA := (lastIdx - firstIdx) + 1
	sliceA := slice.data[firstIdx : firstIdx+lenA]
	var totalRange, sliceB []T
	if newFirstIdx < firstIdx {
		totalRange = slice.data[newFirstIdx : lastIdx+1]
		sliceB = slice.data[newFirstIdx:firstIdx]
	} else {
		totalRange = slice.data[firstIdx : newFirstIdx+lenA]
		sliceB = slice.data[lastIdx+1 : newFirstIdx+lenA]
	}
	slices.Reverse(sliceA)
	slices.Reverse(sliceB)
	slices.Reverse(totalRange)
}

// Return another SliceLike[T, I] that holds values in range [first, last]
//
// Analogous to slice[first:last+1]
func (slice *SliceAdapter[T]) Slice(firstIdx int, lastIdx int) (newSlice SliceLike[T, int]) {
	s := SliceAdapter[T]{
		data: slice.data[firstIdx : lastIdx+1],
	}
	return &s
}

// Return the first index in the slice.
//
// If the slice is empty, the index returned should
// result in `IdxValid(idx) == false`
func (slice *SliceAdapter[T]) FirstIdx() (idx int) {
	return 0
}

// Return the last index in the slice.
//
// If the slice is empty, the index returned should
// result in `IdxValid(idx) == false`
func (slice *SliceAdapter[T]) LastIdx() (idx int) {
	return len(slice.data) - 1
}

// Return the next index after the current index in the slice.
//
// If the given index is invalid or no next index exists,
// the index returned should result in `IdxValid(idx) == false`
func (slice *SliceAdapter[T]) NextIdx(thisIdx int) (nextIdx int) {
	return thisIdx + 1
}

// Return the index `n` places after the current index in the slice.
//
// If the given index is invalid or no nth next index exists,
// the index returned should result in `IdxValid(idx) == false`
func (slice *SliceAdapter[T]) NthNextIdx(thisIdx int, n int) (nthNextIdx int) {
	return thisIdx + n
}

// Return the prev index before the current index in the slice.
//
// If the given index is invalid or no prev index exists,
// the index returned should result in `IdxValid(idx) == false`
func (slice *SliceAdapter[T]) PrevIdx(thisIdx int) (prevIdx int) {
	return thisIdx - 1
}

// Return the index `n` places before the current index in the slice.
//
// If the given index is invalid or no nth previous index exists,
// the index returned should result in `IdxValid(idx) == false`
func (slice *SliceAdapter[T]) NthPrevIdx(thisIdx int, n int) (nthPrevIdx int) {
	return thisIdx - n

}

// Return the current number of values in the slice/list
//
// It is not guaranteed that all indexes less than `len` are valid for the slice
func (slice *SliceAdapter[T]) Len() int {
	return len(slice.data)
}

// Return the number of items between (and including) `firstIdx` and `lastIdx`
func (slice *SliceAdapter[T]) LenBetween(firstIdx int, lastIdx int) int {
	return min((lastIdx-firstIdx)+1, len(slice.data))
}

// Ensure at least `n` empty capacity spaces exist to add new items without reallocating
// the memory or perform any other expensive reorganization procedure
//
// If free space cannot be ensured and attempting to add `nMoreItems`
// will definitely fail or cause undefined behaviour, `ok == false`
func (slice *SliceAdapter[T]) TryEnsureFreeSlots(nMoreItems int) (ok bool) {
	slice.data = slices.Grow(slice.data, nMoreItems)
	ok = true
	return
}

// Insert `n` new slots at existing index, shifting all existing items
// after them forward. Returns the first new slot and the last new slot, inclusive.
//
// The implementation should assume no checks need need to be made to ensure free space exists,
// and calling this should not perform any reallocation
//
// May or may not work if attempting to insert items at an invalid indexIDX
func (slice *SliceAdapter[T]) InsertSlotsAssumeCapacity(idx int, count int) (firstNewSlot int, lastNewSlot int) {
	firstNewSlot = idx
	lastNewSlot = firstNewSlot + count - 1
	slice.data = slices.Insert(slice.data, idx, make([]T, count)...)
	return
}

// Append `n` new slots at the end of the list.
//
// Returns the first new slot and the last new slot, inclusive.
//
// The implementation should assume no checks need need to be made to ensure free space exists,
// and calling this should not perform any reallocation
func (slice *SliceAdapter[T]) AppendSlotsAssumeCapacity(count int) (firstNewSlot int, lastNewSlot int) {
	firstNewSlot = len(slice.data)
	lastNewSlot = firstNewSlot + count - 1
	slice.data = (slice.data)[:len(slice.data)+count]
	return
}

// Remove all items between `firstRemoveIdx` and `lastRemovedIdx`, inclusive
//
// All items after `lastRemovedIdx` are shifted backward
func (slice *SliceAdapter[T]) DeleteRange(firstRemovedIdx int, lastRemovedIdx int) {
	slice.data = slices.Delete(slice.data, firstRemovedIdx, lastRemovedIdx+1)
}

// Reset list to an empty state. The list's capacity may or may not be retained.
func (slice *SliceAdapter[T]) Clear() {
	slice.data = slice.data[:0]
}

// Return the total number of values the slice/list can hold
func (slice *SliceAdapter[T]) Cap() int {
	return cap(slice.data)
}

// GetPtr implements MemQueueLike.
func (slice *SliceAdapter[T]) GetPtr(idx int) *T {
	return &slice.data[idx]
}

// OffsetStart implements MemQueueLike.
func (slice *SliceAdapter[T]) IncrementStart(n int) {
	slice.data = slice.data[n:]
}

func (slice *SliceAdapter[T]) WriteAt(src []T, off int64) (n int, err error) {
	if off < 0 || off > int64(len(slice.data)) {
		return 0, io.EOF
	}
	end := int(off) + len(src)
	if end >= slice.Len() {
		slice.data = slice.data[:off]
		slice.data = append(slice.data, src...)
		n = len(src)
	} else {
		n = copy(slice.data[off:], src)
	}
	return n, nil
}

func (slice *SliceAdapter[T]) ReadAt(dest []T, off int64) (n int, err error) {
	n = copy(dest, slice.data[off:])
	if n < len(dest) {
		err = io.EOF
	}
	return n, err
}

func (slice *SliceAdapter[T]) Write(src []T) (n int, err error) {
	slice.data = append(slice.data, src...)
	return len(src), nil
}

func (slice *SliceAdapter[T]) Read(dest []T) (n int, err error) {
	n = copy(dest, slice.data)
	if n == 0 {
		err = io.EOF
	}
	return n, err
}

func (slice *SliceAdapter[T]) GoSlice() []T {
	return slice.data
}

var _ MemQueueLike[byte, int] = (*SliceAdapter[byte])(nil)
var _ MemListLike[byte, int] = (*SliceAdapter[byte])(nil)
var _ GoSliceLike[byte] = (*SliceAdapter[byte])(nil)
var _ io.Reader = (*SliceAdapter[byte])(nil)
var _ io.Writer = (*SliceAdapter[byte])(nil)
var _ io.ReaderAt = (*SliceAdapter[byte])(nil)
var _ io.WriterAt = (*SliceAdapter[byte])(nil)
