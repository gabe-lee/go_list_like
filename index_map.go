package go_list_like

// type SortedIndexMap[T any, IDX Integer, L ListLike[T, IDX], LL ListLike[IDX, IDX]] struct {
// 	GreaterThan func(a, b T) bool
// 	Indexes     LL
// }

// func NewIndexMapsVar[T any, IDX Integer, L ListLike[T, IDX], LL ListLike[IDX, IDX], M SortedIndexMap[T, IDX, L, LL]](maps ...M) SliceAdapter[M] {
// 	return NewSliceAdapter(maps)
// }

// func NewIndexMap[T any, IDX Integer, L ListLike[T, IDX], LL ListLike[IDX, IDX]](indexes LL, greaterThan func(a, b T) bool) SortedIndexMap[T, IDX, L, LL] {
// 	return SortedIndexMap[T, IDX, L, LL]{
// 		GreaterThan: greaterThan,
// 		Indexes:     indexes,
// 	}
// }

// func (m SortedIndexMap[T, IDX, L, LL]) rearrangeIndexesNoValueChange(list L, firstMoveDown, lastMoveDown, moveDownCount IDX, hasMoveDowns bool, firstMoveUp, lastMoveUp, moveUpCount IDX, hasMoveUps bool, firstDelete, lastDelete IDX, hasDeletes bool) {
// 	var idx IDX
// 	var thisIdxDone, firstDone, nextIsFirst bool
// 	mIdx := m.Indexes.FirstIdx()
// 	ok := m.Indexes.IdxValid(mIdx)
// 	for ok {
// 		idx = m.Indexes.Get(mIdx)
// 		thisIdxDone = false
// 		if list.ConsecutiveIndexesInOrder() {
// 			if hasMoveDowns && idx >= firstMoveDown && idx <= lastMoveDown {
// 				idx -= moveDownCount
// 				m.Indexes.Set(mIdx, idx)
// 			} else if hasMoveUps && idx >= firstMoveUp && idx <= lastMoveUp {
// 				idx += moveUpCount
// 				m.Indexes.Set(mIdx, idx)
// 			} else if hasDeletes && idx >= firstDelete && idx <= lastDelete {
// 				dIdx := mIdx
// 				if !firstDone {
// 					nextIsFirst = true
// 				} else {
// 					mIdx = m.Indexes.PrevIdx(mIdx)
// 				}
// 				DeleteRange(m.Indexes, dIdx, dIdx)
// 			}
// 		} else {
// 			if hasMoveDowns {
// 				DoActionOnItemsInRangeUntilFalse(list, firstMoveDown, lastMoveDown, func(list L, downIdx IDX, item T) bool {
// 					if idx == downIdx {
// 						idx -= moveDownCount
// 						m.Indexes.Set(mIdx, idx)
// 						thisIdxDone = true
// 						if downIdx == firstMoveDown {
// 							if downIdx == lastMoveDown {
// 								hasMoveDowns = false
// 							} else {
// 								firstMoveDown = list.NextIdx(firstMoveDown)
// 							}
// 						} else if downIdx == lastMoveDown {
// 							lastMoveDown = list.PrevIdx(lastMoveDown)
// 						}
// 						return false
// 					}
// 					return true
// 				})
// 			}
// 			if hasMoveUps && !thisIdxDone {
// 				DoActionOnItemsInRangeUntilFalse(list, firstMoveUp, lastMoveUp, func(list L, upIdx IDX, item T) bool {
// 					if idx == upIdx {
// 						idx += moveUpCount
// 						m.Indexes.Set(mIdx, idx)
// 						thisIdxDone = true
// 						if upIdx == firstMoveUp {
// 							if upIdx == lastMoveUp {
// 								hasMoveUps = false
// 							} else {
// 								firstMoveUp = list.NextIdx(firstMoveUp)
// 							}
// 						} else if upIdx == lastMoveUp {
// 							lastMoveUp = list.PrevIdx(lastMoveUp)
// 						}
// 						return false
// 					}
// 					return true
// 				})
// 			}
// 			if hasDeletes && !thisIdxDone {
// 				DoActionOnItemsInRangeUntilFalse(list, firstDelete, lastDelete, func(list L, delIdx IDX, item T) bool {
// 					if idx == delIdx {
// 						dIdx := mIdx
// 						if !firstDone {
// 							nextIsFirst = true
// 						} else {
// 							mIdx = m.Indexes.PrevIdx(mIdx)
// 						}
// 						DeleteRange(m.Indexes, dIdx, dIdx)
// 						thisIdxDone = true
// 						if delIdx == firstDelete {
// 							if delIdx == lastDelete {
// 								hasDeletes = false
// 							} else {
// 								firstDelete = list.NextIdx(firstDelete)
// 							}
// 						} else if delIdx == lastDelete {
// 							lastDelete = list.PrevIdx(lastDelete)
// 						}
// 						return false
// 					}
// 					return true
// 				})
// 			}
// 		}
// 		ok = mIdx != m.Indexes.LastIdx()
// 		if nextIsFirst {
// 			mIdx = m.Indexes.FirstIdx()
// 			nextIsFirst = false
// 		} else {
// 			mIdx = m.Indexes.NextIdx(mIdx)
// 			firstDone = true
// 		}
// 		ok = ok && m.Indexes.IdxValid(mIdx)
// 	}
// }

