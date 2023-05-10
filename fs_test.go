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
	t.Log(Join("a", "b"))
	t.Log(Join("/a", "b", "c"))
	t.Log(Join("/a", "aa", "/b", "c"))
}

func TestLookupFrom(t *testing.T) {
	a := ParseCrossExtra("/pkg/mod/github.com/polydawn/refmt", "github.com/polydawn/refmt/shared")
	if a != "/shared" {
		t.Fatal("match error: ", a)
	}
}

func TestSizeOf(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)
	
	t.Log(pwd)
	s, err := SizeOf(pwd)
	require.NoError(t, err)
	
	t.Log(humanize.Bytes(s))
}
