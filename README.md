# go_list_like
Interface for structs that can operate lise a slice or list

## Why?
This package provides a simple interface for things that can behave like a 'List' or 'Vector', but may be user-defined data structures with complex memory layouts. It provides a simple wrapper type for Golang slices themselves to adapt them automatically.

Implement:
```golang
type SliceLike[T any] interface {
	GetPtr(idx int) *T
	Len() int
}
```
To get:
```golang
func Len[T any](sliceLike SliceLike[T]) int
func Get[T any](sliceLike SliceLike[T], idx int) T
func GetPtr[T any](sliceLike SliceLike[T], idx int) *T
func GetLast[T any](sliceLike SliceLike[T]) T 
func GetLastPtr[T any](sliceLike SliceLike[T]) *T
func Set[T any](sliceLike SliceLike[T], idx int, val T)
func SetLast[T any](sliceLike SliceLike[T], val T)
func SetFrom[T any](dest SliceLike[T], destIdx int, source SliceLike[T], srcIdx int)
func Swap[T any](sliceLike SliceLike[T], idxA int, idxB int)
func Move[T any](sliceLike SliceLike[T], oldIdx int, newIdx int)
func Copy[T any](dest SliceLike[T], destStart, destLen int, source SliceLike[T], srcStart, srcLen int) (n int)
```
And implement:
```golang
type ListLike[T any] interface {
	SliceLike[T]
	OffsetLen(delta int)
    Cap() int
}
```
To get:
```golang
// SliceLike[T] funcs...
func Cap[T any](listLike ListLike[T]) int
func Append[T any](listLike ListLike[T], vals ...T)
func Insert[T any](listLike ListLike[T], idx int, vals ...T)
func Delete[T any](listLike ListLike[T], idx int, count int)
func Remove[T any](listLike ListLike[T], idx int, count int) []T
func Replace[T any](dest ListLike[T], destStart, destLen int, source SliceLike[T], srcStart, srcLen int) (delta int)
func Pop[T any](listLike ListLike[T]) T
func GrowLen[T any](listLike ListLike[T], grow int)
func ShrinkLen[T any](listLike ListLike[T], shrink int)
func GrowCapIfNeeded[T any](listLike ListLike[T], nMoreItems int)
```
## Installation
Run this command from your project directory
```
go get github.com/gabe-lee/go_list_like@latest
```

## Example
```golang
import (
    "fmt"
    ll "github.com/gabe-lee/go_list_like"
)

func main() {
    mySlice := []byte("Hello World")
	mySliceLike := ll.New(&mySlice)
	fmt.Printf("%s\n", mySlice)
	ll.Append(&mySliceLike, '!')
	ll.Delete(&mySliceLike, 5, 1)
	fmt.Printf("%s\n", mySlice)
}
```
Output:
```
Hello World
HelloWorld!
```