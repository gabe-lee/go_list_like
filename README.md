# go_slice_like
Interface for structs that can operate lise a slice or list

## Why?
This package provides a simple interface for things that can behave like a standard Golang slice, but may be user-defined data structures. It provides a simple wrapper around Golang slices themselves to adapt them automatically.

Implement:
```golang
type SliceLike[T any] interface {
	GetPtr(idx int) *T
	Len() int
}
```
To get:
```golang
func Get[T any](sliceLike SliceLike[T], idx int) T
func GetPtr[T any](sliceLike SliceLike[T], idx int) *T
func Set[T any](sliceLike SliceLike[T], idx int, val T)
func Len[T any](sliceLike SliceLike[T]) int
```
And implement:
```golang
type ListLike[T any] interface {
	SliceLike[T]
	AddLen(delta int)
}
```
To get:
```golang
func Insert[T any](listLike ListLike[T], idx int, vals ...T)
func Delete[T any](listLike ListLike[T], idx int, count int)
```

## Installation
Run this command from your project directory
```
go get github.com/gabe-lee/go_slice_like@latest
```

## Example
```golang
import (
    "fmt"
    sl "github.com/gabe-lee/go_slice_like"
)

func main() {
    mySlice := []byte("Hello World")
	mySliceLike := sl.New(&mySlice)
	fmt.Printf("%s\n", mySlice)
	sl.Append(&mySliceLike, '!')
	sl.Delete(&mySliceLike, 5, 1)
	fmt.Printf("%s\n", mySlice)
}
```
Output:
```
Hello World
HelloWorld!
```