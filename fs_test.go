package fs

import (
	"github.com/gozelle/humanize"
	"github.com/gozelle/testify/require"
	"os"
	"testing"
)

func TestMakeDir(t *testing.T) {
	err := MakeDir("./test/a/b/c")
	if err != nil {
		t.Fatal(err)
	}
}

func TestJoin(t *testing.T) {
	t.Log(MergeJoin("a", "b"))
	t.Log(MergeJoin("/a/b/c", "b/d"))
	t.Log(MergeJoin("a", "/b", "c"))
	t.Log(MergeJoin("/a", "aa", "/b", "c"))
	t.Log(MergeJoin("/a", "aa", "//b", "c"))
}

func TestLookupJoin(t *testing.T) {
	t.Log(LookupJoin("a/b/c.txt", "a/b/d"))
	t.Log(LookupJoin("/a/b/c/e", "b/d"))
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
	err := Scan("./test/test.sql", func(line []byte) error {
		t.Log(string(line))
		return nil
	})
	require.NoError(t, err)
}
