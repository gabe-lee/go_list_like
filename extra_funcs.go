package go_list_like

import (
	"math"
	"unsafe"
)

type Integer interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

type Equatable interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64 | ~unsafe.Pointer | ~bool | ~string
}

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// slice[idx] = slice[idx] + val
func SetAdd[T Number, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)+val)
}

// return slice[idx] + val
func GetAdd[T Number, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) + val
}

// slice[idx] = slice[idx] - val
func SetSubtract[T Number, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)-val)
}

// return slice[idx] - val
func GetSubtract[T Number, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) - val
}

// slice[idx] = slice[idx] * val
func SetMultiply[T Number, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)*val)
}

// return slice[idx] * val
func GetMultiply[T Number, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) * val
}

// slice[idx] = slice[idx] / val
func SetDivide[T Number, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)/val)
}

// return slice[idx] / val
func GetDivide[T Number, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) / val
}

// slice[idx] = slice[idx] % val
func SetModulo[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)%val)
}

// return slice[idx] % val
func GetModulo[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) % val
}

// return slice[idx] % val, slice[idx] - (slice[idx] % val)
func GetModRem[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (mod T, rem T) {
	v := slice.Get(idx)
	mod = v % val
	return mod, v - mod
}

// slice[idx] = math.Mod(slice[idx], val)
func SetFModulo[T Float, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, T(math.Mod(float64(slice.Get(idx)), float64(val))))
}

// return math.Mod(slice[idx], val)
func GetFModulo[T Float, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return T(math.Mod(float64(slice.Get(idx)), float64(val)))
}

// return math.Mod(slice[idx], val), slice[idx] - math.Mod(slice[idx], val)
func GetFModRem[T Float, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (mod T, rem T) {
	v := slice.Get(idx)
	mod = T(math.Mod(float64(slice.Get(idx)), float64(val)))
	return mod, v - mod
}

// slice[idx] = slice[idx] & val
func SetBitAnd[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)&val)
}

// return slice[idx] & val
func GetBitAnd[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) & val
}

// slice[idx] = slice[idx] & val
func SetBitOr[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)|val)
}

// return slice[idx] & val
func GetBitOr[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) | val
}

// slice[idx] = slice[idx] ^ val
func SetBitXor[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)^val)
}

// return slice[idx] ^ val
func GetBitXor[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) ^ val
}

// slice[idx] = ^slice[idx]
func SetBitInvert[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX) {
	slice.Set(idx, ^slice.Get(idx))
}

// return ^slice[idx]
func GetBitInvert[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX) (result T) {
	return ^slice.Get(idx)
}

// slice[idx] = slice[idx] << val
func SetBitLsh[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)<<val)
}

// return slice[idx] << val
func GetBitLsh[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) << val
}

// slice[idx] = slice[idx] >> val
func SetBitRsh[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) {
	slice.Set(idx, slice.Get(idx)>>val)
}

// return slice[idx] >> val
func GetBitRsh[T Integer, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (result T) {
	return slice.Get(idx) >> val
}

// return slice[idx] < val
func GetLessThan[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) bool {
	return slice.Get(idx) < val
}

// return slice[idx1] < slice[idx2]
func GetLessThan2[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx1 IDX, idx2 IDX) bool {
	return slice.Get(idx1) < slice.Get(idx2)
}

// return slice[idx] <= val
func GetLessThanEqual[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) bool {
	return slice.Get(idx) <= val
}

// return slice[idx1] <= slice[idx2]
func GetLessThanEqual2[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx1 IDX, idx2 IDX) bool {
	return slice.Get(idx1) <= slice.Get(idx2)
}

// return slice[idx] > val
func GetGreaterThan[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) bool {
	return slice.Get(idx) > val
}

// return slice[idx1] > slice[idx2]
func GetGreaterThan2[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx1 IDX, idx2 IDX) bool {
	return slice.Get(idx1) > slice.Get(idx2)
}

// return slice[idx] >= val
func GetGreaterThanEqual[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) bool {
	return slice.Get(idx) >= val
}

// return slice[idx1] >= slice[idx2]
func GetGreaterThanEqual2[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx1 IDX, idx2 IDX) bool {
	return slice.Get(idx1) >= slice.Get(idx2)
}

// return slice[idx] == val
func GetEquals[T Equatable, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) bool {
	return slice.Get(idx) == val
}

// return slice[idx1] == slice[idx2]
func GetEquals2[T Equatable, IDX Integer, S SliceLike[T, IDX]](slice S, idx1 IDX, idx2 IDX) bool {
	return slice.Get(idx1) == slice.Get(idx2)
}

// return slice[idx] != val
func GetNotEquals[T Equatable, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) bool {
	return slice.Get(idx) != val
}

// return slice[idx1] != slice[idx2]
func GetNotEquals2[T Equatable, IDX Integer, S SliceLike[T, IDX]](slice S, idx1 IDX, idx2 IDX) bool {
	return slice.Get(idx1) != slice.Get(idx2)
}

// return min(slice[indexes[0]], slice[indexes[1]], ...)
func GetMinV[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, indexes ...IDX) T {
	idxSlice := NewSliceAdapter(indexes)
	return GetMin(slice, &idxSlice)
}

// return min(slice[indexes[0]], slice[indexes[1]], ...)
func GetMin[T Ordered, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[IDX1, IDX2]](slice S1, indexes S2) (minVal T) {
	i := indexes.FirstIdx()
	ok1 := indexes.IdxValid(i)
	if !ok1 {
		return
	}
	ii := indexes.Get(i)
	ok2 := slice.IdxValid(ii)
	if !ok2 {
		return
	}
	minVal = slice.Get(ii)
	i = indexes.NextIdx(i)
	ok1 = indexes.IdxValid(i)
	if !ok1 {
		return
	}
	ii = indexes.Get(i)
	ok2 = slice.IdxValid(ii)
	for ok1 && ok2 {
		minVal = min(minVal, slice.Get(ii))
		i = indexes.NextIdx(i)
		ok1 = indexes.IdxValid(i)
		if !ok1 {
			return
		}
		ii = indexes.Get(i)
		ok2 = slice.IdxValid(ii)
	}
	return
}

// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMinV[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, setIdx IDX, indexes ...IDX) {
	idxSlice := NewSliceAdapter(indexes)
	SetMin(slice, setIdx, &idxSlice)
}

// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMin[T Ordered, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[IDX1, IDX2]](slice S1, setIdx IDX1, indexes S2) {
	v := GetMin(slice, indexes)
	slice.Set(setIdx, v)
}

// return max(slice[indexes[0]], slice[indexes[1]], ...)
func GetMaxV[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, indexes ...IDX) T {
	idxSlice := NewSliceAdapter(indexes)
	return GetMax(slice, &idxSlice)
}

// return max(slice[indexes[0]], slice[indexes[1]], ...)
func GetMax[T Ordered, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[IDX1, IDX2]](slice S1, indexes S2) (maxVal T) {
	i := indexes.FirstIdx()
	ok1 := indexes.IdxValid(i)
	if !ok1 {
		return
	}
	ii := indexes.Get(i)
	ok2 := slice.IdxValid(ii)
	if !ok2 {
		return
	}
	maxVal = slice.Get(ii)
	i = indexes.NextIdx(i)
	ok1 = indexes.IdxValid(i)
	if !ok1 {
		return
	}
	ii = indexes.Get(i)
	ok2 = slice.IdxValid(ii)
	for ok1 && ok2 {
		maxVal = max(maxVal, slice.Get(ii))
		i = indexes.NextIdx(i)
		ok1 = indexes.IdxValid(i)
		if !ok1 {
			return
		}
		ii = indexes.Get(i)
		ok2 = slice.IdxValid(ii)
	}
	return
}

// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMaxV[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, setIdx IDX, indexes ...IDX) {
	idxSlice := NewSliceAdapter(indexes)
	SetMax(slice, setIdx, &idxSlice)
}

// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMax[T Ordered, IDX1 Integer, IDX2 Integer, S1 SliceLike[T, IDX1], S2 SliceLike[IDX1, IDX2]](slice S1, setIdx IDX1, indexes S2) {
	v := GetMax(slice, indexes)
	slice.Set(setIdx, v)
}

// return min(maxVal, max(slice[idx], minVal))
func GetClamped[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, minVal T, maxVal T) T {
	v := slice.Get(idx)
	return min(maxVal, max(v, minVal))
}

// slice[idx] = min(maxVal, max(slice[idx], minVal))
func SetClamped[T Ordered, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, minVal T, maxVal T) {
	v := GetClamped(slice, idx, minVal, maxVal)
	slice.Set(idx, v)
}

// oldVal := slice[idx]
// slice[idx] = val
// return oldVal != val
func SetChanged[T Equatable, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val T) (didChange bool) {
	v := slice.Get(idx)
	slice.Set(idx, val)
	return v != val
}

// return *(*TT)(unsafe.Pointer(&slice[idx]))
func GetUnsafeCast[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX) (val TT) {
	v := slice.Get(idx)
	return *(*TT)(unsafe.Pointer(&v))
}

// slice[idx] = *(*T)(unsafe.Pointer(&val))
func SetUnsafeCast[T any, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val TT) {
	v := *(*T)(unsafe.Pointer(&val))
	*(*TT)(unsafe.Pointer(&v)) = val
	slice.Set(idx, v)
}

// val_T := *(*T)(unsafe.Pointer(&val))
// oldVal_TT := *(*TT)(unsafe.Pointer(&slice[idx]))
// slice[idx] = val_T
// return oldVal_TT != val
func SetUnsafeCastChanged[T any, TT Equatable, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val TT) (didChange bool) {
	castNewVal := *(*T)(unsafe.Pointer(&val))
	v := slice.Get(idx)
	oldVal := *(*TT)(unsafe.Pointer(&v))
	slice.Set(idx, castNewVal)
	return oldVal != val
}

// val_T := *(*T)(unsafe.Pointer(&val))
// oldVal_T := slice[idx]
// slice[idx] = val_T
// return oldVal_T != val_T
func SetUnsafeCastChangedAlt[T Equatable, TT any, IDX Integer, S SliceLike[T, IDX]](slice S, idx IDX, val TT) (didChange bool) {
	castNewVal := *(*T)(unsafe.Pointer(&val))
	v := slice.Get(idx)
	slice.Set(idx, castNewVal)
	return v != castNewVal
}

// return *(*TT)(unsafe.Pointer(&slice[idx]))
func GetUnsafePtrCast[T any, TT any, IDX Integer, S MemSliceLike[T, IDX]](slice S, idx IDX) (val *TT) {
	v := slice.GetPtr(idx)
	return (*TT)(unsafe.Pointer(v))
}
