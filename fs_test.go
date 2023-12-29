package fs

import (
	"os"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/require"
)

func TestMakeDir(t *testing.T) {
	err := MakeDir("./test/a/b/c")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLookupJoin(t *testing.T) {
	t.Log(Join("a/b/c.txt", "a/b/d"))
	t.Log(Join("/a/b/c/e", "b/d"))
	t.Log(Join("/a/b/c", "d"))
	t.Log(Join("/a/b/c", "/d"))
	t.Log(Join("/a", "/b", "c"))
}

func TestTrimCrossPrefix(t *testing.T) {
	a := TrimCrossPrefix("/pkg/mod/github.com/polydawn/refmt", "github.com/polydawn/refmt/shared")
	require.Equal(t, "/shared", a)
}

func TestSizeOf(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	t.Log(pwd)
	s, err := SizeOf(pwd)
	require.NoError(t, err)

	t.Log(humanize.Bytes(s))
}

func TestScan(t *testing.T) {
	err := Scan("./test/test.sql", func(line []byte, percent float32) error {
		t.Log(string(line), percent)
		return nil
	})
	require.NoError(t, err)
}
