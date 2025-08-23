package go_list_like

import (
	"cmp"
	"math"
	"unsafe"
)

type integer interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~int | ~int8 | ~int16 | ~int32 | ~int64
}

type float interface {
	~float32 | ~float64
}

type number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// slice[idx] = slice[idx] + val
func SetAdd[T number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)+val)
}

// return slice[idx] + val
func GetAdd[T number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) + val
}

// slice[idx] = slice[idx] - val
func SetSubtract[T number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)-val)
}

// return slice[idx] - val
func GetSubtract[T number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) - val
}

// slice[idx] = slice[idx] * val
func SetMultiply[T number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)*val)
}

// return slice[idx] * val
func GetMultiply[T number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) * val
}

// slice[idx] = slice[idx] / val
func SetDivide[T number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)/val)
}

// return slice[idx] / val
func GetDivide[T number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) / val
}

// slice[idx] = slice[idx] % val
func SetModulo[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)%val)
}

// return slice[idx] % val
func GetModulo[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) % val
}

// return slice[idx] % val, slice[idx] - (slice[idx] % val)
func GetModRem[T integer](slice SliceLike[T], idx int, val T) (mod T, rem T) {
	v := slice.Get(idx)
	mod = v % val
	return mod, v - mod
}

// slice[idx] = math.Mod(slice[idx], val)
func SetFModulo[T float](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, T(math.Mod(float64(slice.Get(idx)), float64(val))))
}

// return math.Mod(slice[idx], val)
func GetFModulo[T float](slice SliceLike[T], idx int, val T) (result T) {
	return T(math.Mod(float64(slice.Get(idx)), float64(val)))
}

// return math.Mod(slice[idx], val), slice[idx] - math.Mod(slice[idx], val)
func GetFModRem[T float](slice SliceLike[T], idx int, val T) (mod T, rem T) {
	v := slice.Get(idx)
	mod = T(math.Mod(float64(slice.Get(idx)), float64(val)))
	return mod, v - mod
}

// slice[idx] = slice[idx] & val
func SetBitAnd[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)&val)
}

// return slice[idx] & val
func GetBitAnd[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) & val
}

// slice[idx] = slice[idx] & val
func SetBitOr[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)|val)
}

// return slice[idx] & val
func GetBitOr[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) | val
}

// slice[idx] = slice[idx] ^ val
func SetBitXor[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)^val)
}

// return slice[idx] ^ val
func GetBitXor[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) ^ val
}

// slice[idx] = ^slice[idx]
func SetBitInvert[T integer](slice SliceLike[T], idx int) {
	slice.Set(idx, ^slice.Get(idx))
}

// return ^slice[idx]
func GetBitInvert[T integer](slice SliceLike[T], idx int) (result T) {
	return ^slice.Get(idx)
}

// slice[idx] = slice[idx] << val
func SetBitLsh[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)<<val)
}

// return slice[idx] << val
func GetBitLsh[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) << val
}

// slice[idx] = slice[idx] >> val
func SetBitRsh[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)>>val)
}

// return slice[idx] >> val
func GetBitRsh[T integer](slice SliceLike[T], idx int, val T) (result T) {
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
func GetEquals[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool {
	return slice.Get(idx) == val
}

// return slice[idx1] == slice[idx2]
func GetEquals2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool {
	return slice.Get(idx1) == slice.Get(idx2)
}

// return slice[idx] != val
func GetNotEquals[T cmp.Ordered](slice SliceLike[T], idx int, val T) bool {
	return slice.Get(idx) != val
}

// return slice[idx1] != slice[idx2]
func GetNotEquals2[T cmp.Ordered](slice SliceLike[T], idx1 int, idx2 int) bool {
	return slice.Get(idx1) != slice.Get(idx2)
}

// return *(*TT)(unsafe.Pointer(&slice[idx]))
func GetUnsafeCast[T any, TT any](slice SliceLike[T], idx int) (val TT) {
	v := slice.Get(idx)
	return *(*TT)(unsafe.Pointer(&v))
}

// return *(*TT)(unsafe.Pointer(&slice[idx]))
func GetUnsafePtrCast[T any, TT any](slice MemSliceLike[T], idx int) (val *TT) {
	v := slice.GetPtr(idx)
	return (*TT)(unsafe.Pointer(v))
}
