package implementation_test

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"testing"

	LL "github.com/gabe-lee/go_list_like"
)

const (
	_MaxValuesPerOp byte = 7
)

// Set this to `true` to insert a `runtime.Breakpoint()` at the start of each
// implementation operation test
var EnableBreakpoints bool = false

const (
	// SliceLike ops

	_ImplOpGet byte = iota
	_ImplOpSet
	_ImplOpMove
	_ImplOpMoveRange
	_ImplOpSlice
	_ImplOpFirst
	_ImplOpLast
	_ImplOpNext
	_ImplOpNthNext
	_ImplOpPrev
	_ImplOpNthPrev

	// ListLike ops

	_ImplOpInsert
	_ImplOpAppend
	_ImplOpDelete

	// opCounts

	_opCount
	_sliceOpCount = _ImplOpSlice + 1
)

var _OpNames = [_opCount]string{
	_ImplOpGet:       "Get",
	_ImplOpSet:       "Set",
	_ImplOpMove:      "Move",
	_ImplOpMoveRange: "MoveRange",
	_ImplOpSlice:     "Slice",
	_ImplOpFirst:     "FirstIdx",
	_ImplOpLast:      "LastIdx",
	_ImplOpNext:      "NextIdx",
	_ImplOpNthNext:   "NthNextIdx",
	_ImplOpPrev:      "PrevIdx",
	_ImplOpNthPrev:   "NthPrevIdx",
	_ImplOpInsert:    "InsertSlotsAssumeCapacity",
	_ImplOpAppend:    "AppendSlotsAssumeCapacity",
	_ImplOpDelete:    "DeleteRange",
}

