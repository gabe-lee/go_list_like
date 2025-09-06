package go_list_like

import (
	"io"
	"os"
	"syscall"
	"time"
)

type FileAdapter struct {
	File *os.File
}

func NewFileAdapter(file *os.File) FileAdapter {
	return FileAdapter{
		File: file,
	}
}

// SliceLike

func (f FileAdapter) PreferLinearOps() bool {
	return false
}

func (slice FileAdapter) ConsecutiveIndexesInOrder() bool {
	return true
}
func (slice FileAdapter) AllIndexesLessThanLenValid() bool {
	return true
}

// Returns whether the given index is valid for the slice
func (f FileAdapter) IdxValid(idx int) bool {
	l, ok := f.LenIfValid()
	if !ok {
		return false
	}
	return idx >= 0 && idx < l
}

// Returns whether the given index range is valid for the slice
//
// The following MUST be true:
//   - `firstIdx` comes logically before OR is equal to `lastIdx`
//   - all indexes including and between `firstIdx` and `lastIdx` are valid for the slice
func (f FileAdapter) RangeValid(firstIdx int, lastIdx int) bool {
	l, ok := f.LenIfValid()
	if !ok || l == 0 {
		return false
	}
	return firstIdx >= 0 && firstIdx <= lastIdx && lastIdx < l
}

// Split an index range in half, returning the index in the middle of the range
//
// Assumes `RangeValid(firstIdx, lastIdx) == true`
func (f FileAdapter) SplitRange(firstIdx int, lastIdx int) (middleIdx int) {
	return firstIdx + ((lastIdx - firstIdx) >> 1)
}

// Get the value at the provided index
func (f FileAdapter) Get(idx int) (val byte) {
	var b [1]byte
	f.ReadAt(b[:], int64(idx))
	return b[0]
}

// Set the value at the provided index to the given value
func (f FileAdapter) Set(idx int, val byte) {
	var b [1]byte = [1]byte{val}
	f.WriteAt(b[:], int64(idx))
}

// Move only the data contained at `oldIdx` to `newIdx`,
// overwriting the existing data at `newIdx`
func (f FileAdapter) Move(oldIdx int, newIdx int) {
	val := f.Get(oldIdx)
	if newIdx < oldIdx {
		prevIdx := oldIdx - 1
		for oldIdx > newIdx {
			Overwrite(f, prevIdx, oldIdx)
			oldIdx = prevIdx
			prevIdx -= 1
		}
	} else {
		nextIdx := oldIdx + 1
		for oldIdx < newIdx {
			Overwrite(f, nextIdx, oldIdx)
			oldIdx = nextIdx
			nextIdx += 1
		}
	}
	f.Set(newIdx, val)
}

// Remove all data contained in range `firstIdx` to `lastIdx` (inclusive),
// and re-insert it at the `newFirstIdx` position
func (f FileAdapter) MoveRange(firstIdx int, lastIdx int, newFirstIdx int) {
	lenA := (lastIdx - firstIdx) + 1
	sliceA := f.Slice(firstIdx, (firstIdx+lenA)-1)
	var totalRange, sliceB SliceLike[byte, int]
	if newFirstIdx < firstIdx {
		totalRange = f.Slice(newFirstIdx, lastIdx)
		sliceB = f.Slice(newFirstIdx, firstIdx-1)
	} else {
		totalRange = f.Slice(firstIdx, (newFirstIdx+lenA)-1)
		sliceB = f.Slice(lastIdx+1, (newFirstIdx+lenA)-1)
	}
	Reverse(sliceA)
	Reverse(sliceB)
	Reverse(totalRange)
}

// Return another SliceLike[T, I] that holds values in range [first, last]
//
// Analogous to slice[first:last+1]
func (f FileAdapter) Slice(firstIdx int, lastIdx int) (slice SliceLike[byte, int]) {
	return &FileSliceAdapter{
		FAdapter: f,
		start:    firstIdx,
		len:      (lastIdx - firstIdx) + 1,
	}
}

// Return the first index in the slice.
//
// If the slice is empty, the index returned should
// result in `IdxValid(idx) == false`
func (f FileAdapter) FirstIdx() (idx int) {
	l, ok := f.LenIfValid()
	if !ok || l == 0 {
		return -1
	}
	return 0
}

// Return the last index in the slice.
//
// If the slice is empty, the index returned should
// result in `IdxValid(idx) == false`
func (f FileAdapter) LastIdx() (idx int) {
	l, ok := f.LenIfValid()
	if !ok || l == 0 {
		return -1
	}
	return l - 1
}

// Return the next index after the current index in the slice.
//
// If the given index is invalid or no next index exists,
// the index returned should result in `IdxValid(idx) == false`
func (f FileAdapter) NextIdx(thisIdx int) (nextIdx int) {
	return thisIdx + 1
}

