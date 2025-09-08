package go_list_like

func SortedInsert[T any, IDX Integer, L ListLike[T, IDX]](list L, val T, equalOrder func(a, b T) bool, greaterThan func(a, b T) bool) (insertIdx IDX) {
	if list.PreferLinearOps() {
		insertIdx = sorted_LinearInsert(list, val, equalOrder, greaterThan)
	} else {
		insertIdx = sorted_BinaryInsert(list, val, equalOrder, greaterThan)
	}
	return
}

func SortedInsertIndex[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, val TT, equalOrder func(a T, b TT) bool, greaterThan func(a T, b TT) bool) (insertIdx IDX, append bool) {
	if slice.PreferLinearOps() {
		insertIdx, append = sorted_LinearInsertIndex(slice, val, equalOrder, greaterThan)
	} else {
		insertIdx, append = sorted_BinaryInsertIndex(slice, val, equalOrder, greaterThan)
	}
	return
}

func SortedSearch[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, val TT, equalOrder func(a T, b TT) bool, greaterThan func(a T, b TT) bool) (foundIdx IDX, found bool) {
	if slice.PreferLinearOps() {
		foundIdx, found = sorted_LinearSearch(slice, val, equalOrder, greaterThan)
	} else {
		foundIdx, found = sorted_BinarySearch(slice, val, equalOrder, greaterThan)
	}
	return
}

func sorted_BinaryInsert[T any, IDX Integer, L ListLike[T, IDX]](list L, val T, equalOrder func(a, b T) bool, greaterThan func(a, b T) bool) (insertIdx IDX) {
	var append bool
	insertIdx, append = sorted_BinaryInsertIndex(list, val, equalOrder, greaterThan)
	if append {
		insertIdx, _ = AppendVar(list, val)
	} else {
		insertIdx, _ = InsertVar(list, insertIdx, val)
	}
	return
}

func sorted_LinearInsert[T any, IDX Integer, L ListLike[T, IDX]](list L, val T, equalOrder func(a, b T) bool, greaterThan func(a, b T) bool) (insertIdx IDX) {
	var append bool
	insertIdx, append = sorted_LinearInsertIndex(list, val, equalOrder, greaterThan)
	if append {
		insertIdx, _ = AppendVar(list, val)
	} else {
		insertIdx, _ = InsertVar(list, insertIdx, val)
	}
	return
}

// func Sorted_BinaryInsertWithMaps[T any, IDX Integer, L ListLike[T, IDX], LL ListLike[IDX, IDX]](list L, val T, equalOrder func(a, b T) bool, greaterThan func(a, b T) bool, indexMaps []IndexMap[T, IDX, L, LL]) (insertIdx IDX) {
// 	insertIdx = sorted_BinaryInsert(list, val, equalOrder, greaterThan)
// 	for _, m := range indexMaps {
// 		m.MoveIndexesUp(list, insertIdx)
// 		newMapIdx, _ := AppendVar(m.Indexes, insertIdx)
// 		m.ResortAltered(list, newMapIdx)
// 	}
// 	return
// }

func sorted_BinaryInsertIndex[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, val TT, equalOrder func(a T, b TT) bool, greaterThan func(a T, b TT) bool) (insertIdx IDX, append bool) {
	var lo IDX = slice.FirstIdx()
	var hi IDX = slice.LastIdx()
	var ok bool = slice.IdxValid(lo) && slice.IdxValid(hi)
	if !ok {
		return
	}
	var exitHi, found bool
	insertIdx, found, exitHi, _ = sorted_BinaryLocate(slice, lo, hi, val, equalOrder, greaterThan)

	if !found && exitHi {
		append = true
	}
	return
}

func sorted_LinearInsertIndex[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, val TT, equalOrder func(a T, b TT) bool, greaterThan func(a T, b TT) bool) (insertIdx IDX, append bool) {
	var lo IDX = slice.FirstIdx()
	var hi IDX = slice.LastIdx()
	var ok bool = slice.IdxValid(lo) && slice.IdxValid(hi)
	if !ok {
		return
	}
	var exitHi, found bool
	insertIdx, found, exitHi, _ = sorted_LinearLocate(slice, lo, hi, val, equalOrder, greaterThan)
	if !found && exitHi {
		append = true
	}
	return
}

func sorted_BinarySearch[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, val TT, equalValue func(a T, b TT) bool, greaterThan func(a T, b TT) bool) (idx IDX, found bool) {
	var lo IDX = slice.FirstIdx()
	var hi IDX = slice.LastIdx()
	var ok bool = slice.IdxValid(lo) && slice.IdxValid(hi)
	if !ok {
		return
	}
	idx, found, _, _ = sorted_BinaryLocate(slice, lo, hi, val, equalValue, greaterThan)
	return
}