func InitImplementationFuzz(f *testing.F) {
	var nilBytes []byte
	var emptyBytes []byte = make([]byte, 0, 10)
	f.Add(nilBytes, []byte{0})
	f.Add(emptyBytes, []byte{0})
	f.Add([]byte{0}, []byte{0})
	f.Add([]byte{1}, []byte{1})
	f.Add([]byte{255}, []byte{255})
	f.Add([]byte{254}, []byte{255})
	f.Add(
		[]byte{129, 169, 201, 61, 232, 249, 94, 132, 135, 203, 186, 17},
		[]byte{38, 221, 92, 182, 48, 203, 165, 12, 225, 40, 46, 189, 150, 53, 180, 191, 127, 39, 23, 194, 58, 196, 95, 217, 66, 109, 127, 208, 250, 68, 70, 49},
	)
	f.Add(
		[]byte{146, 62, 12, 146, 72, 47, 250, 30, 22, 59, 218, 80, 47, 186, 248, 129, 89, 48, 103, 108, 145, 68, 236, 45, 24, 102, 95, 166, 186, 174, 201, 80},
		[]byte{195, 25, 215, 144, 245, 229, 189, 184, 121, 152, 2, 212, 121, 49, 148, 242, 208, 38, 105, 139, 235, 96, 244, 174, 4, 105, 117, 226, 131, 170, 101, 11},
	)
	f.Add(
		[]byte{219, 248, 124, 167, 82, 50, 207, 91, 114, 5, 54, 190, 52, 116, 214, 16, 29, 119, 21, 188, 138, 57, 130, 160, 115, 172, 219, 78, 3, 224, 17, 222, 195, 1, 190, 203, 215, 202, 89, 12, 206, 236, 231, 189, 177, 211, 204, 94, 49, 69, 126, 102, 75, 226, 231, 228, 103, 58, 34, 45, 60, 118, 206, 1, 52, 103, 99},
		[]byte{102, 86, 158, 20, 0, 73, 69, 222, 151, 223, 184, 252, 81, 43, 136, 148, 251, 190, 18, 27, 120, 79, 176, 83, 54, 108, 219, 172, 168, 122, 227, 87, 114, 91, 107, 120, 83, 91, 182, 213, 118, 15, 46, 120, 107, 125, 252, 172, 244, 48, 250, 160, 83, 51, 247, 169, 85, 207, 146, 235, 253, 183, 122, 90, 187, 43, 45},
	)
	f.Add(
		[]byte{6, 43, 44, 41, 95, 174, 27, 84, 144, 109, 62, 157, 110, 62, 164, 48, 152, 126, 161, 178, 4, 53, 93, 249, 22, 246, 37, 183, 190, 196, 109, 238, 20, 162, 82, 80, 213, 173, 221, 110, 214, 89, 232, 142, 219, 253, 40, 56, 50, 232, 92, 140, 145, 137, 182, 235, 47, 38, 246, 158, 228, 181, 190, 24, 234, 132, 249},
		[]byte{59, 48, 197, 235, 24, 36, 108, 204, 90, 12, 235, 175, 51, 46, 60, 96, 108, 212, 182, 240, 165, 180, 217, 159, 144, 155, 173, 159, 205, 196, 215, 208, 121, 162, 20, 131, 18, 196, 92, 74, 160, 236, 206, 144, 157, 249, 48, 238, 201, 102, 88, 169, 144, 195, 19, 94, 155, 99, 236, 199, 40, 95, 108, 67, 224, 186, 187},
	)
	f.Add(
		[]byte{0, 180, 33, 132, 134, 3, 154, 103, 99, 5, 10, 176, 99, 93, 76, 195, 205, 184, 84, 140, 110, 20, 27, 171, 44, 151, 3, 246, 161, 171, 142, 218, 74, 164, 248, 231, 32, 75, 120, 37, 180, 138, 144, 14, 251, 125, 150, 63, 193, 177, 41, 31, 66, 31, 2, 99, 15, 49, 250, 88, 23, 135, 158, 74, 189, 203, 102, 72, 201, 212, 176, 236, 196, 79, 122, 170, 74, 236, 47, 237, 179, 214, 139, 126, 111, 90, 29, 200, 110, 196, 127, 68, 151, 215, 117, 170, 154, 34, 248, 235, 164, 99, 168, 6, 98, 152, 255, 178, 31, 178, 127, 123, 82, 164, 85, 121, 205, 240, 236, 8, 234, 49, 153, 125, 246, 59, 243, 186, 32, 208, 62, 247, 209, 54, 155, 80, 218, 169, 3, 177, 108, 123, 120, 238, 213, 93, 52, 84, 100, 122, 152, 44, 225, 118, 20, 5, 123, 98, 158, 205, 121, 66, 134, 173, 118, 113, 57, 196, 114, 93, 51, 236, 143, 246, 3, 194, 244, 83, 36, 226, 117, 138, 124, 217, 238, 124, 9, 111, 183, 99, 131, 142, 35, 186, 199, 59, 25, 40, 56, 62, 71, 87, 96, 88, 237, 39, 93, 6, 253, 249, 60, 100, 248, 155, 90, 200, 129, 76, 92, 27, 134, 117, 22, 98, 126, 53, 79, 192, 29, 21, 223, 128, 88, 128, 49, 21, 165, 157, 82, 167, 201, 230, 82, 169, 30, 63, 41, 60, 103, 231, 45, 211, 174, 236, 198, 78, 126, 110, 63, 42, 99, 222, 39, 170, 173, 126, 216, 136, 37, 239, 28, 89, 170, 43, 199, 27, 252, 76, 80, 107, 3, 173, 27, 253, 139, 158, 89, 215, 241, 121, 50, 133, 42, 203, 44, 54, 97, 223, 108, 100, 55, 4, 253, 29, 23, 219, 242, 131, 154, 100, 63, 170, 118, 27, 87, 104, 105, 55, 46, 79, 199, 30, 67, 197, 23, 207, 47, 11, 143, 39, 53, 114, 199, 78, 255, 145, 85, 111, 206, 129, 10, 195, 96, 191, 86, 170, 187, 164, 142, 137, 173, 233, 11, 235, 40, 138, 87, 77, 65, 54, 88, 89, 206, 3, 139, 68, 206, 204, 223, 187, 247, 72, 253, 28, 138, 228, 190, 99, 157, 7, 133, 250, 81, 116, 178, 239, 79, 203, 152, 179, 112, 194, 105, 34, 231, 104, 210, 226, 153, 170, 59, 74, 203, 14, 77, 183, 210, 102, 40, 85, 63, 36, 255, 198, 104, 152, 22, 20, 178, 177, 56, 50, 200, 114, 62, 219, 220, 27, 192, 232, 134, 210, 13, 111, 178, 174, 154, 148, 30, 54, 39, 195, 215, 77, 22, 160, 255, 135, 99, 114, 190, 233, 240, 154, 148, 110, 225, 10, 68, 132, 218, 40, 220, 246, 102, 99, 119, 129, 40, 84, 84, 203, 183, 156, 238, 220, 234, 233, 159, 167, 243, 120, 191, 224, 10, 57, 131, 206, 123, 190, 177, 217, 92, 57, 218, 159, 10, 180, 170, 47, 186, 54, 147, 240, 16, 254, 18, 52, 189, 15, 50, 45, 96, 163, 182, 102, 87, 119, 184, 148, 254, 65, 240, 121, 149, 219, 248, 45, 36, 127, 233, 61, 175, 209, 174, 115, 141, 130, 211, 78, 185, 147, 61, 155, 248, 11, 8, 24, 218, 83, 73, 32, 58, 50, 73, 14, 113, 209, 118, 32, 177, 169, 221, 10, 205, 2, 75, 199, 226, 175, 86, 84, 73, 234, 241, 97, 128, 100, 54, 8, 253, 172, 132, 161, 169, 99, 0, 196, 162, 241, 25, 245, 144, 55, 75, 103, 226, 240, 120, 232, 250, 97, 242, 95, 240, 158, 185, 141, 20, 87, 179, 182, 171, 35, 179, 209, 174, 4, 191, 95, 140, 219, 69, 146, 69, 89, 225, 150, 164, 236, 101, 201, 86, 104, 44, 212, 135, 87, 87, 165, 150, 171, 160, 241, 70, 196, 196, 150, 238, 62, 91, 93, 89, 82, 67, 254, 165, 134, 107, 214, 134, 178, 200, 150, 171, 69, 51, 79, 183, 59, 33, 101, 160, 94, 113, 120, 116, 202, 48, 246, 166, 161, 231, 128, 113, 86, 238, 104, 23, 16, 172, 245, 72, 232, 20, 97, 20, 117, 199, 33, 132, 220, 17, 55, 24, 0, 191, 149, 124, 12, 38, 222, 4, 213, 17, 97, 23, 27, 36, 19, 4, 153, 55, 166, 58, 122, 192, 216, 134, 193, 78, 193, 120, 171, 167, 46, 154, 83, 9, 14, 36, 125, 223, 53, 61, 48, 7, 146, 115, 180, 79, 20, 175, 253, 188, 59, 241, 100, 199, 58, 197, 159, 240, 70, 164, 142, 164, 219, 121, 171, 45, 37, 196, 9, 96, 9, 219, 66, 207, 203, 186, 116, 35, 101, 154, 154, 234, 74, 139, 76, 36, 200, 240, 130, 244, 15, 96, 231, 44, 195, 6, 166, 221, 209, 145, 241, 70, 180, 208, 16, 79, 127, 218, 202, 67, 204, 20, 221, 150, 64, 222, 123, 81, 250, 201, 99, 240, 48, 76, 46, 214, 118, 112, 60, 168, 166, 28, 162, 129, 124, 209, 146, 180, 83, 177, 146, 135, 178, 91, 133, 209, 137, 23, 42, 22, 111, 209, 36, 141, 183, 208, 1, 183, 171, 46, 116, 133, 63, 82, 30, 239, 126, 138, 108, 140, 34, 135, 180, 253, 56, 117, 32, 8, 190, 235, 126, 54, 254, 142, 20, 226, 88, 4, 47, 221, 195, 110, 176, 51, 82, 201, 103, 112, 102, 187, 138, 179, 118, 202, 145, 207, 231, 188, 160, 32, 250, 217, 34, 19, 197, 114, 63, 253, 176, 252, 246, 19, 222, 230, 140, 102, 41, 49, 84, 201, 23, 226, 84, 79, 222, 52, 162, 202, 180, 121, 157, 55, 101, 208, 129, 221, 72, 226, 96, 18, 255, 228, 141, 119, 245, 50, 21, 8, 232, 196, 15, 181, 173, 21, 32, 91, 233, 173, 97, 113, 169, 58, 34, 28, 207, 7, 39, 157, 232, 145, 87, 130, 208, 200, 24, 12, 24, 170, 174, 94, 190, 165, 128, 193, 22},
		[]byte{199, 38, 103, 199, 155, 181, 110, 201, 139, 137, 106, 17, 116, 3, 74, 125, 125, 45, 92, 210, 44, 166, 187, 31, 144, 239, 211, 7, 159, 33, 64, 108, 167, 217, 84, 101, 193, 143, 5, 238, 127, 172, 116, 178, 254, 83, 156, 158, 53, 31, 246, 66, 169, 52, 118, 190, 25, 38, 117, 14, 121, 197, 18, 131, 153, 127, 164, 21, 172, 169, 34, 13, 163, 183, 9, 177, 245, 232, 141, 95, 78, 51, 64, 154, 236, 211, 15, 181, 178, 161, 92, 191, 35, 186, 248, 40, 101, 204, 189, 45, 196, 148, 22, 227, 48, 120, 8, 94, 208, 76, 57, 103, 47, 19, 89, 64, 221, 33, 94, 33, 133, 159, 83, 117, 93, 205, 210, 145, 176, 72, 218, 9, 95, 181, 33, 233, 155, 138, 134, 51, 169, 236, 231, 0, 90, 25, 32, 204, 136, 208, 133, 174, 68, 191, 186, 146, 82, 34, 6, 186, 14, 83, 135, 182, 158, 179, 77, 177, 4, 74, 53, 246, 94, 17, 162, 169, 0, 107, 56, 75, 247, 213, 240, 74, 186, 204, 192, 129, 62, 116, 85, 2, 86, 212, 103, 108, 201, 62, 192, 150, 102, 132, 26, 132, 103, 81, 128, 244, 233, 241, 171, 69, 102, 239, 11, 245, 28, 8, 239, 78, 6, 12, 193, 42, 34, 214, 193, 45, 118, 164, 209, 152, 110, 83, 231, 11, 124, 26, 138, 9, 144, 190, 103, 68, 135, 25, 18, 98, 188, 41, 53, 243, 103, 137, 230, 254, 9, 160, 199, 4, 121, 190, 132, 1, 201, 180, 147, 183, 105, 115, 21, 135, 116, 103, 81, 224, 212, 71, 38, 245, 243, 169, 206, 228, 66, 167, 164, 0, 29, 143, 143, 222, 239, 162, 160, 208, 171, 33, 135, 159, 79, 132, 197, 25, 111, 67, 86, 141, 50, 176, 190, 28, 96, 2, 163, 83, 166, 161, 18, 99, 144, 212, 78, 187, 160, 50, 88, 110, 61, 197, 3, 145, 206, 91, 35, 115, 172, 207, 139, 1, 101, 97, 28, 216, 213, 98, 74, 155, 74, 11, 226, 226, 132, 99, 64, 143, 188, 213, 164, 209, 27, 17, 246, 160, 231, 43, 139, 181, 61, 133, 69, 33, 214, 8, 37, 157, 194, 4, 111, 48, 69, 128, 131, 101, 6, 48, 168, 34, 76, 218, 182, 74, 6, 41, 168, 9, 168, 251, 112, 11, 225, 72, 244, 100, 177, 234, 195, 234, 190, 12, 26, 213, 27, 178, 150, 6, 78, 78, 63, 57, 12, 196, 237, 70, 207, 228, 112, 196, 240, 181, 170, 179, 212, 204, 105, 17, 79, 102, 50, 184, 168, 231, 95, 59, 69, 113, 15, 40, 187, 120, 73, 53, 113, 217, 1, 24, 79, 223, 39, 110, 30, 25, 216, 100, 249, 224, 216, 240, 193, 0, 91, 106, 186, 195, 156, 151, 223, 67, 27, 28, 93, 77, 166, 235, 227, 11, 212, 16, 217, 105, 152, 49, 169, 122, 86, 71, 171, 116, 177, 241, 96, 60, 126, 141, 149, 171, 73, 66, 30, 23, 239, 181, 39, 129, 207, 28, 223, 34, 120, 135, 90, 42, 29, 244, 110, 37, 194, 24, 108, 191, 209, 57, 199, 146, 13, 163, 126, 160, 199, 137, 196, 208, 189, 50, 13, 140, 101, 29, 6, 249, 205, 68, 41, 142, 54, 49, 233, 160, 117, 52, 9, 123, 188, 8, 230, 220, 40, 162, 60, 201, 244, 209, 235, 103, 206, 9, 12, 2, 52, 171, 42, 112, 72, 116, 247, 118, 98, 45, 15, 31, 40, 237, 87, 63, 68, 153, 104, 237, 62, 58, 107, 22, 228, 151, 28, 11, 214, 11, 20, 11, 51, 86, 93, 148, 105, 76, 82, 243, 168, 77, 173, 232, 209, 71, 236, 73, 34, 117, 159, 244, 126, 228, 160, 89, 246, 211, 133, 85, 219, 204, 34, 250, 80, 62, 8, 67, 48, 65, 83, 51, 255, 8, 99, 3, 4, 45, 239, 53, 20, 27, 87, 247, 163, 218, 250, 204, 245, 21, 250, 208, 204, 246, 237, 18, 108, 101, 106, 211, 161, 60, 27, 18, 86, 250, 99, 165, 61, 164, 104, 10, 97, 142, 5, 215, 249, 202, 2, 172, 127, 107, 38, 254, 20, 51, 144, 243, 4, 79, 3, 34, 27, 149, 95, 134, 45, 71, 23, 20, 211, 44, 23, 200, 7, 94, 190, 74, 248, 254, 130, 121, 75, 54, 181, 157, 27, 206, 228, 198, 114, 196, 213, 203, 188, 69, 99, 6, 166, 235, 42, 219, 130, 211, 153, 68, 25, 241, 234, 76, 26, 38, 249, 182, 134, 83, 39, 100, 9, 169, 199, 245, 49, 114, 84, 195, 243, 114, 85, 45, 250, 83, 47, 160, 126, 220, 38, 158, 164, 237, 135, 70, 230, 141, 145, 54, 31, 201, 89, 0, 31, 45, 249, 209, 245, 14, 37, 114, 70, 33, 129, 140, 152, 154, 49, 142, 195, 50, 159, 92, 236, 173, 156, 136, 127, 55, 181, 215, 127, 55, 163, 161, 63, 159, 13, 63, 217, 20, 22, 235, 24, 141, 104, 143, 235, 197, 148, 103, 45, 24, 155, 222, 230, 89, 203, 80, 157, 133, 13, 59, 36, 9, 146, 204, 120, 62, 123, 132, 138, 36, 213, 100, 71, 237, 167, 88, 136, 94, 21, 58, 60, 106, 149, 224, 200, 190, 221, 126, 164, 243, 177, 27, 161, 250, 7, 94, 133, 135, 144, 181, 52, 117, 186, 158, 60, 167, 61, 147, 236, 199, 201, 28, 194, 46, 215, 190, 186, 20, 190, 208, 189, 57, 162, 67, 164, 168, 245, 171, 41, 210, 33, 236, 114, 49, 50, 96, 212, 5, 123, 121, 176, 249, 241, 102, 184, 28, 195, 141, 59, 65, 220, 124, 1, 242, 240, 138, 156, 74, 203, 126, 110, 171, 23, 143, 90, 99, 88, 97, 39, 161, 241, 43, 212, 65, 121, 80, 84, 166, 80, 4, 9, 163, 175, 76, 124, 148, 48, 64, 152, 242, 222, 75, 89, 75, 169, 239, 227, 254, 45, 95, 160, 240},
	)
	f.Add(
		[]byte{117, 162, 32, 132, 36, 130, 62, 108, 178, 114, 148, 42, 150, 12, 171, 166, 38, 251, 150, 32, 172, 176, 20, 125, 205, 84, 46, 162, 87, 147, 229, 116, 242, 121, 248, 11, 130, 18, 217, 197, 190, 70, 185, 129, 100, 229, 170, 91, 178, 38, 131, 143, 118, 140, 13, 227, 228, 121, 140, 228, 51, 211, 134, 64, 54, 100, 187, 253, 235, 162, 69, 182, 38, 132, 85, 218, 126, 123, 83, 126, 99, 212, 137, 186, 116, 207, 42, 145, 131, 185, 171, 113, 17, 215, 194, 174, 230, 8, 166, 200, 145, 34, 250, 17, 31, 66, 3, 12, 1, 10, 234, 156, 233, 206, 225, 176, 213, 17, 179, 218, 162, 198, 210, 59, 111, 141, 183, 92, 62, 196, 86, 21, 10, 34, 16, 76, 119, 94, 45, 120, 138, 81, 54, 115, 169, 43, 16, 136, 137, 24, 148, 246, 92, 96, 237, 46, 3, 77, 118, 124, 122, 22, 95, 206, 207, 6, 186, 55, 55, 16, 7, 196, 8, 82, 38, 115, 19, 17, 18, 185, 30, 111, 244, 176, 198, 67, 161, 188, 28, 232, 101, 172, 252, 39, 171, 205, 18, 3, 175, 146, 20, 108, 189, 115, 222, 217, 113, 121, 235, 169, 9, 163, 156, 5, 220, 161, 46, 148, 198, 97, 71, 157, 147, 109, 224, 209, 72, 74, 113, 8, 239, 57, 163, 56, 62, 3, 251, 44, 176, 202, 114, 172, 216, 160, 255, 136, 17, 44, 37, 129, 212, 121, 199, 132, 103, 243, 239, 196, 57, 27, 8, 146, 34, 45, 91, 193, 119, 141, 61, 24, 141, 133, 233, 87, 150, 215, 157, 161, 207, 26, 220, 237, 151, 45, 222, 67, 18, 107, 53, 112, 103, 148, 72, 183, 145, 137, 90, 1, 50, 59, 103, 96, 160, 254, 83, 154, 191, 114, 98, 226, 207, 40, 104, 53, 215, 191, 51, 119, 144, 213, 92, 96, 88, 232, 88, 91, 237, 175, 163, 61, 203, 56, 130, 229, 171, 7, 182, 35, 57, 106, 137, 139, 226, 36, 154, 142, 44, 219, 18, 199, 120, 103, 101, 2, 222, 65, 47, 68, 221, 132, 7, 69, 144, 254, 144, 0, 49, 82, 110, 77, 148, 119, 58, 130, 16, 164, 65, 28, 203, 135, 2, 105, 205, 178, 144, 66, 150, 4, 80, 109, 71, 85, 236, 22, 60, 134, 134, 179, 184, 159, 168, 83, 30, 239, 202, 238, 167, 49, 218, 251, 181, 224, 45, 30, 193, 245, 47, 214, 32, 37, 185, 178, 46, 201, 118, 119, 122, 30, 0, 57, 105, 58, 184, 15, 84, 166, 30, 182, 64, 211, 127, 9, 24, 216, 208, 172, 171, 31, 105, 63, 61, 209, 58, 12, 66, 232, 166, 4, 69, 237, 5, 36, 116, 32, 198, 64, 215, 83, 33, 131, 29, 100, 62, 41, 102, 180, 53, 233, 136, 153, 164, 1, 4, 230, 170, 175, 77, 16, 167, 85, 27, 149, 128, 225, 217, 160, 48, 173, 91, 217, 49, 26, 89, 204, 84, 243, 230, 212, 14, 177, 173, 9, 100, 143, 227, 88, 228, 200, 157, 64, 179, 109, 4, 34, 220, 158, 159, 170, 236, 4, 254, 1, 10, 101, 243, 220, 136, 223, 220, 182, 253, 67, 223, 168, 149, 116, 79, 90, 69, 216, 173, 48, 167, 48, 114, 244, 222, 54, 26, 76, 93, 246, 115, 185, 136, 212, 89, 136, 9, 101, 96, 110, 218, 92, 56, 187, 79, 204, 205, 208, 242, 213, 128, 91, 46, 107, 95, 133, 99, 128, 94, 178, 5, 179, 206, 98, 19, 138, 78, 120, 122, 252, 158, 179, 110, 247, 51, 136, 32, 57, 58, 65, 183, 63, 163, 140, 63, 61, 158, 8, 82, 150, 24, 16, 161, 199, 228, 141, 59, 77, 38, 7, 21, 7, 172, 128, 101, 159, 23, 63, 110, 191, 7, 255, 52, 80, 85, 19, 44, 5, 177, 109, 190, 140, 170, 35, 250, 203, 24, 57, 209, 83, 89, 58, 111, 203, 67, 52, 84, 33, 49, 85, 126, 53, 69, 238, 239, 244, 161, 106, 224, 171, 142, 148, 243, 187, 84, 62, 132, 129, 129, 29, 105, 232, 143, 107, 167, 137, 48, 128, 75, 101, 215, 45, 151, 190, 2, 6, 213, 8, 85, 234, 13, 112, 149, 67, 78, 133, 84, 5, 179, 165, 31, 161, 146, 109, 187, 221, 54, 74, 97, 61, 27, 180, 175, 255, 161, 48, 78, 238, 8, 9, 102, 30, 251, 171, 167, 85, 41, 33, 224, 111, 45, 243, 252, 137, 202, 139, 195, 50, 222, 22, 243, 193, 38, 63, 95, 125, 104, 79, 6, 158, 207, 107, 10, 3, 131, 243, 208, 141, 174, 247, 191, 51, 196, 173, 27, 38, 75, 104, 110, 193, 176, 126, 181, 14, 227, 45, 235, 179, 170, 236, 140, 35, 143, 138, 6, 140, 84, 245, 83, 253, 81, 139, 145, 195, 6, 161, 33, 6, 157, 138, 159, 196, 139, 138, 169, 5, 84, 165, 141, 91, 250, 126, 80, 160, 6, 197, 172, 29, 119, 111, 29, 248, 161, 222, 57, 114, 111, 91, 95, 237, 10, 242, 30, 253, 9, 232, 210, 146, 16, 111, 133, 83, 0, 227, 214, 88, 4, 174, 196, 17, 155, 41, 153, 234, 135, 83, 177, 53, 31, 51, 203, 252, 207, 45, 116, 117, 188, 51, 31, 68, 90, 27, 0, 45, 9, 123, 119, 153, 58, 39, 55, 206, 168, 46, 8, 64, 29, 86, 3, 196, 202, 255, 125, 124, 127, 109, 57, 143, 65, 31, 32, 48, 50, 210, 70, 206, 241, 156, 128, 160, 75, 17, 190, 248, 95, 44, 188, 5, 37, 151, 137, 204, 249, 231, 160, 29, 247, 192, 211, 103, 207, 233, 175, 169, 13, 12, 223, 84, 129, 177, 104, 167, 105, 21, 118, 222, 74, 111, 80, 138, 147, 164, 12, 39, 33, 141, 112, 138, 101, 38, 202, 121, 85, 152, 0, 166, 46, 149, 115, 44, 202, 225, 188, 188, 12, 139, 237, 148, 58},
		[]byte{185, 240, 38, 228, 92, 207, 38, 234, 13, 104, 2, 242, 83, 48, 147, 210, 210, 252, 86, 176, 110, 249, 97, 225, 247, 203, 155, 187, 155, 139, 138, 213, 140, 167, 48, 25, 18, 133, 90, 81, 212, 5, 203, 159, 247, 75, 141, 61, 199, 197, 194, 179, 197, 217, 199, 245, 173, 7, 219, 50, 123, 68, 171, 188, 175, 111, 128, 218, 100, 8, 167, 96, 134, 247, 245, 254, 159, 199, 205, 168, 222, 22, 23, 105, 9, 147, 193, 225, 4, 160, 74, 115, 224, 77, 66, 140, 151, 225, 223, 215, 40, 160, 29, 50, 65, 134, 134, 49, 38, 44, 152, 43, 228, 31, 7, 116, 159, 83, 239, 98, 106, 223, 190, 103, 70, 26, 4, 77, 173, 26, 217, 212, 112, 120, 246, 6, 154, 185, 147, 46, 253, 99, 209, 236, 201, 244, 61, 6, 115, 251, 94, 203, 96, 111, 195, 15, 172, 85, 149, 42, 183, 9, 193, 16, 29, 30, 120, 131, 85, 96, 137, 211, 68, 119, 3, 255, 44, 225, 74, 149, 216, 241, 10, 135, 245, 18, 69, 52, 83, 216, 254, 18, 177, 18, 185, 26, 95, 232, 244, 250, 195, 132, 133, 54, 8, 186, 26, 229, 83, 27, 233, 160, 186, 37, 213, 165, 67, 224, 97, 229, 160, 223, 85, 13, 92, 91, 24, 195, 227, 239, 63, 158, 78, 72, 216, 98, 251, 117, 194, 150, 63, 24, 25, 179, 213, 235, 164, 52, 152, 66, 43, 189, 255, 29, 198, 119, 142, 235, 90, 182, 249, 93, 174, 24, 33, 247, 251, 10, 179, 185, 14, 139, 162, 56, 248, 172, 86, 21, 232, 12, 65, 13, 169, 172, 253, 39, 14, 50, 162, 24, 128, 111, 169, 24, 98, 231, 49, 219, 23, 95, 214, 67, 238, 178, 63, 81, 249, 83, 180, 89, 13, 10, 36, 20, 226, 28, 149, 239, 14, 60, 45, 208, 204, 19, 104, 149, 85, 84, 213, 43, 249, 220, 13, 235, 228, 72, 252, 196, 200, 37, 69, 2, 193, 113, 107, 123, 51, 109, 152, 198, 251, 227, 95, 219, 224, 252, 36, 26, 134, 156, 129, 95, 15, 192, 52, 20, 212, 155, 188, 187, 93, 210, 29, 218, 126, 178, 100, 41, 148, 32, 36, 141, 186, 48, 25, 32, 230, 145, 147, 191, 44, 241, 166, 29, 227, 227, 81, 226, 210, 241, 23, 158, 238, 53, 29, 94, 53, 196, 62, 175, 126, 81, 200, 253, 36, 222, 125, 229, 39, 7, 113, 24, 19, 61, 106, 247, 24, 61, 159, 150, 87, 76, 230, 8, 156, 150, 149, 18, 60, 185, 72, 153, 61, 152, 55, 74, 178, 217, 178, 193, 14, 245, 63, 9, 151, 7, 72, 211, 5, 34, 214, 115, 218, 201, 145, 8, 188, 239, 197, 169, 179, 200, 2, 2, 0, 97, 248, 108, 157, 149, 104, 175, 2, 76, 167, 152, 11, 218, 125, 70, 140, 193, 123, 36, 209, 40, 184, 77, 23, 47, 210, 7, 105, 187, 36, 32, 168, 13, 145, 187, 103, 135, 60, 254, 19, 90, 138, 86, 90, 201, 178, 155, 78, 51, 105, 102, 84, 237, 144, 227, 60, 216, 40, 16, 71, 58, 161, 52, 178, 184, 34, 123, 67, 14, 42, 236, 67, 61, 250, 243, 81, 223, 128, 138, 176, 67, 180, 221, 36, 70, 140, 64, 124, 50, 253, 250, 247, 57, 172, 228, 36, 242, 3, 54, 65, 111, 86, 188, 176, 214, 111, 122, 143, 229, 200, 220, 108, 55, 195, 186, 18, 192, 111, 101, 151, 69, 207, 115, 253, 11, 182, 79, 60, 222, 3, 179, 206, 81, 118, 243, 254, 18, 106, 80, 64, 88, 21, 220, 255, 78, 79, 7, 111, 230, 227, 174, 59, 97, 209, 85, 142, 184, 9, 19, 195, 92, 159, 206, 16, 81, 93, 28, 8, 41, 38, 50, 30, 29, 251, 124, 186, 141, 151, 219, 199, 99, 2, 165, 201, 235, 89, 131, 184, 188, 241, 66, 90, 128, 36, 139, 253, 59, 47, 176, 161, 17, 64, 126, 72, 9, 93, 22, 47, 216, 69, 138, 115, 192, 220, 254, 139, 170, 213, 48, 85, 209, 177, 139, 153, 192, 117, 156, 191, 113, 127, 106, 79, 227, 237, 197, 94, 225, 168, 64, 49, 144, 42, 133, 71, 64, 238, 165, 197, 57, 53, 182, 40, 145, 189, 181, 114, 98, 231, 231, 43, 219, 15, 222, 67, 210, 208, 186, 234, 50, 13, 176, 43, 62, 88, 60, 200, 123, 22, 117, 225, 183, 37, 33, 16, 57, 91, 229, 87, 242, 35, 32, 150, 164, 180, 39, 2, 87, 205, 238, 139, 115, 77, 72, 134, 118, 219, 54, 110, 222, 122, 144, 106, 175, 236, 220, 37, 171, 243, 251, 197, 149, 6, 100, 109, 232, 195, 141, 141, 12, 44, 173, 1, 79, 232, 117, 252, 49, 26, 188, 203, 83, 153, 111, 209, 244, 131, 26, 12, 252, 3, 100, 254, 108, 209, 22, 5, 185, 6, 144, 156, 184, 227, 53, 26, 113, 46, 166, 158, 74, 250, 212, 59, 216, 78, 99, 124, 163, 87, 243, 141, 172, 54, 210, 149, 23, 169, 123, 15, 2, 131, 140, 188, 158, 9, 95, 158, 107, 2, 160, 80, 237, 24, 34, 185, 216, 150, 18, 36, 45, 90, 45, 190, 26, 242, 171, 45, 141, 175, 44, 87, 157, 247, 156, 143, 113, 17, 191, 120, 234, 237, 127, 86, 157, 86, 169, 209, 8, 44, 107, 173, 170, 85, 198, 18, 227, 243, 250, 221, 224, 31, 222, 218, 18, 206, 255, 93, 241, 81, 47, 153, 0, 186, 67, 11, 142, 92, 124, 144, 219, 6, 96, 144, 41, 122, 43, 39, 232, 137, 221, 92, 96, 17, 240, 84, 222, 168, 132, 109, 19, 115, 99, 222, 87, 12, 109, 7, 74, 16, 60, 58, 237, 47, 93, 114, 50, 183, 151, 185, 213, 143, 111, 207, 63, 21, 136, 122, 27, 218, 94, 80, 251, 155, 106, 127, 133, 6},
	)
}