// Return the index `n` places after the current index in the slice.
//
// If the given index is invalid or no nth next index exists,
// the index returned should result in `IdxValid(idx) == false`
func (f FileAdapter) NthNextIdx(thisIdx int, n int) (nthNextIdx int) {
	return thisIdx + n
}

// Return the prev index before the current index in the slice.
//
// If the given index is invalid or no prev index exists,
// the index returned should result in `IdxValid(idx) == false`
func (f FileAdapter) PrevIdx(thisIdx int) (prevIdx int) {
	return thisIdx - 1
}

// Return the index `n` places before the current index in the slice.
//
// If the given index is invalid or no nth previous index exists,
// the index returned should result in `IdxValid(idx) == false`
func (f FileAdapter) NthPrevIdx(thisIdx int, n int) (nthPrevIdx int) {
	return thisIdx - n
}

// Return the current number of values in the slice/list
//
// It is not guaranteed that all indexes less than `len` are valid for the slice
func (f FileAdapter) Len() int {
	l, _ := f.LenIfValid()
	return l
}

// Returns the file length (size), and whether the file is valid
func (f FileAdapter) LenIfValid() (l int, valid bool) {
	stat, err := f.File.Stat()
	if err != nil {
		return 0, false
	}
	return int(stat.Size()), true
}

// Return the number of items between (and including) `firstIdx` and `lastIdx`
func (f FileAdapter) LenBetween(firstIdx int, lastIdx int) int {
	l, ok := f.LenIfValid()
	if !ok || l == 0 {
		return 0
	}
	return (lastIdx - firstIdx) + 1
}

// ListLike

// Ensure at least `n` empty capacity spaces exist to add new items without reallocating
// the memory or perform any other expensive reorganization procedure
//
// If free space cannot be ensured and attempting to add `nMoreItems`
// will definitely fail or cause undefined behaviour, `ok == false`
func (f FileAdapter) TryEnsureFreeSlots(nMoreItems int) (ok bool) {
	_, ok = f.LenIfValid()
	return
}

// Insert `n` new slots directly before existing index, shifting all existing items
// after them forward.
//
// Returns the first new slot and the last new slot, inclusive.
//
// The implementation should assume no checks need need to be made to ensure free space exists,
// and calling this should not perform any reallocation(ok bool)
//
// May or may not work if attempting to insert items at an invalid index
func (f FileAdapter) InsertSlotsAssumeCapacity(idx int, count int) (firstNewSlot int, lastNewSlot int) {
	l, ok := f.LenIfValid()
	if !ok {
		return 0, 0
	}
	ridx := l - 1
	newSize := l + count
	f.File.Truncate(int64(newSize))
	widx := newSize - 1
	for ridx >= idx {
		Overwrite(f, ridx, widx)
		ridx -= 1
		widx -= 1
	}
	firstNewSlot = idx
	lastNewSlot = (firstNewSlot + count) - 1
	return
}

// Append `n` new slots at the end of the list.
//
// Returns the first new slot and the last new slot, inclusive.
//
// The implementation should assume no checks need need to be made to ensure free space exists,
// and calling this should not perform any reallocation
func (f FileAdapter) AppendSlotsAssumeCapacity(count int) (firstNewSlot int, lastNewSlot int) {
	l, ok := f.LenIfValid()
	if !ok {
		return 0, 0
	}
	firstNewSlot = l
	newSize := l + count
	f.File.Truncate(int64(newSize))
	lastNewSlot = newSize - 1
	return
}

// Remove all items between `firstRemoveIdx` and `lastRemovedIdx`, inclusive
//
// All items after `lastRemovedIdx` are shifted backward
func (f FileAdapter) DeleteRange(firstRemovedIdx int, lastRemovedIdx int) {
	stat, err := f.File.Stat()
	if err != nil {
		return
	}
	size := int(stat.Size())
	ridx := lastRemovedIdx + 1
	widx := firstRemovedIdx
	for ridx < size {
		Overwrite(f, ridx, widx)
		ridx += 1
		widx += 1
	}
	f.Truncate(int64(size - ((lastRemovedIdx - firstRemovedIdx) + 1)))
}

// Reset list to an empty state. The list's capacity may or may not be retained.
func (f FileAdapter) Clear() {
	f.Truncate(0)
}

// Return the total number of values the slice/list can hold
func (f FileAdapter) Cap() int {
	stat, err := f.File.Stat()
	if err != nil {
		return 0
	}
	return int(stat.Size())
}

// File aliases