func sorted_LinearSearch[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, val TT, equalValue func(a T, b TT) bool, greaterThan func(a T, b TT) bool) (idx IDX, found bool) {
	var lo IDX = slice.FirstIdx()
	var hi IDX = slice.LastIdx()
	var ok bool = slice.IdxValid(lo) && slice.IdxValid(hi)
	if !ok {
		return
	}
	idx, found, _, _ = sorted_LinearLocate(slice, lo, hi, val, equalValue, greaterThan)
	return
}

func sorted_BinaryLocate[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, lo, hi IDX, locateVal TT, equalValue func(a T, b TT) bool, greaterThan func(a T, b TT) bool) (idx IDX, found bool, exitHi bool, exitLo bool) {
	originalHi := hi
	originalLo := lo
	var val T
	for {
		idx = slice.SplitRange(lo, hi)
		val = slice.Get(idx)
		if equalValue(val, locateVal) {
			found = true
			return
		}
		if greaterThan(val, locateVal) {
			if idx == lo {
				exitLo = idx == originalLo
				return
			}
			hi = slice.PrevIdx(idx)
		} else {
			if idx == hi {
				exitHi = idx == originalHi
				if !exitHi {
					idx = slice.NextIdx(hi)
				}
				return
			}
			lo = slice.NextIdx(idx)
		}
	}
}

func sorted_LinearLocate[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, lo, hi IDX, locateVal TT, equalValue func(a T, b TT) bool, greaterThan func(a T, b TT) bool) (idx IDX, found bool, exitHi bool, exitLo bool) {
	var val T
	idx = lo
	for {
		val = slice.Get(idx)
		if equalValue(val, locateVal) {
			found = true
			return
		}
		if greaterThan(val, locateVal) {
			return
		} else {
			if idx == hi {
				exitHi = true
				return
			}
			idx = slice.NextIdx(idx)
		}
	}
}

// func sorted_BinaryFindMappedIdx[T any, IDX Integer, S SliceLike[T, IDX], L ListLike[IDX, IDX]](slice S, valIdx IDX, indexMap IndexMap[T, IDX, S, L]) (mapIdx IDX, found bool) {
// 	var lo IDX = indexMap.Indexes.FirstIdx()
// 	var hi IDX = indexMap.Indexes.LastIdx()
// 	var midMapIdx IDX = lo
// 	var midVal T
// 	val := slice.Get(valIdx)
// 	if !indexMap.Indexes.RangeValid(lo, hi) {
// 		return 0, false
// 	}
// 	for {
// 		if hi == lo {
// 			midVal, _ = indexMap.GetVal(slice, hi)
// 			if !indexMap.GreaterThan(midVal, val) && !indexMap.GreaterThan(val, midVal) {
// 				return hi, true
// 			}
// 			return 0, false
// 		}
// 		midMapIdx = indexMap.Indexes.SplitRange(lo, hi)
// 		midVal, _ = indexMap.GetVal(slice, midMapIdx)
// 		if indexMap.GreaterThan(midVal, val) {
// 			hi = midMapIdx
// 		} else if indexMap.GreaterThan(val, midVal) {
// 			lo = indexMap.Indexes.NextIdx(midMapIdx)
// 		} else {
// 			return midMapIdx, true
// 		}
// 	}
// }

func SortedSetAndResort[T any, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T, greaterThan func(a, b T) bool) (newIdx IDX) {
	slice.Set(idx, val)
	newIdx = idx
	var testVal T
	testIdx := slice.NextIdx(idx)
	ok := slice.IdxValid(testIdx)
	if ok {
		testVal = slice.Get(testIdx)
		for ok && greaterThan(val, testVal) {
			newIdx = testIdx
			testIdx = slice.NextIdx(testIdx)
			ok = slice.IdxValid(testIdx)
			if ok {
				testVal = slice.Get(testIdx)
			}
		}
		if newIdx != idx {
			slice.Move(idx, newIdx)
			return
		}
	}
	testIdx = slice.PrevIdx(idx)
	ok = slice.IdxValid(testIdx)
	if ok {
		testVal = slice.Get(testIdx)
		for ok && greaterThan(testVal, val) {
			newIdx = testIdx
			testIdx = slice.PrevIdx(testIdx)
			ok = slice.IdxValid(testIdx)
			if ok {
				testVal = slice.Get(testIdx)
			}
		}
		if newIdx != idx {
			slice.Move(idx, newIdx)
			return
		}
	}
	return idx
}

func GreaterThanImplicit[T Ordered](a, b T) bool {
	return a > b
}
func LessThanImplicit[T Ordered](a, b T) bool {
	return a < b
}
func EqualImplicit[T Equatable](a, b T) bool {
	return a == b
}
