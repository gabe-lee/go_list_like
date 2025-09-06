package go_list_like

type GoSliceLike[T any] interface {
	// Return the underlying golang slice that holds the data
	GoSlice() []T
}
