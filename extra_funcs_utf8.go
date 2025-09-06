package go_list_like

import "unicode/utf8"

const (
	// this const table copied from [unicode/utf8]

	surrogateMin = 0xD800
	surrogateMax = 0xDFFF

	// t1 = 0b00000000
	tx = 0b10000000
	t2 = 0b11000000
	t3 = 0b11100000
	t4 = 0b11110000
	// t5 = 0b11111000

	maskx = 0b00111111
	mask2 = 0b00011111
	mask3 = 0b00001111
	mask4 = 0b00000111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1

	// The default lowest and highest continuation byte.
	locb = 0b10000000
	hicb = 0b10111111

	// These names of these constants are chosen to give nice alignment in the
	// table below. The first nibble is an index into acceptRanges or F for
	// special one-byte cases. The second nibble is the Rune length or the
	// Status for the special one-byte case.

	xx = 0xF1 // invalid: size 1
	as = 0xF0 // ASCII: size 1
	s1 = 0x02 // accept 0, size 2
	s2 = 0x13 // accept 1, size 3
	s3 = 0x03 // accept 0, size 3
	s4 = 0x23 // accept 2, size 3
	s5 = 0x34 // accept 3, size 4
	s6 = 0x04 // accept 0, size 4
	s7 = 0x44 // accept 4, size 4

	runeErrorByte0 = t3 | (utf8.RuneError >> 12)
	runeErrorByte1 = tx | (utf8.RuneError>>6)&maskx
	runeErrorByte2 = tx | utf8.RuneError&maskx
)

// first is information about the first byte in a UTF-8 sequence.
var utf8_first = [256]uint8{
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x00-0x0F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x10-0x1F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x20-0x2F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x30-0x3F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x40-0x4F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x50-0x5F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x60-0x6F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x70-0x7F
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x80-0x8F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x90-0x9F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xA0-0xAF
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xB0-0xBF
	xx, xx, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xC0-0xCF
	s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xD0-0xDF
	s2, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s4, s3, s3, // 0xE0-0xEF
	s5, s6, s6, s6, s7, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xF0-0xFF
}

// acceptRange gives the range of valid values for the second byte in a UTF-8
// sequence.
type acceptRange struct {
	lo uint8 // lowest value for second byte.
	hi uint8 // highest value for second byte.
}

// acceptRanges has size 16 to avoid bounds checks in the code that uses it.
var acceptRanges = [16]acceptRange{
	0: {locb, hicb},
	1: {0xA0, hicb},
	2: {locb, 0x9F},
	3: {0x90, hicb},
	4: {locb, 0x8F},
}

// Decode a rune from the byte slice at the given index
func ReadRune[IDX Integer, S SliceLike[byte, IDX]](slice S, idx IDX) (r rune, bytes IDX, ok bool) {
	n := slice.Len()
	if n < 1 {
		return utf8.RuneError, 0, false
	}
	b0 := slice.Get(idx)
	x := utf8_first[b0]
	if x >= as {
		// The following code simulates an additional check for x == xx and
		// handling the ASCII and invalid cases accordingly. This mask-and-or
		// approach prevents an additional branch.
		mask := rune(x) << 31 >> 31 // Create 0x0000 or 0xFFFF.
		r = rune(b0)&^mask | utf8.RuneError&mask
		return r, 1, r == utf8.RuneError
	}
	sz := IDX(x & 7)
	accept := acceptRanges[x>>4]
	if n < sz {
		return utf8.RuneError, 1, false
	}
	b1 := slice.Get(idx + 1)
	if b1 < accept.lo || accept.hi < b1 {
		return utf8.RuneError, 1, false
	}
	if sz <= 2 { // <= instead of == to help the compiler eliminate some bounds checks
		return rune(b0&mask2)<<6 | rune(b1&maskx), 2, true
	}
	b2 := slice.Get(idx + 2)
	if b2 < locb || hicb < b2 {
		return utf8.RuneError, 1, false
	}
	if sz <= 3 {
		return rune(b0&mask3)<<12 | rune(b1&maskx)<<6 | rune(b2&maskx), 3, true
	}
	b3 := slice.Get(idx + 3)
	if b3 < locb || hicb < b3 {
		return utf8.RuneError, 1, false
	}
	return rune(b0&mask4)<<18 | rune(b1&maskx)<<12 | rune(b2&maskx)<<6 | rune(b3&maskx), 4, true
}

// Write a rune to the byte slice at given index
func WriteRune[IDX Integer, S SliceLike[byte, IDX]](slice S, idx IDX, r rune) (bytes IDX, ok bool) {
	if uint32(r) <= rune1Max {
		slice.Set(idx, byte(r))
		return 1, true
	}
	return setRuneNonASCII(slice, idx, r)
}

// Append a rune to the end of the byte slice
func AppendRune[IDX Integer, L ListLike[byte, IDX]](list L, r rune) (bytes IDX, ok bool) {
	len := IDX(utf8.RuneLen(r))
	ok = list.TryEnsureFreeSlots(len)
	if !ok {
		return
	}
	idx, _ := list.AppendSlotsAssumeCapacity(len)
	bytes, ok = WriteRune(list, idx, r)
	return
}

func setRuneNonASCII[IDX Integer, S SliceLike[byte, IDX]](slice S, idx IDX, r rune) (bytes IDX, ok bool) {
	switch i := uint32(r); {
	case i <= rune2Max:
		slice.Set(idx, t2|byte(r>>6))
		slice.Set(idx+1, tx|byte(r)&maskx)
		return 2, true
	case i < surrogateMin, surrogateMax < i && i <= rune3Max:
		slice.Set(idx, t3|byte(r>>12))
		slice.Set(idx+1, tx|byte(r>>6)&maskx)
		slice.Set(idx+2, tx|byte(r)&maskx)
		return 3, true
	case i > rune3Max && i <= utf8.MaxRune:
		slice.Set(idx, t4|byte(r>>18))
		slice.Set(idx+1, tx|byte(r>>12)&maskx)
		slice.Set(idx+2, tx|byte(r>>6)&maskx)
		slice.Set(idx+3, tx|byte(r)&maskx)
		return 4, true
	default:
		slice.Set(idx, runeErrorByte0)
		slice.Set(idx+1, runeErrorByte1)
		slice.Set(idx+2, runeErrorByte2)
		return 3, false
	}
}
