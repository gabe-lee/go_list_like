package go_list_like

import (
	"cmp"
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

// slice[idx] = slice[idx] + val
func SetAdd[T Number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)+val)
}

// return slice[idx] + val
func GetAdd[T Number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) + val
}

// slice[idx] = slice[idx] - val
func SetSubtract[T Number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)-val)
}

// return slice[idx] - val
func GetSubtract[T Number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) - val
}

// slice[idx] = slice[idx] * val
func SetMultiply[T Number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)*val)
}

// return slice[idx] * val
func GetMultiply[T Number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) * val
}

// slice[idx] = slice[idx] / val
func SetDivide[T Number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)/val)
}

// return slice[idx] / val
func GetDivide[T Number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) / val
}

// slice[idx] = slice[idx] % val
func SetModulo[T Integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)%val)
}

// return slice[idx] % val
func GetModulo[T Integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) % val
}

// return slice[idx] % val, slice[idx] - (slice[idx] % val)
func GetModRem[T Integer](slice SliceLike[T], idx int, val T) (mod T, rem T) {
	v := slice.Get(idx)
	mod = v % val
	return mod, v - mod
}

// slice[idx] = math.Mod(slice[idx], val)
func SetFModulo[T Float](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, T(math.Mod(float64(slice.Get(idx)), float64(val))))
}

// return math.Mod(slice[idx], val)
func GetFModulo[T Float](slice SliceLike[T], idx int, val T) (result T) {
	return T(math.Mod(float64(slice.Get(idx)), float64(val)))
}

// return math.Mod(slice[idx], val), slice[idx] - math.Mod(slice[idx], val)
func GetFModRem[T Float](slice SliceLike[T], idx int, val T) (mod T, rem T) {
	v := slice.Get(idx)
	mod = T(math.Mod(float64(slice.Get(idx)), float64(val)))
	return mod, v - mod
}

// slice[idx] = slice[idx] & val
func SetBitAnd[T Integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)&val)
}

// return slice[idx] & val
func GetBitAnd[T Integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) & val
}

// slice[idx] = slice[idx] & val
func SetBitOr[T Integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)|val)
}

// return slice[idx] & val
func GetBitOr[T Integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) | val
}

// slice[idx] = slice[idx] ^ val
func SetBitXor[T Integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)^val)
}

// return slice[idx] ^ val
func GetBitXor[T Integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) ^ val
}

// slice[idx] = ^slice[idx]
func SetBitInvert[T Integer](slice SliceLike[T], idx int) {
	slice.Set(idx, ^slice.Get(idx))
}

// return ^slice[idx]
func GetBitInvert[T Integer](slice SliceLike[T], idx int) (result T) {
	return ^slice.Get(idx)
}

// slice[idx] = slice[idx] << val
func SetBitLsh[T Integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)<<val)
}

// return slice[idx] << val
func GetBitLsh[T Integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) << val
}

// slice[idx] = slice[idx] >> val
func SetBitRsh[T Integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)>>val)
}

// return slice[idx] >> val
func GetBitRsh[T Integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) >> val
}

