package go_list_like

import (
	"io"
	"os"
	"syscall"
	"time"
)

type FileAdapter struct {
	File os.File
}

func (f FileAdapter) Cap() int {
	stat, err := f.File.Stat()
	if err != nil {
		return 0
	}
	return int(stat.Size())
}

func (f FileAdapter) ChangeLen(delta int) {
	f.File.Truncate(int64(f.Len() + delta))
}

func (f FileAdapter) Get(idx int) (val byte) {
	var b [1]byte
	f.File.ReadAt(b[:], int64(idx))
	return b[0]
}

func (f FileAdapter) GrowCap(n int) {}

func (f FileAdapter) Len() int {
	stat, err := f.File.Stat()
	if err != nil {
		return 0
	}
	return int(stat.Size())
}

func (f FileAdapter) Set(idx int, val byte) {
	var b = [1]byte{val}
	f.File.WriteAt(b[:], int64(idx))
}

func (f FileAdapter) Chdir() error {
	return f.File.Chdir()
}
func (f FileAdapter) Chmod(mode os.FileMode) error {
	return f.File.Chmod(mode)
}
func (f FileAdapter) Chown(uid int, gid int) error {
	return f.File.Chown(uid, gid)
}
func (f FileAdapter) Close() error {
	return f.File.Close()
}
func (f FileAdapter) Fd() uintptr {
	return f.File.Fd()
}
func (f FileAdapter) Name() string {
	return f.File.Name()
}
func (f FileAdapter) Read(b []byte) (n int, err error) {
	return f.File.Read(b)
}
func (f FileAdapter) ReadAt(b []byte, off int64) (n int, err error) {
	return f.File.ReadAt(b, off)
}
func (f FileAdapter) ReadDir(n int) ([]os.DirEntry, error) {
	return f.File.ReadDir(n)
}
func (f FileAdapter) ReadFrom(r io.Reader) (n int64, err error) {
	return f.File.ReadFrom(r)
}
func (f FileAdapter) Readdir(n int) ([]os.FileInfo, error) {
	return f.File.Readdir(n)
}
func (f FileAdapter) Readdirnames(n int) (names []string, err error) {
	return f.File.Readdirnames(n)
}
func (f FileAdapter) Seek(offset int64, whence int) (ret int64, err error) {
	return f.File.Seek(offset, whence)
}
func (f FileAdapter) SetDeadline(t time.Time) error {
	return f.File.SetDeadline(t)
}
func (f FileAdapter) SetReadDeadline(t time.Time) error {
	return f.File.SetReadDeadline(t)
}
func (f FileAdapter) SetWriteDeadline(t time.Time) error {
	return f.File.SetWriteDeadline(t)
}
func (f FileAdapter) Stat() (os.FileInfo, error) {
	return f.File.Stat()
}
func (f FileAdapter) Sync() error {
	return f.File.Sync()
}
func (f FileAdapter) SyscallConn() (syscall.RawConn, error) {
	return f.File.SyscallConn()
}
func (f FileAdapter) Truncate(size int64) error {
	return f.File.Truncate(size)
}
func (f FileAdapter) Write(b []byte) (n int, err error) {
	return f.File.Write(b)
}
func (f FileAdapter) WriteAt(b []byte, off int64) (n int, err error) {
	return f.File.WriteAt(b, off)
}
func (f FileAdapter) WriteString(s string) (n int, err error) {
	return f.File.WriteString(s)
}
func (f FileAdapter) WriteTo(w io.Writer) (n int64, err error) {
	return f.File.WriteTo(w)
}

var _ ListLike[byte] = FileAdapter{}
var _ io.Reader = FileAdapter{}
var _ io.Writer = FileAdapter{}
var _ io.ReaderAt = FileAdapter{}
var _ io.WriterAt = FileAdapter{}
