package go_list_like

import (
	"slices"
	"testing"
)

func Fuzz_InsertionSort_(f *testing.F) {
	var nilSlice []byte
	var emptySlice []byte = make([]byte, 0, 10)
	f.Add(nilSlice)
	f.Add(emptySlice)
	f.Add([]byte{0, 1, 2, 3, 4})
	f.Add([]byte{4, 3, 2, 1, 0})
	f.Add([]byte{1, 9, 37, 223})
	f.Add([]byte{223, 37, 9, 1})
	f.Add([]byte{56, 42, 3, 77, 22, 5, 109})
	f.Fuzz(func(t *testing.T, a []byte) {
		oldLen := len(a)
		var oldSum uint64 = 0
		for _, b := range a {
			oldSum += uint64(b) + 1
		}
		aa := NewSliceAdapterIndirect(&a)
		InsertionSortImplicit(aa)
		newLen := len(a)
		if oldLen != newLen {
			t.Errorf("\ntest case failed: len mismatch\nSLICE: %v\nEXP LEN: %d\nGOT LEN: %d\n", a, oldLen, newLen)
		}
		var newSum uint64 = 0
		for _, b := range a {
			newSum += uint64(b) + 1
		}
		if oldLen != newLen {
			t.Errorf("\ntest case failed: sum mismatch\nSLICE: %v\nEXP SUM: %d\nGOT SUM: %d\n", a, oldSum, newSum)
		}
		var ii int = 1
		var i int = 0
		for ii < oldLen {
			if a[i] > a[ii] {
				t.Errorf("\ntest case failed: not sorted\nSLICE: %v\nVAL : %d > %d\nIDX: %d < %d\n", a, a[i], a[ii], i, ii)
			}
			i += 1
			ii += 1
		}
	})
}

func Fuzz_SortedInsert_(f *testing.F) {
	var nilSlice []byte
	var emptySlice []byte = make([]byte, 0, 10)
	f.Add(nilSlice, byte(5))
	f.Add(emptySlice, byte(5))
	f.Add([]byte{0, 1, 2, 3, 4}, byte(5))
	f.Add([]byte{1, 2, 3, 4}, byte(5))
	f.Add([]byte{0, 1, 2, 3, 4, 6, 7, 8, 9, 10}, byte(5))
	f.Add([]byte{0, 1, 2, 3, 4, 6, 7, 8, 10}, byte(5))
	f.Add([]byte{5, 5, 5, 5, 6}, byte(5))
	f.Add([]byte{5, 5, 5, 6}, byte(5))
	f.Add([]byte{6, 6, 6, 6}, byte(5))
	f.Add([]byte{6, 6, 6, 6, 6}, byte(5))
	f.Fuzz(func(t *testing.T, a []byte, b byte) {
		slices.Sort(a)
		var sum uint64 = 0
		for _, b := range a {
			sum += uint64(b) + 1
		}
		var expSum = sum + uint64(b) + 1
		var expLen = len(a) + 1
		aa := NewSliceAdapter(a)
		SortedInsert(&aa, b, EqualImplicit, GreaterThanImplicit)
		var gotSum uint64 = 0
		var gotLen = aa.Len()
		if gotLen != expLen {
			t.Errorf("\ntest case failed: len mismatch\nINSERT VAL: %d\nOLD SLICE: %v\nNEW SLICE: %v\nEXP LEN: %d\nGOT LEN: %d\n", b, a, aa, expLen, gotLen)
		}
		for i := range aa.Len() {
			b := Get(&aa, i)
			gotSum += uint64(b) + 1
		}
		if gotSum != expSum {
			t.Errorf("\ntest case failed: sum mismatch\nINSERT VAL: %d\nOLD SLICE: %v\nNEW SLICE: %v\nEXP SUM: %d\nGOT SUM: %d\n", b, a, aa, expSum, gotSum)
		}
		var ii int = 1
		var i int = 0
		for ii < gotLen {
			if Get(&aa, i) > Get(&aa, ii) {
				t.Errorf("\ntest case failed: not sorted\nINSERT VAL: %d\nOLD SLICE: %v\nNEW SLICE: %v\nVAL : %d > %d\nIDX: %d < %d\n", b, a, aa, Get(&aa, i), Get(&aa, ii), i, ii)
			}
			i += 1
			ii += 1
		}
	})
}