func AddCustomImplementationFuzzCase(f *testing.F, initialSlice []byte, opData []byte) {
	f.Add(initialSlice, opData)
}

// PARAMS:
//   - ```f *testing.F```: the fuzzing object from the testing package
//   - ```implementationTypeName string```: a string name describing the type (will print in errors)
//   - ```initializerFunc func(initialData []byte) L```: a function that initializes the ListLike with the data from the provided slice as an initial state
//   - ```cleanupFunc func(t *testing.T, list L)```: a function that performs any necessary cleanup on the ListLike (for example if a file was created for testing)
func PerformListImplementationFuzz[IDX LL.Integer, L LL.ListLike[byte, IDX]](f *testing.F, implementationTypeName string, initializerFunc func(t *testing.T, initialData []byte) L, cleanupFunc func(t *testing.T, list L)) {
	f.Fuzz(func(t *testing.T, initialSlice []byte, opData []byte) {
		var op byte
		list := initializerFunc(t, initialSlice)
		defer cleanupFunc(t, list)
		expect := initialSlice
		var keepGoing bool = true
		for keepGoing {
			op, keepGoing = getOneOpDataVal(&opData)
			if !keepGoing {
				break
			}
			keepGoing = doListOp(t, implementationTypeName, &expect, list, op, &opData)
			if !keepGoing {
				break
			}
		}
	})
}

