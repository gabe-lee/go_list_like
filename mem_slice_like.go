package go_list_like

type MemSliceLike[T any, IDX Integer] interface {
	SliceLike[T, IDX]
	// Get the a pointer to the value at the provided index
	GetPtr(idx IDX) (ptr *T)
}

func GetPtr[T any, IDX Integer, S MemSliceLike[T, IDX]](memSliceLike S, idx IDX) (ptr *T) {
	ptr = memSliceLike.GetPtr(idx)
	return
}
func TryGetPtr[T any, IDX Integer, S MemSliceLike[T, IDX]](memSliceLike S, idx IDX) (ptr *T, ok bool) {
	ok = memSliceLike.IdxValid(idx)
	if !ok {
		return
	}
	ptr = memSliceLike.GetPtr(idx)
	return
}

type MemListLike[T any, IDX Integer] interface {
	MemSliceLike[T, IDX]
	ListLike[T, IDX]
}

type MemQueueLike[T any, IDX Integer] interface {
	MemSliceLike[T, IDX]
	QueueLike[T, IDX]
}
