package go_list_like

import (
	"math"
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
func AddSet[T number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)+val)
}

// return slice[idx] + val
func GetAdd[T number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) + val
}

// slice[idx] = slice[idx] - val
func SubSet[T number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)-val)
}

// return slice[idx] - val
func GetSub[T number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) - val
}

// slice[idx] = slice[idx] * val
func MultSet[T number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)*val)
}

// return slice[idx] * val
func GetMult[T number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) * val
}

// slice[idx] = slice[idx] / val
func DivSet[T number](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)/val)
}

// return slice[idx] / val
func GetDiv[T number](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) / val
}

// slice[idx] = slice[idx] % val
func ModSet[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)%val)
}

// return slice[idx] % val
func GetMod[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) % val
}

// return slice[idx] % val, slice[idx] - (slice[idx] % val)
func GetModRem[T integer](slice SliceLike[T], idx int, val T) (mod T, rem T) {
	v := slice.Get(idx)
	mod = v % val
	return mod, v - mod
}

// slice[idx] = math.Mod(slice[idx], val)
func FModSet[T float](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, T(math.Mod(float64(slice.Get(idx)), float64(val))))
}

// return math.Mod(slice[idx], val)
func GetFMod[T float](slice SliceLike[T], idx int, val T) (result T) {
	return T(math.Mod(float64(slice.Get(idx)), float64(val)))
}

// return math.Mod(slice[idx], val), slice[idx] - math.Mod(slice[idx], val)
func GetFModRem[T float](slice SliceLike[T], idx int, val T) (mod T, rem T) {
	v := slice.Get(idx)
	mod = T(math.Mod(float64(slice.Get(idx)), float64(val)))
	return mod, v - mod
}

// slice[idx] = slice[idx] & val
func BitAndSet[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)&val)
}

// return slice[idx] & val
func GetBitAnd[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) & val
}

// slice[idx] = slice[idx] & val
func BitOrSet[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)|val)
}

// return slice[idx] & val
func GetBitOr[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) | val
}

// slice[idx] = slice[idx] ^ val
func BitXorSet[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)^val)
}

// return slice[idx] ^ val
func GetBitXor[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) ^ val
}

// slice[idx] = ^slice[idx]
func BitNotSet[T integer](slice SliceLike[T], idx int) {
	slice.Set(idx, ^slice.Get(idx))
}

// return ^slice[idx]
func GetBitNot[T integer](slice SliceLike[T], idx int) (result T) {
	return ^slice.Get(idx)
}

// slice[idx] = slice[idx] << val
func BitLshSet[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)<<val)
}

// return slice[idx] << val
func GetBitLsh[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) << val
}

// slice[idx] = slice[idx] >> val
func BitRshSet[T integer](slice SliceLike[T], idx int, val T) {
	slice.Set(idx, slice.Get(idx)>>val)
}

// return slice[idx] >> val
func GetBitRsh[T integer](slice SliceLike[T], idx int, val T) (result T) {
	return slice.Get(idx) >> val
}