func logImplError[IDX LL.Integer, S LL.SliceLike[byte, IDX]](t *testing.T, typeName string, expSlice *[]byte, list S, reason string, op byte, opArgs ...any) {
	listSlice := make([]byte, 0, list.Len())
	LL.DoActionOnAllItems(list, func(slice S, idx IDX, item byte) {
		listSlice = append(listSlice, item)
	})
	if len(opArgs) > 1 {
		t.Errorf("go_list_like: implementation test: operation `%s.%s(%s)` failed test:\n\tREASON: %s\n\tEXP: %v\n\tGOT: %v\n", typeName, _OpNames[op], stringifyParamList(opArgs...), reason, *expSlice, listSlice)
	} else {
		t.Errorf("go_list_like: implementation test: operation `%s.%s()` failed test:\n\tREASON: %s\n\tEXP: %v\n\tGOT: %v\n", typeName, _OpNames[op], reason, *expSlice, listSlice)
	}
}

func checkParity[IDX LL.Integer, S LL.SliceLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list S, op byte, opArgs ...any) (ok bool) {
	ok = true
	if IDX(len(*goSlice)) != list.Len() {
		ok = false
		logImplError(t, typeName, goSlice, list, fmt.Sprintf("expected length (%d) does not equal reported length (%d)", len(*goSlice), list.Len()), op, opArgs...)
		return
	}
	if list.LenBetween(list.FirstIdx(), list.LastIdx()) != list.Len() {
		ok = false
		logImplError(t, typeName, goSlice, list, fmt.Sprintf("reported length (%d) does not equal reported length-between first and last indexes (%d)", list.Len(), list.LenBetween(list.FirstIdx(), list.LastIdx())), op, opArgs...)
		return
	}
	llIdx := list.FirstIdx()
	llOk := list.IdxValid(llIdx)
	goIdx := 0
	goOk := goIdx < len(*goSlice)
	var goVal, llVal byte
	for {
		if !llOk && !goOk {
			break
		}
		if !llOk && goOk {
			ok = false
			logImplError(t, typeName, goSlice, list, "ListLike implementation returned an invalid idx where a valid index was expected", op, opArgs...)
			break
		}
		if llOk && !goOk {
			ok = false
			logImplError(t, typeName, goSlice, list, "ListLike implementation returned a valid idx where an invalid index was expected", op, opArgs...)
			break
		}
		goVal = (*goSlice)[goIdx]
		llVal = list.Get(llIdx)
		if goVal != llVal {
			ok = false
			logImplError(t, typeName, goSlice, list, fmt.Sprintf("value at index `%d` places from beginning (idx %d) did not match the expected value, EXP %d, GOT %d", goIdx, llIdx, goVal, llVal), op, opArgs...)
			break
		}
		goIdx += 1
		goOk = goIdx < len(*goSlice)
		llIdx = list.NextIdx(llIdx)
		llOk = list.IdxValid(llIdx)
	}
	return
}
func getOneOpIndex(opData *[]byte) (index int, hasData bool) {
	if len(*opData) <= 0 {
		return
	}
	hasData = true
	index = int((*opData)[0])
	if len(*opData) > 1 {
		index |= int((*opData)[1]) << 8
		*opData = (*opData)[2:]
	} else {
		*opData = (*opData)[1:]
	}
	return
}
func getOneOpDataVal(opData *[]byte) (val byte, hasData bool) {
	if len(*opData) <= 0 {
		return
	}
	hasData = true
	val = (*opData)[0]
	*opData = (*opData)[1:]
	return
}
func getManyOpDataVals(opData *[]byte) (vals []byte, hasData bool) {
	if len(*opData) < 2 {
		return
	}
	n, _ := getOneOpDataVal(opData)
	n = max(min(n, _MaxValuesPerOp, byte(len(*opData))), 1)
	vals = (*opData)[:n]
	*opData = (*opData)[n:]
	hasData = len(vals) > 0
	return
}