// return slice[idx] < val
func GetLessThan[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool {
	return slice.Get(idx) < val
}

// return slice[idx1] < slice[idx2]
func GetLessThan2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool {
	return slice.Get(idx1) < slice.Get(idx2)
}

// return slice[idx] <= val
func GetLessThanEqual[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool {
	return slice.Get(idx) <= val
}

// return slice[idx1] <= slice[idx2]
func GetLessThanEqual2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool {
	return slice.Get(idx1) <= slice.Get(idx2)
}

// return slice[idx] > val
func GetGreaterThan[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool {
	return slice.Get(idx) > val
}

// return slice[idx1] > slice[idx2]
func GetGreaterThan2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool {
	return slice.Get(idx1) > slice.Get(idx2)
}

// return slice[idx] >= val
func GetGreaterThanEqual[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool {
	return slice.Get(idx) >= val
}

// return slice[idx1] >= slice[idx2]
func GetGreaterThanEqual2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool {
	return slice.Get(idx1) >= slice.Get(idx2)
}

// return slice[idx] == val
func GetEquals[T Equatable](slice SliceLike[T], idx int, val T) bool {
	return slice.Get(idx) == val
}

// return slice[idx1] == slice[idx2]
func GetEquals2[T Equatable](slice SliceLike[T], idx1 int, idx2 int) bool {
	return slice.Get(idx1) == slice.Get(idx2)
}

// return slice[idx] != val
func GetNotEquals[T Equatable](slice SliceLike[T], idx int, val T) bool {
	return slice.Get(idx) != val
}

// return slice[idx1] != slice[idx2]
func GetNotEquals2[T Equatable](slice SliceLike[T], idx1 int, idx2 int) bool {
	return slice.Get(idx1) != slice.Get(idx2)
}

// return min(slice[indexes[0]], slice[indexes[1]], ...)
func GetMinV[T cmp.Ordered](slice SliceLike[T], indexes ...int) T {
	idxSlice := NewSliceAdapter(indexes)
	return GetMin(slice, &idxSlice)
}

// return min(slice[indexes[0]], slice[indexes[1]], ...)
func GetMin[T cmp.Ordered, I Index](slice SliceLike[T], indexes SliceLike[I]) T {
	v := slice.Get(int(indexes.Get(0)))
	i := 1
	limit := indexes.Len()
	for i < limit {
		v = min(v, slice.Get(int(indexes.Get(0))))
	}
	return v
}

// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMinV[T cmp.Ordered](slice SliceLike[T], setIdx int, indexes ...int) {
	idxSlice := NewSliceAdapter(indexes)
	SetMin(slice, setIdx, &idxSlice)
}

// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMin[T cmp.Ordered, I Index](slice SliceLike[T], setIdx int, indexes SliceLike[I]) {
	v := GetMin(slice, indexes)
	slice.Set(setIdx, v)
}

// return max(slice[indexes[0]], slice[indexes[1]], ...)
func GetMaxV[T cmp.Ordered](slice SliceLike[T], indexes ...int) T {
	idxSlice := NewSliceAdapter(indexes)
	return GetMax(slice, &idxSlice)
}

// return max(slice[indexes[0]], slice[indexes[1]], ...)
func GetMax[T cmp.Ordered, I Index](slice SliceLike[T], indexes SliceLike[I]) T {
	v := slice.Get(int(indexes.Get(0)))
	i := 1
	limit := indexes.Len()
	for i < limit {
		v = max(v, slice.Get(int(indexes.Get(0))))
	}
	return v
}

// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMaxV[T cmp.Ordered](slice SliceLike[T], setIdx int, indexes ...int) {
	idxSlice := NewSliceAdapter(indexes)
	SetMax(slice, setIdx, &idxSlice)
}

// slice[setIdx] = min(slice[indexes[0]], slice[indexes[1]], ...)
func SetMax[T cmp.Ordered, I Index](slice SliceLike[T], setIdx int, indexes SliceLike[I]) {
	v := GetMax(slice, indexes)
	slice.Set(setIdx, v)
}

// return min(maxVal, max(slice[idx], minVal))
func GetClamped[T cmp.Ordered](slice SliceLike[T], idx int, minVal T, maxVal T) T {
	v := slice.Get(idx)
	return min(maxVal, max(v, minVal))
}

// slice[idx] = min(maxVal, max(slice[idx], minVal))
func SetClamped[T cmp.Ordered](slice SliceLike[T], idx int, minVal T, maxVal T) {
	v := GetClamped(slice, idx, minVal, maxVal)
	slice.Set(idx, v)
}

// oldVal := slice[idx]
// slice[idx] = val
// return oldVal != val
func SetChanged[T Equatable](slice SliceLike[T], idx int, val T) (didChange bool) {
	v := slice.Get(idx)
	slice.Set(idx, val)
	return v != val
}

// return *(*TT)(unsafe.Pointer(&slice[idx]))
func GetUnsafeCast[T any, TT any](slice SliceLike[T], idx int) (val TT) {
	v := slice.Get(idx)
	return *(*TT)(unsafe.Pointer(&v))
}

// slice[idx] = *(*T)(unsafe.Pointer(&val))
func SetUnsafeCast[T any, TT any](slice SliceLike[T], idx int, val TT) {
	v := *(*T)(unsafe.Pointer(&val))
	*(*TT)(unsafe.Pointer(&v)) = val
}

// val_T := *(*T)(unsafe.Pointer(&val))
// oldVal_TT := *(*TT)(unsafe.Pointer(&slice[idx]))
// slice[idx] = val_T
// return oldVal_TT != val
func SetUnsafeCastChanged[T any, TT Equatable](slice SliceLike[T], idx int, val TT) (didChange bool) {
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
func SetUnsafeCastChangedAlt[T Equatable, TT any](slice SliceLike[T], idx int, val TT) (didChange bool) {
	castNewVal := *(*T)(unsafe.Pointer(&val))
	v := slice.Get(idx)
	slice.Set(idx, castNewVal)
	return v != castNewVal
}

// return *(*TT)(unsafe.Pointer(&slice[idx]))
func GetUnsafePtrCast[T any, TT any](slice MemSliceLike[T], idx int) (val *TT) {
	v := slice.GetPtr(idx)
	return (*TT)(unsafe.Pointer(v))
}