// // func (m IndexMap[T, IDX, S, L]) ShiftIndexUp(slice S, oldMovedIndex, newMovedIndex IDX) (mappedIdx IDX, found bool) {
// // 	var maxCount IDX = slice.LenBetween(oldMovedIndex, newMovedIndex)
// // 	var count IDX = 0
// // 	found = false
// // 	i := m.Indexes.FirstIdx()
// // 	ok := m.Indexes.IdxValid(i)
// // 	for ok {
// // 		v := m.Indexes.Get(i)
// // 		if v == oldMovedIndex {
// // 			found = true
// // 			mappedIdx = i
// // 			if count >= maxCount && found {
// // 				break
// // 			}
// // 		} else if v > oldMovedIndex && v <= newMovedIndex {
// // 			SetSubtract(m.Indexes, i, 1)
// // 			count += 1
// // 			if count >= maxCount && found {
// // 				break
// // 			}
// // 		}
// // 		i = m.Indexes.NextIdx(i)
// // 		ok = m.Indexes.IdxValid(i)
// // 	}
// // 	if found {
// // 		m.Indexes.Set(mappedIdx, newMovedIndex)
// // 	}
// // 	return
// // }

// // func (m IndexMap[T, IDX, S, L]) ShiftIndexDown(slice S, oldMovedIndex, newMovedIndex IDX) (mappedIdx IDX, found bool) {
// // 	var maxCount IDX = slice.LenBetween(newMovedIndex, oldMovedIndex)
// // 	var count IDX = 0
// // 	found = false
// // 	i := m.Indexes.FirstIdx()
// // 	ok := m.Indexes.IdxValid(i)
// // 	for ok {
// // 		v := m.Indexes.Get(i)
// // 		if v == oldMovedIndex {
// // 			found = true
// // 			mappedIdx = i
// // 			if count >= maxCount && found {
// // 				break
// // 			}
// // 		} else if v < oldMovedIndex && v >= newMovedIndex {
// // 			SetSubtract(m.Indexes, i, 1)
// // 			count += 1
// // 			if count >= maxCount && found {
// // 				break
// // 			}
// // 		}
// // 		i = m.Indexes.NextIdx(i)
// // 		ok = m.Indexes.IdxValid(i)
// // 	}
// // 	if found {
// // 		m.Indexes.Set(mappedIdx, newMovedIndex)
// // 	}
// // 	return
// // }

// // func (m ListIndexMap[T, IDX, L, LL]) DeleteIndexesAndMoveOthersDown(list L, Sli) {
// // 	var maxCount IDX = list.LenBetween(deleteIdx, list.LastIdx())
// // 	var count IDX = 0
// // 	i := m.Indexes.FirstIdx()
// // 	ok := m.Indexes.IdxValid(i)
// // 	for ok {
// // 		v := m.Indexes.Get(i)
// // 		if v >= startIdx {
// // 			SetSubtract(m.Indexes, i, 1)
// // 			count += 1
// // 			if count >= maxCount {
// // 				break
// // 			}
// // 		}
// // 		i = m.Indexes.NextIdx(i)
// // 		ok = m.Indexes.IdxValid(i)
// // 	}
// // }

// // func (m IndexMap[T, IDX, S, L]) MoveIndexesUp(slice S, startIdx IDX) {
// // 	var maxCount IDX = slice.Len() - startIdx
// // 	var count IDX = 0
// // 	i := m.Indexes.FirstIdx()
// // 	ok := m.Indexes.IdxValid(i)
// // 	for ok {
// // 		v := m.Indexes.Get(i)
// // 		if v >= startIdx {
// // 			SetAdd(m.Indexes, i, 1)
// // 			count += 1
// // 			if count >= maxCount {
// // 				break
// // 			}
// // 		}
// // 		i = m.Indexes.NextIdx(i)
// // 		ok = m.Indexes.IdxValid(i)
// // 	}
// // }