func checkIdxParity[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, nthIdx int, op byte, llIdx IDX, llValid bool) (valid bool) {
	var goVal, llVal byte
	if nthIdx >= len(*goSlice) || nthIdx < 0 {
		if llValid {
			logImplError(t, typeName, goSlice, list, fmt.Sprintf("index that is %d positions from beginning (%d) should have been an invalid index, but wasn't (len == %d)", nthIdx, llIdx, len(*goSlice)), op, llIdx)
		}
		valid = false
	} else {
		if !llValid {
			logImplError(t, typeName, goSlice, list, fmt.Sprintf("index that is %d positions from beginning (%d) should have been a valid index, but wasn't (len == %d)", nthIdx, llIdx, len(*goSlice)), op, llIdx)
		} else {
			goVal = (*goSlice)[nthIdx]
			llVal = list.Get(llIdx)
			valid = goVal == llVal
			if !valid {
				logImplError(t, typeName, goSlice, list, fmt.Sprintf("index that is %d positions from beginning (%d) did not have the same value as the one %d positions from the start of the expected slice, EXP %d, GOT %d", nthIdx, llIdx, nthIdx, goVal, llVal), op, llIdx)
			}
		}
	}
	return
}

func checkRangeAndSwapIfNeeded[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, nthA int, nthB int, llNthA, llNthB IDX, op byte) (nthAA, nthBB int, llNthAA, llNthBB IDX, valid bool) {
	valid = true
	if nthA > nthB {
		wasInvalid := !list.RangeValid(llNthA, llNthB)
		if !wasInvalid {
			valid = false
			logImplError(t, typeName, goSlice, list, fmt.Sprintf("range (nth-position) [%d, %d] (true index [%d, %d]) should have reported as an invalid range but reported valid", nthA, nthB, llNthA, llNthB), op, llNthA, llNthB)
		}
		nthAA = nthB
		nthBB = nthA
		llNthAA = llNthB
		llNthBB = llNthA
	} else {
		nthAA = nthA
		nthBB = nthB
		llNthAA = llNthA
		llNthBB = llNthB
	}
	if nthA < 0 || int(nthB) >= len(*goSlice) {
		wasInvalid := !list.RangeValid(llNthA, llNthB)
		if !wasInvalid {
			valid = false
			logImplError(t, typeName, goSlice, list, fmt.Sprintf("range (nth-position) [%d, %d] (true index [%d, %d]) should have reported as an invalid range but reported valid", nthA, nthB, llNthA, llNthB), op, llNthA, llNthB)
		}
	} else {
		wasValid := list.RangeValid(llNthAA, llNthBB)
		if !wasValid {
			valid = false
			logImplError(t, typeName, goSlice, list, fmt.Sprintf("range (nth-position) [%d, %d] (true index [%d, %d]) should have reported as a valid range but reported invalid", nthA, nthB, llNthA, llNthB), op, llNthA, llNthB)
		}
	}
	return
}