func (f FileAdapter) Chdir() error {
	return f.File.Chdir()
}
func (f FileAdapter) Chmod(mode os.FileMode) error {
	return f.File.Chmod(mode)
}
func (f FileAdapter) Chown(uid int, gid int) error {
	return f.File.Chown(uid, gid)
}
func (f FileAdapter) Close() error {
	return f.File.Close()
}
func (f FileAdapter) Fd() uintptr {
	return f.File.Fd()
}
func (f FileAdapter) Name() string {
	return f.File.Name()
}
func (f FileAdapter) Read(b []byte) (n int, err error) {
	return f.File.Read(b)
}
func (f FileAdapter) ReadAt(b []byte, off int64) (n int, err error) {
	return f.File.ReadAt(b, off)
}
func (f FileAdapter) ReadDir(n int) ([]os.DirEntry, error) {
	return f.File.ReadDir(n)
}
func (f FileAdapter) ReadFrom(r io.Reader) (n int64, err error) {
	return f.File.ReadFrom(r)
}
func (f FileAdapter) Readdir(n int) ([]os.FileInfo, error) {
	return f.File.Readdir(n)
}
func (f FileAdapter) Readdirnames(n int) (names []string, err error) {
	return f.File.Readdirnames(n)
}
func (f FileAdapter) Seek(offset int64, whence int) (ret int64, err error) {
	return f.File.Seek(offset, whence)
}
func (f FileAdapter) SetDeadline(t time.Time) error {
	return f.File.SetDeadline(t)
}
func (f FileAdapter) SetReadDeadline(t time.Time) error {
	return f.File.SetReadDeadline(t)
}
func (f FileAdapter) SetWriteDeadline(t time.Time) error {
	return f.File.SetWriteDeadline(t)
}
func (f FileAdapter) Stat() (os.FileInfo, error) {
	return f.File.Stat()
}
func (f FileAdapter) Sync() error {
	return f.File.Sync()
}
func (f FileAdapter) SyscallConn() (syscall.RawConn, error) {
	return f.File.SyscallConn()
}
func (f FileAdapter) Truncate(size int64) error {
	return f.File.Truncate(size)
}
func (f FileAdapter) Write(b []byte) (n int, err error) {
	return f.File.Write(b)
}
func (f FileAdapter) WriteAt(b []byte, off int64) (n int, err error) {
	return f.File.WriteAt(b, off)
}
func (f FileAdapter) WriteString(s string) (n int, err error) {
	return f.File.WriteString(s)
}
func (f FileAdapter) WriteTo(w io.Writer) (n int64, err error) {
	return f.File.WriteTo(w)
}

var _ ListLike[byte, int] = FileAdapter{}
var _ io.Reader = FileAdapter{}
var _ io.Writer = FileAdapter{}
var _ io.ReaderAt = FileAdapter{}
var _ io.WriterAt = FileAdapter{}

type FileSliceAdapter struct {
	FAdapter FileAdapter
	start    int
	len      int
}

func (f *FileSliceAdapter) PreferLinearOps() bool {
	return false
}

func (slice *FileSliceAdapter) ConsecutiveIndexesInOrder() bool {
	return true
}
func (slice *FileSliceAdapter) AllIndexesLessThanLenValid() bool {
	return true
}

// Returns whether the given index is valid for the slice
func (f *FileSliceAdapter) IdxValid(idx int) bool {
	return idx >= 0 && idx < f.Len()
}

// Returns whether the given index range is valid for the slice
//
// The following MUST be true:
//   - `firstIdx` comes logically before OR is equal to `lastIdx`
//   - all indexes including and between `firstIdx` and `lastIdx` are valid for the slice
func (f *FileSliceAdapter) RangeValid(firstIdx int, lastIdx int) bool {
	return firstIdx >= 0 && firstIdx <= lastIdx && lastIdx < f.Len()
}

// Split an index range in half, returning the index in the middle of the range
//
// Assumes `RangeValid(firstIdx, lastIdx) == true`
func (f *FileSliceAdapter) SplitRange(firstIdx int, lastIdx int) (middleIdx int) {
	return firstIdx + ((lastIdx - firstIdx) >> 1)
}

// Get the value at the provided index
func (f *FileSliceAdapter) Get(idx int) (val byte) {
	var b [1]byte
	f.FAdapter.ReadAt(b[:], int64(f.start+idx))
	return b[0]
}

// Set the value at the provided index to the given value
func (f *FileSliceAdapter) Set(idx int, val byte) {
	var b [1]byte = [1]byte{val}
	f.FAdapter.WriteAt(b[:], int64(f.start+idx))
}

