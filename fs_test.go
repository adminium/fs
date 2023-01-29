package fs

import (
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
	a := ParseCrossPath("/pkg/mod/github.com/polydawn/refmt", "github.com/polydawn/refmt/shared")
	if a != "/pkg/mod/github.com/polydawn/refmt/shared" {
		t.Fatal("match error: ", a)
	}
}