func checkAndGetNthIdx[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, nthIdx int, op byte, opData *[]byte) (idx IDX, valid bool) {
	idx, valid = LL.TryNthIdx(list, IDX(nthIdx))
	valid = valid && checkIdxParity(t, typeName, goSlice, list, int(nthIdx), op, idx, valid)
	return
}

func checkAndGetFirstIdx[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, op byte, opData *[]byte) (idx IDX, valid bool) {
	idx, valid = LL.TryFirstIdx(list)
	valid = valid && checkIdxParity(t, typeName, goSlice, list, 0, op, idx, valid)
	return
}
func checkAndGetLastIdx[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, op byte, opData *[]byte) (idx IDX, valid bool) {
	idx, valid = LL.TryLastIdx(list)
	valid = valid && checkIdxParity(t, typeName, goSlice, list, len(*goSlice)-1, op, idx, valid)
	return
}

func checkAndGetPrevIdx[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, startIdx int, op byte, opData *[]byte) (idx IDX, valid bool) {
	idx, valid = LL.TryNthIdx(list, IDX(startIdx))
	valid = valid && checkIdxParity(t, typeName, goSlice, list, int(startIdx), op, idx, valid)
	if !valid {
		return
	}
	idx, valid = LL.TryPrevIdx(list, idx)
	valid = valid && checkIdxParity(t, typeName, goSlice, list, int(startIdx)-1, op, idx, valid)
	return
}
func checkAndGetNextIdx[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, startIdx int, op byte, opData *[]byte) (idx IDX, valid bool) {
	idx, valid = LL.TryNthIdx(list, IDX(startIdx))
	valid = valid && checkIdxParity(t, typeName, goSlice, list, int(startIdx), op, idx, valid)
	if !valid {
		return
	}
	idx, valid = LL.TryNextIdx(list, idx)
	valid = valid && checkIdxParity(t, typeName, goSlice, list, int(startIdx)+1, op, idx, valid)
	return
}
func checkAndGetNthPrevIdx[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, startIdx int, n int, op byte, opData *[]byte) (idx IDX, valid bool) {
	idx, valid = LL.TryNthIdx(list, IDX(startIdx))
	valid = valid && checkIdxParity(t, typeName, goSlice, list, int(startIdx), op, idx, valid)
	if !valid {
		return
	}
	idx, valid = LL.TryNthPrevIdx(list, idx, IDX(n))
	valid = valid && checkIdxParity(t, typeName, goSlice, list, int(startIdx-n), op, idx, valid)
	return
}
func checkAndGetNthNextIdx[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, startIdx int, n int, op byte, opData *[]byte) (idx IDX, valid bool) {
	idx, valid = LL.TryNthIdx(list, IDX(startIdx))
	valid = valid && checkIdxParity(t, typeName, goSlice, list, int(startIdx), op, idx, valid)
	if !valid {
		return
	}
	idx, valid = LL.TryNthNextIdx(list, idx, IDX(n))
	valid = valid && checkIdxParity(t, typeName, goSlice, list, int(startIdx+n), op, idx, valid)
	return
}

