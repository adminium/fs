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

func TestJoin(t *testing.T) {

	type Test struct {
		Elems  []string
		Expect string
	}

	tests := []Test{
		{Elems: []string{"a/b/c.txt", "a/b/d"}, Expect: "a/b/d"},
		{Elems: []string{"/a/b/c/e", "b/d"}, Expect: "/a/b/d"},
		{Elems: []string{"/a/b/c", "d"}, Expect: "/a/b/c/d"},
		{Elems: []string{"/a/b/c", "/d"}, Expect: "/d"},
		{Elems: []string{"/a", "/b", "c"}, Expect: "/b/c"},
	}

	for _, v := range tests {
		require.Equal(t, v.Expect, Join(v.Elems...), v.Elems)
	}
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

func TestFiles(t *testing.T) {
	require.NoError(t, MakeDir("tests"))
	require.NoError(t, Write("tests/a.txt", []byte("a")))
	require.NoError(t, MakeDir("tests/b"))
	require.NoError(t, Write("tests/b/b.txt", []byte("b")))

	files, err := Files(true, "tests")
	require.NoError(t, err)

	require.Equal(t, "tests/a.txt", files[0])
	require.Equal(t, "tests/b/b.txt", files[1])

	require.NoError(t, Remove("tests"))
}
