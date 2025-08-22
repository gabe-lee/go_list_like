# go_list_like
Interface for user-defined types or wrappers around 3rd-party types that can behave like a traditional 'List' or 'Vector'
  - [Why](#why)
  - [Installation](#installation)
  - [Interfaces](#interfaces)
    - [SliceLike[T]](#sliceliket)
    - [ListLike[T]](#listliket)
    - [QueueLike[T]](#queueliket)
    - [FwdTraversable[T]](#fwdtraversablet)
	- [FwdLinkedListLike[T]](#fwdlinkedlistliket) (not yet used)
	- [RevTraversable[T]](#revtraversablet) (not yet used)
	- [RevLinkedListLike[T]](#revlinkedlistliket) (not yet used)
  - [Example](#example)
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

As a special case, the provided wrapper type `SliceAdapter[T]` can be used with a standard golang slice `[]T`, which implements:
  - SliceLike[T]
  - ListLike[T]
  - QueueLike[T]
  - FwdTraversable[T]
  - RevTraversable[T]

[Back to Top](#go_list_like)
## Interfaces
### SliceLike[T]
Implement:
```golang
type SliceLike[T any] interface {
	// Return a pointer to the value at index `idx`
	GetPtr(idx int) *T
	// Return the current number of values in the slice/list
	//
	// All indexes less than this value should be valid for `GetPtr(idx)`
	Len() int
}
```
To get:
```golang
func Len[T any](sliceLike SliceLike[T]) int
func Get[T any](sliceLike SliceLike[T], idx int) T
func TryGet[T any](sliceLike SliceLike[T], idx int) (val T, ok bool)
func GetPtr[T any](sliceLike SliceLike[T], idx int) *T
func TryGetPtr[T any](sliceLike SliceLike[T], idx int) (val *T, ok bool)
func GetLast[T any](sliceLike SliceLike[T]) T 
func TryGetLast[T any](sliceLike SliceLike[T]) (val T, ok bool)
func GetLastPtr[T any](sliceLike SliceLike[T]) *T
func TryGetLastPtr[T any](sliceLike SliceLike[T]) (val *T, ok bool)
func Set[T any](sliceLike SliceLike[T], idx int, val T)
func TrySet[T any](sliceLike SliceLike[T], idx int, val T) (ok bool)
func SetLast[T any](sliceLike SliceLike[T], val T)
func TrySetLast[T any](sliceLike SliceLike[T], val T) (ok bool) 
func SetFrom[T any](dest SliceLike[T], destIdx int, source SliceLike[T], srcIdx int)
func Swap[T any](sliceLike SliceLike[T], idxA int, idxB int)
func Move[T any](sliceLike SliceLike[T], oldIdx int, newIdx int)
func Copy[T any](dest SliceLike[T], destStart, destLen int, source SliceLike[T], srcStart, srcLen int) (n int)
func IsSorted[T any](sliceLike SliceLike[T], greaterThan func(a *T, b *T) bool) bool
func IsSortedImplicit[T cmp.Ordered](sliceLike SliceLike[T]) bool
func Sort[T any](sliceLike SliceLike[T], greaterThan func(a *T, b T) bool)
func SortImplicit[T cmp.Ordered](sliceLike SliceLike[T])
```

[Back to Top](#go_list_like)
### ListLike[T]
Implement:
```golang
type ListLike[T any] interface {
	SliceLike[T]
	// Increase or decrease the length of the slice/list by `delta` elements,
	// possibly reallocating/resizing the data if needed
	OffsetLen(delta int)
	// Return the total number of values the slice/list can hold
	Cap() int
}
```
To get:
```golang
// SliceLike[T] funcs...
func OffsetLen[T any](listLike ListLike[T], delta int)
func GrowLen[T any](listLike ListLike[T], grow int)
func ShrinkLen[T any](listLike ListLike[T], shrink int)
func Cap[T any](listLike ListLike[T]) int
func GrowCapIfNeeded[T any](listLike ListLike[T], nMoreItems int)
func Append[T any](listLike ListLike[T], vals ...T)
func Clear[T any](listLike ListLike[T]) 
func Insert[T any](listLike ListLike[T], idx int, vals ...T)
func Delete[T any](listLike ListLike[T], idx int, count int)
func DeleteSparse[T any, I Index](listLike ListLike[T], deleteIndexSlice SliceLike[I], sortDeleteIndexes bool)
func Remove[T any](listLike ListLike[T], idx int, count int, outputList ListLike[T])
func RemoveSparse[T any, I Index](listLike ListLike[T], removeIndexSlice SliceLike[I], sortRemoveIndexes bool, outputList ListLike[T])
func Replace[T any](dest ListLike[T], destStart, destLen int, source SliceLike[T], srcStart, srcLen int) (delta int)
func Pop[T any](listLike ListLike[T]) T
```

[Back to Top](#go_list_like)
### QueueLike[T]
Implement:
```golang
type QueueLike[T any] interface {
	ListLike[T]
	// Offset the start location (index/pointer/etc.) of this queue by
	// the given delta. The new 'first' item in the queue should be the item
	// previously located at `queue.GetPtr(0+delta)`.
	OffsetStart(delta int)
}
```
To get:
```golang
// ListLike[T] funcs...
func Dequeue[T any](queueLike QueueLike[T], count int, outputList ListLike[T])
```

[Back to Top](#go_list_like)
### FwdTraversable[T]
Implement:
```golang
type FwdTraversable[T any] interface {
	SliceLike[T]
	// Return the first index in the slice/list
	FirstIdx() (firstIdx int, hasFirst bool)
	// Return the next idx after this one, and whether the next idx is valid/exists
	NextIdx(thisIdx int) (nextIdx int, hasNext bool)
}
```
To get:
```golang
// SliceLike[T] funcs...
func DoActionOnItemsUntilFalse[T any](slice FwdTraversable[T], action func(slice FwdTraversable[T], idx int, item *T) (shouldContinue bool)) (prevIdx int, stopIdx int, stoppedAtEnd bool)
func DoActionOnItemsUntilFalseWithUserdata[T any, U any](slice FwdTraversable[T], action func(slice FwdTraversable[T], idx int, item *T, userdata *U) (shouldContinue bool), userdata *U) (prevIdx int, stopIdx int, stoppedAtEnd bool) 
func DoActionOnAllItems[T any](slice FwdTraversable[T], action func(slice FwdTraversable[T], idx int, item *T)) 
func DoActionOnAllItemsWithUserdata[T any, U any](slice FwdTraversable[T], action func(slice FwdTraversable[T], idx int, item *T, userdata *U), userdata *U)
func FilterIndexes[T any, I Index](slice FwdTraversable[T], selectFunc func(slice FwdTraversable[T], idx I, item *T) bool, outputList ListLike[I])
func FilterIndexesWithUserdata[T any, I Index, U any](slice FwdTraversable[T], selectFunc func(slice FwdTraversable[T], idx I, item *T, userdata *U) bool, outputList ListLike[I], userdata *U)
func MapValues[T any, TT any](slice FwdTraversable[T], mapFunc func(slice FwdTraversable[T], idx int, item *T) TT, outputList ListLike[TT]) 
func MapValuesWithUserdata[T any, TT any, U any](slice FwdTraversable[T], mapFunc func(slice FwdTraversable[T], idx int, item *T, userdata *U) TT, outputList ListLike[TT], userdata *U)
func Accumulate[T any, TT any](slice FwdTraversable[T], initialAccumulation TT, accumulate func(slice FwdTraversable[T], idx int, item *T, currentAccumulation TT) (newAccumulation TT)) (finalAccumulation TT)
func AccumulateWithUserdata[T any, TT any, U any](slice FwdTraversable[T], initialAccumulation TT, accumulate func(slice FwdTraversable[T], idx int, item *T, currentAccumulation TT, userdata *U) (newAccumulation TT), userdata *U) (finalAccumulation TT)
```

[Back to Top](#go_list_like)
### FwdLinkedListLike[T]
No additional functions yet, but included here for future use:
```golang
type FwdLinkedListLike[T any] interface {
	FwdTraversable[T]
	// Set the next idx after this one on the type located at this idx
	SetNextIdx(thisIdx int, nextIdx int)
}
```

[Back to Top](#go_list_like)
### RevTraversable[T]
No additional functions yet, but included here for future use:
```golang
type RevTraversable[T any] interface {
	SliceLike[T]
	// Return the last index in the slice/list
	LastIdx() (lastIdx int, hasLast bool)
	// Return the prev idx before this one, and whether the prev idx is valid/exists
	PrevIdx(thisIdx int) (prevIdx int, hasPrev bool)
}
```

[Back to Top](#go_list_like)
### RevLinkedListLike[T]
No additional functions yet, but included here for future use:
```golang
type RevLinkedListLike[T any] interface {
	RevTraversable[T]
	// Set the prev idx before this one on the type located at this idx
	SetPrevIdx(thisIdx int, prevIdx int)
}
```

[Back to Top](#go_list_like)
## Example
```golang
import (
    "fmt"
    ll "github.com/gabe-lee/go_list_like"
)

func main() {
    mySlice := []byte("Hello World")
	mySliceLike := ll.NewSliceAdapter(&mySlice)
	fmt.Printf("%s\n", mySlice)
	ll.Append(mySliceLike, '!')
	ll.Delete(mySliceLike, 5, 1)
	fmt.Printf("%s\n", mySlice)
}
```
Output:
```
Hello World
HelloWorld!
```
