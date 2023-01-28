package fs

import "testing"

func TestMakeDir(t *testing.T) {
	err := MakeDir("./test/a/b/c")
	if err != nil {
		t.Fatal(err)
	}
}