// Move only the data contained at `oldIdx` to `newIdx`,
// overwriting the existing data at `newIdx`
func (f *FileSliceAdapter) Move(oldIdx int, newIdx int) {
	val := f.Get(oldIdx)
	if newIdx < oldIdx {
		prevIdx := oldIdx - 1
		for oldIdx >= newIdx {
			Overwrite(f, prevIdx, oldIdx)
			oldIdx = prevIdx
			prevIdx -= 1
		}
	} else {
		nextIdx := oldIdx + 1
		for oldIdx <= newIdx {
			Overwrite(f, nextIdx, oldIdx)
			oldIdx = nextIdx
			nextIdx += 1
		}
	}
	f.Set(newIdx, val)
}

// Remove all data contained in range `firstIdx` to `lastIdx` (inclusive),
// and re-insert it at the `newFirstIdx` position
func (f *FileSliceAdapter) MoveRange(firstIdx int, lastIdx int, newFirstIdx int) {
	f.FAdapter.MoveRange(f.start+firstIdx, f.start+lastIdx, f.start+newFirstIdx)
}

// Return another SliceLike[T, I] that holds values in range [first, last]
//
// Analogous to slice[first:last+1]
func (f *FileSliceAdapter) Slice(firstIdx int, lastIdx int) (slice SliceLike[byte, int]) {
	return &FileSliceAdapter{
		FAdapter: f.FAdapter,
		start:    firstIdx,
		len:      (lastIdx - firstIdx) + 1,
	}
}

// Return the first index in the slice.
//
// If the slice is empty, the index returned should
// result in `IdxValid(idx) == false`
func (f *FileSliceAdapter) FirstIdx() (idx int) {
	return 0
}

// Return the last index in the slice.
//
// If the slice is empty, the index returned should
// result in `IdxValid(idx) == false`
func (f *FileSliceAdapter) LastIdx() (idx int) {
	return f.Len() - 1
}

// Return the next index after the current index in the slice.
//
// If the given index is invalid or no next index exists,
// the index returned should result in `IdxValid(idx) == false`
func (f *FileSliceAdapter) NextIdx(thisIdx int) (nextIdx int) {
	return thisIdx + 1
}

// Return the index `n` places after the current index in the slice.
//
// If the given index is invalid or no nth next index exists,
// the index returned should result in `IdxValid(idx) == false`
func (f *FileSliceAdapter) NthNextIdx(thisIdx int, n int) (nthNextIdx int) {
	return thisIdx + n
}

// Return the prev index before the current index in the slice.
//
// If the given index is invalid or no prev index exists,
// the index returned should result in `IdxValid(idx) == false`
func (f *FileSliceAdapter) PrevIdx(thisIdx int) (prevIdx int) {
	return thisIdx - 1
}

// Return the index `n` places before the current index in the slice.
//
// If the given index is invalid or no nth previous index exists,
// the index returned should result in `IdxValid(idx) == false`
func (f *FileSliceAdapter) NthPrevIdx(thisIdx int, n int) (nthPrevIdx int) {
	return thisIdx - n
}

// Return the current number of values in the slice/list
//
// It is not guaranteed that all indexes less than `len` are valid for the slice
func (f *FileSliceAdapter) Len() int {
	return f.len
}

// Return the number of items between (and including) `firstIdx` and `lastIdx`
func (f *FileSliceAdapter) LenBetween(firstIdx int, lastIdx int) int {
	return (lastIdx - firstIdx) + 1
}

// Increment the start location (index/pointer/etc.) of this queue by
// `n` positions. The new 'first' item in the queue should be the item
// previously located at index `delta`
func (f *FileSliceAdapter) IncrementStart(n int) {
	f.start += n
	f.len -= n
}

func (f *FileSliceAdapter) Read(b []byte) (n int, err error) {
	minLen := min(f.len, len(b))
	n, err = f.FAdapter.ReadAt(b[:minLen], int64(f.start))
	if err == nil && n == 0 {
		err = io.EOF
	}
	return
}

func (f *FileSliceAdapter) WriteAt(b []byte, off int64) (n int, err error) {
	maxOffset := min(f.len, int(off))
	maxLen := f.len - maxOffset
	n, err = f.FAdapter.WriteAt(b[:maxLen], int64(f.start+maxOffset))
	if err == nil && n < len(b) {
		err = io.EOF
	}
	return
}
func (f *FileSliceAdapter) ReadAt(b []byte, off int64) (n int, err error) {
	maxOffset := min(f.len, int(off))
	maxLen := f.len - maxOffset
	n, err = f.FAdapter.ReadAt(b[:maxLen], int64(f.start+maxOffset))
	if err == nil && n < len(b) {
		err = io.EOF
	}
	return
}

var _ QueueLike[byte, int] = (*FileSliceAdapter)(nil)
var _ io.Reader = (*FileSliceAdapter)(nil)
var _ io.ReaderAt = (*FileSliceAdapter)(nil)
var _ io.WriterAt = (*FileSliceAdapter)(nil)
