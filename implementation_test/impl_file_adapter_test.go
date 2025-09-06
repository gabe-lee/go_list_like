package implementation_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	LL "github.com/gabe-lee/go_list_like"
)

func newFileAdapter(t *testing.T, data []byte) LL.FileAdapter {
	fname := t.Name()
	fname = strings.ReplaceAll(fname, "/", ".")
	file, err := os.CreateTemp("", fname)
	if err != nil {
		panic(fmt.Sprintf("could not create temp file '%s': %s", fname, err))
	}
	file.Truncate(int64(len(data)))
	file.WriteAt(data, 0)
	return LL.NewFileAdapter(file)
}

func cleanupFileAdapter(t *testing.T, fa LL.FileAdapter) {
	fa.Close()
	os.Remove(fa.Name())
}

func Fuzz_FileAdapter_(f *testing.F) {
	InitImplementationFuzz(f)
	PerformListImplementationFuzz(f, "FileAdapter", newFileAdapter, cleanupFileAdapter)
}