func Fuzz_SortedSearch_(f *testing.F) {
	var nilSlice []byte
	var emptySlice []byte = make([]byte, 0, 10)
	f.Add(nilSlice, byte(5))
	f.Add(emptySlice, byte(5))
	f.Add([]byte{0, 1, 2, 3, 4}, byte(5))
	f.Add([]byte{1, 2, 3, 4}, byte(5))
	f.Add([]byte{0, 1, 2, 3, 4, 6, 7, 8, 9, 10}, byte(5))
	f.Add([]byte{0, 1, 2, 3, 4, 6, 7, 8, 10}, byte(5))
	f.Add([]byte{5, 5, 5, 5, 6}, byte(5))
	f.Add([]byte{5, 5, 5, 6}, byte(5))
	f.Add([]byte{6, 6, 6, 6}, byte(5))
	f.Add([]byte{6, 6, 6, 6, 6}, byte(5))
	f.Fuzz(func(t *testing.T, a []byte, b byte) {
		slices.Sort(a)
		var minValidIdx int
		var maxValidIdx int
		var existsInList bool
		for i, bb := range a {
			if !existsInList {
				if bb == b {
					existsInList = true
					minValidIdx = i
					maxValidIdx = i
				}
			} else {
				if bb > b {
					maxValidIdx = i - 1
					break
				} else {
					maxValidIdx = i
				}
			}
		}
		aa := NewSliceAdapter(a)
		foundIdx, found := SortedSearch(&aa, b, EqualImplicit, GreaterThanImplicit)
		if found && !existsInList {
			t.Errorf("\ntest case failed: value does not exist in list but was 'found' by BinarySearch\nSEARCH VAL: %d\nSLICE: %v\nBAD FOUND IDX: %d\n", b, a, foundIdx)
		}
		if !found && existsInList {
			t.Errorf("\ntest case failed: value exists in list but was not 'found' by BinarySearch\nSEARCH VAL: %d\nSLICE: %v\n", b, a)
		}
		if found && (foundIdx < minValidIdx || foundIdx > maxValidIdx) {
			t.Errorf("\ntest case failed: found value idx outside valid range of matching bytes\nSEARCH VAL: %d\nSLICE: %v\nMIN VALID IDX: %d\nMAX VALID IDX: %d\n'FOUND' IDX: %d", b, a, minValidIdx, maxValidIdx, foundIdx)
		}
	})
}

// func Fuzz_SortedSetAndResort(f *testing.F) {
// 	f.Add([]byte{0, 1, 2, 3, 4}, byte(5), int(1))
// 	f.Add([]byte{1, 2, 3, 4}, byte(5), int(2))
// 	f.Add([]byte{0, 1, 2, 3, 4, 6, 7, 8, 9, 10}, byte(5), int(3))
// 	f.Add([]byte{0, 1, 2, 3, 4, 6, 7, 8, 10}, byte(5), int(4))
// 	f.Add([]byte{5, 5, 5, 5, 6}, byte(5), int(5))
// 	f.Add([]byte{5, 5, 5, 6}, byte(5), int(6))
// 	f.Add([]byte{6, 6, 6, 6}, byte(5), int(7))
// 	f.Add([]byte{6, 6, 6, 6, 6}, byte(5), int(8))
// 	f.Fuzz(func(t *testing.T, goSlice []byte, setVal byte, setIdx int) {
// 		if len(goSlice) == 0 {
// 			return
// 		}
// 		slices.Sort(goSlice)
// 		llSliceData := slices.Clone(goSlice)
// 		llSlice := NewSliceAdapter(llSliceData)
// 		setIdx = setIdx % len(goSlice)
// 		// go slice set and resort
// 		goSlice[setIdx] = setVal
// 		slices.Sort(goSlice)
// 		// ll slice set and resort
// 		SortedSetAndResort(&llSlice, setIdx, setVal, GreaterThanImplicit)

// 		if found && !existsInList {
// 			t.Errorf("\ntest case failed: value does not exist in list but was 'found' by BinarySearch\nSEARCH VAL: %d\nSLICE: %v\nBAD FOUND IDX: %d\n", setVal, goSlice, foundIdx)
// 		}
// 		if !found && existsInList {
// 			t.Errorf("\ntest case failed: value exists in list but was not 'found' by BinarySearch\nSEARCH VAL: %d\nSLICE: %v\n", setVal, goSlice)
// 		}
// 		if found && (foundIdx < minValidIdx || foundIdx > maxValidIdx) {
// 			t.Errorf("\ntest case failed: found value idx outside valid range of matching bytes\nSEARCH VAL: %d\nSLICE: %v\nMIN VALID IDX: %d\nMAX VALID IDX: %d\n'FOUND' IDX: %d", setVal, goSlice, minValidIdx, maxValidIdx, foundIdx)
// 		}
// 	})
// }