// // func (m IndexMap[T, IDX, S, L]) MoveIndexesDown(slice S, startIdx IDX) {
// // 	var maxCount IDX = slice.Len() - startIdx
// // 	var count IDX = 0
// // 	i := m.Indexes.FirstIdx()
// // 	ok := m.Indexes.IdxValid(i)
// // 	for ok {
// // 		v := m.Indexes.Get(i)
// // 		if v >= startIdx {
// // 			SetSubtract(m.Indexes, i, 1)
// // 			count += 1
// // 			if count >= maxCount {
// // 				break
// // 			}
// // 		}
// // 		i = m.Indexes.NextIdx(i)
// // 		ok = m.Indexes.IdxValid(i)
// // 	}
// // }

// // func (m IndexMap[T, IDX, S, L]) TryResortAltered(slice S, mapIdx IDX) (ok bool) {
// // 	var v, vv T
// // 	var i, ii IDX
// // 	var j IDX
// // 	var sorted bool
// // 	i = mapIdx
// // 	v, j, ok = m.TryGetVal(slice, i)
// // 	if !ok {
// // 		return
// // 	}
// // 	ii = m.Indexes.PrevIdx(i)
// // 	vv, _, ok = m.TryGetVal(slice, ii)
// // 	for ok && m.GreaterThan(vv, v) {
// // 		sorted = true
// // 		Overwrite(m.Indexes, ii, i)
// // 		i = ii
// // 		ii = m.Indexes.PrevIdx(ii)
// // 		vv, _, ok = m.TryGetVal(slice, ii)
// // 	}
// // 	if !sorted {
// // 		i = mapIdx
// // 		ii = m.Indexes.NextIdx(i)
// // 		vv, _, ok = m.TryGetVal(slice, ii)
// // 		for ok && m.GreaterThan(v, vv) {
// // 			sorted = true
// // 			Overwrite(m.Indexes, ii, i)
// // 			i = ii
// // 			ii = m.Indexes.NextIdx(ii)
// // 			vv, _, ok = m.TryGetVal(slice, ii)
// // 		}
// // 	}
// // 	m.Indexes.Set(i, j)
// // 	ok = true
// // 	return
// // }

// func (m SortedIndexMap[T, IDX, L, LL]) resortAltered(list L, mapIdx IDX) {
// 	var v, vv T
// 	var i, ii IDX
// 	var j IDX
// 	var ok bool
// 	var sorted bool
// 	i = mapIdx
// 	v, j = m.GetVal(list, i)
// 	ii = m.Indexes.PrevIdx(i)
// 	ok = m.Indexes.IdxValid(ii)
// 	if ok {
// 		vv, _ = m.GetVal(list, ii)
// 		for ok && m.GreaterThan(vv, v) {
// 			sorted = true
// 			Overwrite(m.Indexes, ii, i)
// 			i = ii
// 			ii = m.Indexes.PrevIdx(ii)
// 			ok = m.Indexes.IdxValid(ii)
// 			if ok {
// 				vv, _ = m.GetVal(list, ii)
// 			}
// 		}
// 	}
// 	if !sorted {
// 		i = mapIdx
// 		ii = m.Indexes.NextIdx(i)
// 		ok = m.Indexes.IdxValid(ii)
// 		if ok {
// 			vv, _ = m.GetVal(list, ii)
// 			for ok && m.GreaterThan(v, vv) {
// 				sorted = true
// 				Overwrite(m.Indexes, ii, i)
// 				i = ii
// 				ii = m.Indexes.NextIdx(ii)
// 				ok = m.Indexes.IdxValid(ii)
// 				if ok {
// 					vv, _ = m.GetVal(list, ii)
// 				}
// 			}
// 		}
// 	}
// 	m.Indexes.Set(i, j)
// }

// func (m SortedIndexMap[T, IDX, S, L]) TryGetVal(slice S, mapIdx IDX) (val T, idx IDX, ok bool) {
// 	ok = m.Indexes.IdxValid(mapIdx)
// 	if !ok {
// 		return
// 	}
// 	idx = m.Indexes.Get(mapIdx)
// 	ok = slice.IdxValid(idx)
// 	if !ok {
// 		return
// 	}
// 	val = slice.Get(idx)
// 	return
// }
// func (m SortedIndexMap[T, IDX, S, L]) GetVal(slice S, mapIdx IDX) (val T, idx IDX) {
// 	idx = m.Indexes.Get(mapIdx)
// 	val = slice.Get(idx)
// 	return
// }