func doListOp[IDX LL.Integer, L LL.ListLike[byte, IDX]](t *testing.T, typeName string, goSlice *[]byte, list L, op byte, opData *[]byte) (couldDo bool) {
	op = op % _opCount
	var nthIdxA, nthIdxB, nthIdxC, nthIdxD int
	var llNthIdxA, llNthIdxB, llNthIdxC, llNthIdxD IDX
	var setVal byte
	var insertVals []byte
	switch op {
	case _ImplOpGet:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		_, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxA, op, opData)
		if !couldDo {
			return
		}
	case _ImplOpSet:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		llNthIdxA, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxA, op, opData)
		if !couldDo {
			return
		}
		setVal, couldDo = getOneOpDataVal(opData)
		if !couldDo {
			return
		}
		list.Set(llNthIdxA, setVal)
		(*goSlice)[nthIdxA] = setVal
		couldDo = checkIdxParity(t, typeName, goSlice, list, int(nthIdxA), op, llNthIdxA, true)
		if !couldDo {
			return
		}
	case _ImplOpMove:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		nthIdxB, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		llNthIdxA, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxA, op, opData)
		if !couldDo {
			return
		}
		llNthIdxB, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxB, op, opData)
		if !couldDo {
			return
		}
		list.Move(llNthIdxA, llNthIdxB)
		sliceMoveOne(goSlice, nthIdxA, nthIdxB)
		couldDo = checkParity(t, typeName, goSlice, list, op, llNthIdxA, llNthIdxB)
		if !couldDo {
			return
		}
	case _ImplOpMoveRange:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		nthIdxB, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		nthIdxA, nthIdxB, couldDo = clampRange(goSlice, nthIdxA, nthIdxB)
		if !couldDo {
			return
		}
		llNthIdxA, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxA, op, opData)
		if !couldDo {
			return
		}

		llNthIdxB, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxB, op, opData)
		if !couldDo {
			return
		}
		nthIdxC, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		llNthIdxC, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxC, op, opData)
		if !couldDo {
			return
		}
		nthIdxD = nthIdxC + (nthIdxB - nthIdxA)
		llNthIdxD, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxD, op, opData)
		couldDo = couldDo && list.IdxValid(llNthIdxD)
		if !couldDo {
			return
		}
		nthIdxA, nthIdxB, llNthIdxA, llNthIdxB, couldDo = checkRangeAndSwapIfNeeded(t, typeName, goSlice, list, nthIdxA, nthIdxB, llNthIdxA, llNthIdxB, op)
		if !couldDo {
			return
		}
		list.MoveRange(llNthIdxA, llNthIdxB, llNthIdxC)
		sliceMoveRange(goSlice, nthIdxA, nthIdxB, nthIdxC)
		couldDo = checkParity(t, typeName, goSlice, list, op, llNthIdxA, llNthIdxB, llNthIdxC)
		if !couldDo {
			return
		}
	case _ImplOpSlice:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		llNthIdxA, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxA, op, opData)
		if !couldDo {
			return
		}
		nthIdxB, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		llNthIdxB, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxB, op, opData)
		if !couldDo {
			return
		}
		nthIdxA, nthIdxB, llNthIdxA, llNthIdxB, couldDo = checkRangeAndSwapIfNeeded(t, typeName, goSlice, list, nthIdxA, nthIdxB, llNthIdxA, llNthIdxB, op)
		if !couldDo {
			return
		}
		subList := list.Slice(llNthIdxA, llNthIdxB)
		subSlice := (*goSlice)[nthIdxA : nthIdxB+1]
		couldDo = checkParity(t, typeName, &subSlice, subList, op, llNthIdxA, llNthIdxB)
		if !couldDo {
			return
		}
	case _ImplOpFirst:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		_, couldDo = checkAndGetFirstIdx(t, typeName, goSlice, list, op, opData)
		if !couldDo {
			return
		}
	case _ImplOpLast:
		_, couldDo = checkAndGetLastIdx(t, typeName, goSlice, list, op, opData)
		if !couldDo {
			return
		}
	case _ImplOpNext:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		_, couldDo = checkAndGetNextIdx(t, typeName, goSlice, list, nthIdxA, op, opData)
		if !couldDo {
			return
		}
	case _ImplOpNthNext:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		nthIdxB, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		_, couldDo = checkAndGetNthNextIdx(t, typeName, goSlice, list, nthIdxA, nthIdxB, op, opData)
		if !couldDo {
			return
		}
	case _ImplOpPrev:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		_, couldDo = checkAndGetPrevIdx(t, typeName, goSlice, list, nthIdxA, op, opData)
		if !couldDo {
			return
		}
	case _ImplOpNthPrev:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		nthIdxB, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		_, couldDo = checkAndGetNthPrevIdx(t, typeName, goSlice, list, nthIdxA, nthIdxB, op, opData)
		if !couldDo {
			return
		}
	case _ImplOpInsert:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		llNthIdxA, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxA, op, opData)
		if !couldDo {
			return
		}
		insertVals, couldDo = getManyOpDataVals(opData)
		if !couldDo {
			return
		}
		llInsertVals := LL.NewSliceAdapter(insertVals)
		gotFree := list.TryEnsureFreeSlots(IDX(len(insertVals)))
		if !gotFree {
			couldDo = true
			return
		}
		llNthIdxB, llNthIdxC = list.InsertSlotsAssumeCapacity(llNthIdxA, IDX(len(insertVals)))
		LL.CopyToRange(&llInsertVals, list, llNthIdxB, llNthIdxC)

		*goSlice = slices.Insert(*goSlice, int(nthIdxA), insertVals...)
		couldDo = checkParity(t, typeName, goSlice, list, op, llNthIdxA, IDX(len(insertVals)), insertVals)
		if !couldDo {
			return
		}
	case _ImplOpAppend:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		insertVals, couldDo = getManyOpDataVals(opData)
		if !couldDo {
			return
		}
		llInsertVals := LL.NewSliceAdapter(insertVals)
		gotFree := list.TryEnsureFreeSlots(IDX(len(insertVals)))
		if !gotFree {
			couldDo = true
			return
		}
		llNthIdxB, llNthIdxC = list.AppendSlotsAssumeCapacity(IDX(len(insertVals)))
		LL.CopyToRange(&llInsertVals, list, llNthIdxB, llNthIdxC)
		*goSlice = append(*goSlice, insertVals...)
		couldDo = checkParity(t, typeName, goSlice, list, op, llNthIdxA, insertVals)
		if !couldDo {
			return
		}
	case _ImplOpDelete:
		if EnableBreakpoints {
			runtime.Breakpoint()
		}
		if len(*goSlice) == 0 {
			couldDo = true
			return
		}
		nthIdxA, couldDo = getOneOpIndex(opData)
		if !couldDo {
			return
		}
		llNthIdxA, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxA, op, opData)
		if !couldDo {
			return
		}
		setVal, couldDo = getOneOpDataVal(opData)
		if !couldDo {
			return
		}
		setVal = max(setVal%_MaxValuesPerOp, 1)
		nthIdxB = nthIdxA + int(setVal)
		nthIdxB = min(nthIdxB, int(len(*goSlice)-1))
		if nthIdxA == nthIdxB {
			couldDo = true
			return
		}
		llNthIdxB, couldDo = checkAndGetNthIdx(t, typeName, goSlice, list, nthIdxB, op, opData)
		if !couldDo {
			return
		}
		list.DeleteRange(llNthIdxA, llNthIdxB)

		*goSlice = slices.Delete(*goSlice, int(nthIdxA), int(nthIdxB)+1)
		couldDo = checkParity(t, typeName, goSlice, list, op, llNthIdxA, IDX(len(insertVals)), insertVals)
		if !couldDo {
			return
		}
	default:
		panic(fmt.Sprintf("somehow, an invalid op-code was found (%d, max valid == %d)", op, _opCount-1))
	}
	couldDo = true
	return
}

