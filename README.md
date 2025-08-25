# go_list_like
Interface for user-defined types or wrappers around 3rd-party types that can behave like a traditional 'List' or 'Vector'
  - [Why](#why)
  - [Installation](#installation)
  - [Example](#example)
  - [Interfaces](#interfaces)
    - [SliceLike[T]](#sliceliket)
    - [ListLike[T]](#listliket)
    - [QueueLike[T]](#queueliket)
    - [MemSliceLike[T]](#memsliceliket)
  - [Additionl Functions](#additional-functions)
## Why?
This package provides a simple interface for things that can behave like a 'List' or 'Vector', but may be user-defined data structures with complex memory layouts and/or traversal methods. By implementing these interfaces, one can avoid an intermediate temporary slice to use as an interface between two different incompatible data structures.

It provides a simple wrapper type for Golang slices themselves to adapt them automatically.

[Back to Top](#go_list_like)
## Installation
Run this command from your project directory
```
go get github.com/gabe-lee/go_list_like@latest
```
Then implement one or more of the defined interfaces for your data structure in order to access a collection of additional functions that can automatically operate on your data structure

The provided wrapper types `SliceAdapter[T]` and `SliceAdapterIndirect[T]` can be used with a standard golang slice `[]T`, which implements:
  - SliceLike[T] + MemSliceLike[T]
  - ListLike[T] + MemListLike[T]
  - QueueLike[T] + MemQueueLike[T]
  - io.Reader, io.Writer, io.ReaderAt, and io.WriterAt in a generic manner for all types

Also, the provided wrapper type `FileAdapter` can be used with a standard golang file `os.File`, which implements:
  - SliceLike[byte] 
  - ListLike[byte]
  - Also all other interfaces normally implemented by `os.File`
  - `FileAdapter.Slice(start, end)` returns a `*FileSliceAdapter` that implements:
    - SliceLike[byte]
	- QueueLike[byte]
	- io.Reader, io.ReaderAt, io.WriterAt

[Back to Top](#go_list_like)
## Example
```golang
import (
    "fmt"
    ll "github.com/gabe-lee/go_list_like"
)

func main() {
    mySlice := []byte("Hello World")
	mySliceLike := ll.NewSliceAdapterIndirect(&mySlice)
	fmt.Printf("%s\n", mySlice)
	ll.AppendV(mySliceLike, '!')
	ll.Delete(mySliceLike, 5, 1)
	fmt.Printf("%s\n", mySlice)
}
```
Output:
```
Hello World
HelloWorld!
```

[Back to Top](#go_list_like)
## Interfaces
### SliceLike[T]
Implement:
```golang
type SliceLike[T any] interface {
	// Get the value at the provided index
	Get(idx int) (val T)
	// Set the value at the provided index to the given value
	Set(idx int, val T)
	// Return another SliceLike[T] that holds values in range [start:end),
	// where subSlice[0] == slice[start]
	//
	// Analogous to slice[start:end]
	Slice(start int, end int) SliceLike[T]
	// Return the current number of values in the slice/list
	//
	// All values less than this length MUST be valid for Get() and Set()
	Len() int
}
```
To get:
```golang
func Len[T any](slice SliceLike[T]) int 
func Slice[T any](slice SliceLike[T], start int, end int) SliceLike[T]
func IdxInRange[T any](slice SliceLike[T], idx int) bool
func IdxUnderLen[T any](slice SliceLike[T], idx int) bool
func IsEmpty[T any](slice SliceLike[T]) bool
func Get[T any](slice SliceLike[T], idx int) T
func TryGet[T any](slice SliceLike[T], idx int) (val T, ok bool)
func LastIdx[T any](slice SliceLike[T]) (lastIdx int)
func GetLast[T any](slice SliceLike[T]) T
func TryGetLast[T any](slice SliceLike[T]) (val T, ok bool)
func GetFirst[T any](slice SliceLike[T]) T
func TryGetFirst[T any](slice SliceLike[T]) (val T, ok bool)
func Set[T any](slice SliceLike[T], idx int, val T)
func TrySet[T any](slice SliceLike[T], idx int, val T) (ok bool)
func SetLast[T any](slice SliceLike[T], val T)
func TrySetLast[T any](slice SliceLike[T], val T) (ok bool)
func SetFirst[T any](slice SliceLike[T], val T)
func TrySetFirst[T any](slice SliceLike[T], val T) (ok bool)
func SetFrom[T any](dest SliceLike[T], destIdx int, source SliceLike[T], srcIdx int)
func TrySetFrom[T any](dest SliceLike[T], destIdx int, source SliceLike[T], srcIdx int) (ok bool)
func Swap[T any](slice SliceLike[T], idxA int, idxB int)
func TrySwap[T any](slice SliceLike[T], idxA int, idxB int) (ok bool)
func Move[T any](slice SliceLike[T], oldIdx int, newIdx int)
func TryMove[T any](slice SliceLike[T], oldIdx int, newIdx int) (ok bool)
func Copy[T any](dest SliceLike[T], destStart int, source SliceLike[T], srcStart int, copyLen int) (nCopied int)
func TryCopy[T any](dest SliceLike[T], destStart int, source SliceLike[T], srcStart int, copyLen int) (nCopied int, ok bool) 
func CopyScalar[T any](dest SliceLike[T], destStart int, copyLen int, val T) (nCopied int)
func TryCopyScalar[T any](dest SliceLike[T], destStart int, copyLen int, val T) (nCopied int, ok bool)
func Swizzle[T any, I Index](slices SliceLike[SliceLike[T]], selectors SliceLike[I], outputList ListLike[T])
func TrySwizzle[T any, I Index](slices SliceLike[SliceLike[T]], selectors SliceLike[I], outputList ListLike[T]) (ok bool)
func IsSorted[T any](slice SliceLike[T], greaterThan func(a T, b T) bool) bool
func IsSortedImplicit[T cmp.Ordered](slice SliceLike[T]) bool
func Sort[T any](slice SliceLike[T], greaterThan func(a T, b T) bool)
func SortImplicit[T cmp.Ordered](slice SliceLike[T])
func DoActionOnItemsUntilFalse[T any](slice SliceLike[T], action func(slice SliceLike[T], idx int, item T) (shouldContinue bool)) (stopIdx int)
func DoActionOnItemsUntilFalseWithUserdata[T any, U any](slice SliceLike[T], action func(slice SliceLike[T], idx int, item T, userdata *U) (shouldContinue bool), userdata *U) (stopIdx int)
func DoActionOnAllItems[T any](slice SliceLike[T], action func(slice SliceLike[T], idx int, item T))
func DoActionOnAllItemsWithUserdata[T any, U any](slice SliceLike[T], action func(slice SliceLike[T], idx int, item T, userdata *U), userdata *U) 
func FilterIndexes[T any, I Index](slice SliceLike[T], selectFunc func(slice SliceLike[T], idx I, item T) bool, outputList ListLike[I])
func FilterIndexesWithUserdata[T any, I Index, U any](slice SliceLike[T], selectFunc func(slice SliceLike[T], idx I, item T, userdata *U) bool, outputList ListLike[I], userdata *U)
func MapValues[T any, TT any](slice SliceLike[T], mapFunc func(slice SliceLike[T], idx int, item T) TT, outputList ListLike[TT])
func MapValuesWithUserdata[T any, TT any, U any](slice SliceLike[T], mapFunc func(slice SliceLike[T], idx int, item T, userdata *U) TT, outputList ListLike[TT], userdata *U)
func Accumulate[T any, TT any](slice SliceLike[T], initialAccumulation TT, accumulate func(slice SliceLike[T], idx int, item T, currentAccumulation TT) (newAccumulation TT)) (finalAccumulation TT)
func AccumulateWithUserdata[T any, TT any, U any](slice SliceLike[T], initialAccumulation TT, accumulate func(slice SliceLike[T], idx int, item T, currentAccumulation TT, userdata *U) (newAccumulation TT), userdata *U) (finalAccumulation TT)
```

[Back to Top](#go_list_like)
### ListLike[T]
Implement:
```golang
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
```
To get:
```golang
// SliceLike[T] funcs...
func ChangeLen[T any](list ListLike[T], delta int)
func GrowLen[T any](list ListLike[T], grow int)
func ShrinkLen[T any](list ListLike[T], shrink int)
func GrowCap[T any](list ListLike[T], n int) 
func GrowCapIfNeeded[T any](list ListLike[T], nMoreItems int)
func Cap[T any](list ListLike[T]) int
func Clear[T any](list ListLike[T])
func AppendSlots[T any](list ListLike[T], count int) SliceLike[T]
func AppendV[T any](list ListLike[T], vals ...T) 
func AppendGetStartIdxV[T any](list ListLike[T], vals ...T) (startIdx int)
func Append[T any](list ListLike[T], vals SliceLike[T])
func AppendGetStartIdx[T any](list ListLike[T], vals SliceLike[T]) (startIdx int)
func InsertSlots[T any](list ListLike[T], idx int, count int) SliceLike[T]
func InsertV[T any](list ListLike[T], idx int, vals ...T)
func Insert[T any](list ListLike[T], idx int, vals SliceLike[T])
func Delete[T any](list ListLike[T], idx int, count int)
func DeleteSparse[T any, I Index](list ListLike[T], deleteIndexSlice SliceLike[I], sortDeleteIndexes bool)
func Remove[T any](list ListLike[T], idx int, count int, outputList ListLike[T])
func RemoveSparse[T any, I Index](list ListLike[T], removeIndexSlice SliceLike[I], sortRemoveIndexes bool, outputList ListLike[T])
func Replace[T any](dest ListLike[T], destStart, destLen int, source SliceLike[T], srcStart, srcLen int) (delta int)
func Push[T any](list ListLike[T], val T)
func PushGetIdx[T any](list ListLike[T], val T) (idx int)
func Pop[T any](list ListLike[T]) T
func TryPop[T any](list ListLike[T]) (val T, ok bool)
```

[Back to Top](#go_list_like)
### QueueLike[T]
Implement:
```golang
type QueueLike[T any] interface {
	SliceLike[T]
	// Offset the start location (index/pointer/etc.) of this queue by
	// the given delta. The new 'first' item in the queue should be the item
	// previously located at `queue.GetPtr(0+delta)`.
	OffsetStart(delta int)
}
```
To get:
```golang
// SliceLike[T] funcs...
func Dequeue[T any](queue QueueLike[T], count int, outputList ListLike[T])
func TryDequeue[T any](queue QueueLike[T], count int, outputList ListLike[T]) (ok bool)
func Peek[T any](queue QueueLike[T], count int) (peekSlice SliceLike[T]) 
func TryPeek[T any](queue QueueLike[T], count int) (peekSlice SliceLike[T], ok bool) 
func Discard[T any](queue QueueLike[T], count int)
func TryDiscard[T any](queue QueueLike[T], count int) (ok bool) 
```

[Back to Top](#go_list_like)
### MemSliceLike[T]
Implement:
```golang
type MemSliceLike[T any] interface {
	SliceLike[T]
	// Get the a pointer to the value at the provided index
	GetPtr(idx int) *T
}
```
To get:
```golang
// SliceLike[T] funcs...
func GetPtr[T any](memSliceLike MemSliceLike[T], idx int) *T
func TryGetPtr[T any](memSliceLike MemSliceLike[T], idx int) (ptr *T, ok bool)
```
## Additional Functions
In addition to the normal functions supplied to implementors of the above interfaces, a number of helper functions are included for less cumbersome common use cases:
```golang
// oldVal := slice[idx]
// slice[idx] = val
// return oldVal != val
func SetChanged[T Equatable](slice SliceLike[T], idx int, val T) (didChange bool) 
// slice[idx] = slice[idx] + val
func SetAdd[T Number](slice SliceLike[T], idx int, val T)
// return slice[idx] + val
func GetAdd[T Number](slice SliceLike[T], idx int, val T) (result T)
// slice[idx] = slice[idx] - val
func SetSubtract[T Number](slice SliceLike[T], idx int, val T)
// return slice[idx] - val
func GetSubtract[T Number](slice SliceLike[T], idx int, val T) (result T)
// slice[idx] = slice[idx] * val
func SetMultiply[T Number](slice SliceLike[T], idx int, val T)
// return slice[idx] * val
func GetMultiply[T Number](slice SliceLike[T], idx int, val T) (result T)
// slice[idx] = slice[idx] / val
func SetDivide[T Number](slice SliceLike[T], idx int, val T)
// return slice[idx] / val
func GetDivide[T Number](slice SliceLike[T], idx int, val T) (result T)
// slice[idx] = slice[idx] % val
func SetModulo[T Integer](slice SliceLike[T], idx int, val T)
// return slice[idx] % val
func GetModulo[T Integer](slice SliceLike[T], idx int, val T) (result T)
// return slice[idx] % val, slice[idx] - (slice[idx] % val)
func GetModRem[T Integer](slice SliceLike[T], idx int, val T) (mod T, rem T)
// slice[idx] = math.Mod(slice[idx], val)
func SetFModulo[T Float](slice SliceLike[T], idx int, val T)
// return math.Mod(slice[idx], val)
func GetFModulo[T Float](slice SliceLike[T], idx int, val T) (result T)
// return math.Mod(slice[idx], val), slice[idx] - math.Mod(slice[idx], val)
func GetFModRem[T Float](slice SliceLike[T], idx int, val T) (mod T, rem T)
// slice[idx] = slice[idx] & val
func SetBitAnd[T Integer](slice SliceLike[T], idx int, val T)
// return slice[idx] & val
func GetBitAnd[T Integer](slice SliceLike[T], idx int, val T) (result T)
// slice[idx] = slice[idx] & val
func SetBitOr[T Integer](slice SliceLike[T], idx int, val T)
// return slice[idx] & val
func GetBitOr[T Integer](slice SliceLike[T], idx int, val T) (result T)
// slice[idx] = slice[idx] ^ val
func SetBitXor[T Integer](slice SliceLike[T], idx int, val T)
// return slice[idx] ^ val
func GetBitXor[T Integer](slice SliceLike[T], idx int, val T) (result T)
// slice[idx] = ^slice[idx]
func SetBitInvert[T Integer](slice SliceLike[T], idx int)
// return ^slice[idx]
func GetBitInvert[T Integer](slice SliceLike[T], idx int) (result T)
// slice[idx] = slice[idx] << val
func SetBitLsh[T Integer](slice SliceLike[T], idx int, val T)
// return slice[idx] << val
func GetBitLsh[T Integer](slice SliceLike[T], idx int, val T) (result T)
// slice[idx] = slice[idx] >> val
func SetBitRsh[T Integer](slice SliceLike[T], idx int, val T)
// return slice[idx] >> val
func GetBitRsh[T Integer](slice SliceLike[T], idx int, val T) (result T)
// return slice[idx] < val
func GetLessThan[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool
// return slice[idx1] < slice[idx2]
func GetLessThan2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool
// return slice[idx] <= val
func GetLessThanEqual[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool
// return slice[idx1] <= slice[idx2]
func GetLessThanEqual2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool
// return slice[idx] > val
func GetGreaterThan[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool
// return slice[idx1] > slice[idx2]
func GetGreaterThan2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool
// return slice[idx] >= val
func GetGreaterThanEqual[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool
// return slice[idx1] >= slice[idx2]
func GetGreaterThanEqual2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool
// return slice[idx] == val
func GetEquals[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool
// return slice[idx1] == slice[idx2]
func GetEquals2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool
// return slice[idx] != val
func GetNotEquals[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool
// return slice[idx1] != slice[idx2]
func GetNotEquals2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool
// return min(slice[indexes[0]], slice[indexes[1]], ...)
func GetMinV[T cmp.Ordered](slice SliceLike[T], indexes ...int) T
// return min(slice[indexes[0]], slice[indexes[1]], ...)
func GetMin[T cmp.Ordered, I Index](slice SliceLike[T], indexes SliceLike[I]) T
// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMinV[T cmp.Ordered](slice SliceLike[T], setIdx int, indexes ...int)
// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMin[T cmp.Ordered, I Index](slice SliceLike[T], setIdx int, indexes SliceLike[I])
// return max(slice[indexes[0]], slice[indexes[1]], ...)
func GetMaxV[T cmp.Ordered](slice SliceLike[T], indexes ...int) T
// return max(slice[indexes[0]], slice[indexes[1]], ...)
func GetMax[T cmp.Ordered, I Index](slice SliceLike[T], indexes SliceLike[I]) T
// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMaxV[T cmp.Ordered](slice SliceLike[T], setIdx int, indexes ...int)
// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMax[T cmp.Ordered, I Index](slice SliceLike[T], setIdx int, indexes SliceLike[I])
// return min(maxVal, max(slice[idx], minVal))
func GetClamped[T cmp.Ordered](slice SliceLike[T], idx int, minVal T, maxVal T) T
// slice[idx] = min(maxVal, max(slice[idx], minVal))
func SetClamped[T cmp.Ordered](slice SliceLike[T], idx int, minVal T, maxVal T)
// return *(*TT)(unsafe.Pointer(&slice[idx]))
func GetUnsafeCast[T any, TT any](slice SliceLike[T], idx int) (val TT)
// return *(*TT)(unsafe.Pointer(&slice[idx]))
func GetUnsafePtrCast[T any, TT any](slice MemSliceLike[T], idx int) (val *TT)
// slice[idx] = *(*T)(unsafe.Pointer(&val))
func SetUnsafeCast[T any, TT any](slice SliceLike[T], idx int, val TT)
// val_T := *(*T)(unsafe.Pointer(&val))
// oldVal_TT := *(*TT)(unsafe.Pointer(&slice[idx]))
// slice[idx] = val_T
// return oldVal_TT != val
func SetUnsafeCastChanged[T any, TT Equatable](slice SliceLike[T], idx int, val TT) (didChange bool)
// val_T := *(*T)(unsafe.Pointer(&val))
// oldVal_T := slice[idx]
// slice[idx] = val_T
// return oldVal_T != val_T
func SetUnsafeCastChangedAlt[T Equatable, TT any](slice SliceLike[T], idx int, val TT) (didChange bool)
// Decode a rune from the byte slice at the given index
func GetRune(slice SliceLike[byte], idx int) (r rune, bytes int, ok bool)
// Write a rune to the byte slice at given index
func SetRune(slice SliceLike[byte], idx int, r rune) (bytes int, ok bool)
// Append a rune to the end of the byte slice
func AppendRune(list ListLike[byte], r rune) (bytes int, ok bool)
```