func sliceMoveOne(slice *[]byte, old, new int) {
	val := (*slice)[old]
	if new < old {
		i := old
		ii := old - 1
		for i > new {
			(*slice)[i] = (*slice)[ii]
			i = ii
			ii -= 1
		}
	} else {
		i := old
		ii := old + 1
		for i < new {
			(*slice)[i] = (*slice)[ii]
			i = ii
			ii += 1
		}
	}
	(*slice)[new] = val
}

func sliceMoveRange(slice *[]byte, firstOld, lastOld, firstNew int) {
	lenA := (lastOld - firstOld) + 1
	sliceA := (*slice)[firstOld : firstOld+lenA]
	var totalRange, sliceB []byte
	if firstNew < firstOld {
		totalRange = (*slice)[firstNew : lastOld+1]
		sliceB = (*slice)[firstNew:firstOld]
	} else {
		totalRange = (*slice)[firstOld : firstNew+lenA]
		sliceB = (*slice)[lastOld+1 : firstNew+lenA]
	}
	slices.Reverse(sliceA)
	slices.Reverse(sliceB)
	slices.Reverse(totalRange)
}

func stringifyParamList(args ...any) string {
	builder := strings.Builder{}
	for _, a := range args {
		builder.WriteString(fmt.Sprintf("%v, ", a))
	}
	result := builder.String()
	return result[:len(result)-2]
}

func clampRange(slice *[]byte, first, last int) (ff, ll int, ok bool) {
	ok = true
	ff = first
	ll = last
	if ff > ll {
		ff = last
		ll = first
	}
	cc := ll - ff
	cc = min(cc, int(_MaxValuesPerOp))
	ll = ff + cc
	if int(ll) >= len(*slice) {
		ll = int(len(*slice) - 1)
	}
	if ff > ll || int(ll) >= len(*slice) {
		ok = false
	}
	return
}
