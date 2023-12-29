package fs

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Write content to a file
func Write(path string, content []byte) (err error) {
	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return
	}
	return
}

// Append content to a file
func Append(path string, content []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

// Read a file content
func Read(path string) (data []byte, err error) {
	data, err = os.ReadFile(path)
	if err != nil {
		return
	}
	return
}

// Scan a file content with progress
func Scan(path string, reader func(line []byte, percent float32) error) (err error) {
	if reader == nil {
		err = fmt.Errorf("reader is nil")
		return
	}
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()
	stat, err := os.Stat(path)
	if err != nil {
		return
	}
	totalSize := stat.Size()
	scanSize := 0
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1073741824)
	for scanner.Scan() {
		scanSize += len(scanner.Bytes()) + 1
		percent := float32(scanSize) / float32(totalSize)
		err = reader(scanner.Bytes(), percent)
		if err != nil {
			err = fmt.Errorf("call reader error: %s", err)
			return
		}
	}
	err = scanner.Err()
	if err != nil {
		return
	}
	return
}

// Remove a file or directory
func Remove(path string) error {
	return os.RemoveAll(path)
}

// Exist check a path exist
func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Lookup Start searching for the object from the specified path,
// and if not found, fall back to the previous directory to continue searching.
func Lookup(path, lookup string) (string, error) {
	origin := path
	var err error
	for {
		p := filepath.Join(path, lookup)
		_, err = os.Stat(p)
		if err == nil {
			return p, nil
		}
		if path == "." {
			break
		}
		path = filepath.Join(path, "../")
	}

	return "", fmt.Errorf("lookup path: %s from: %s faild", lookup, origin)
}

// IsFile check path is a file or not
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IsDir check path is a directory or not
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Files get files paths
func Files(deeply bool, dir string, suffix ...string) (files []string, err error) {

	d, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	sep := string(os.PathSeparator)

	for _, fi := range d {
		//if !fi.IsDir() {
		//	if HasSuffix(fi.Name(), suffix...) {
		//		files = append(files, filepath.Join(dir, sep, fi.Name()))
		//	}
		//}
		if deeply && fi.IsDir() {
			var dirFiles []string
			dirFiles, err = Files(deeply, filepath.Join(dir, sep, fi.Name()), suffix...)
			if err != nil {
				return
			}
			files = append(files, dirFiles...)
		}

		if !fi.IsDir() {
			var r []string
			r, err = Files(deeply, filepath.Join(dir, sep, fi.Name()), suffix...)
			if err != nil {
				return
			}
			files = append(files, r...)
		}
	}

	return
}

// HasSuffix check path has suffix or not
func HasSuffix(path string, suffix ...string) bool {
	for _, v := range suffix {
		if strings.HasSuffix(path, v) {
			return true
		}
	}
	return false
}

// MakeDir create a directory
func MakeDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// RealName get a file real name
func RealName(path string) (string, error) {
	stat, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	switch stat.Mode().Type() {
	case fs.ModeSymlink:
		var name string
		name, err = os.Readlink(path)
		if err != nil {
			return "", err
		}
		return name, nil
	}

	return filepath.Base(path), nil
}

// Join
// /a/b/c, b,d  => /a/b/d
// /a/b/c, d    => /a/b/c/d
// /a/b/c, /d   => /d
func Join(elem ...string) string {
	l := len(elem)
	if l == 0 {
		return ""
	}
	if l == 1 {
		return elem[0]
	}
	p := elem[0]
	for i := 1; i < l; i++ {
		p = join(p, elem[i])
	}
	return p
}

func join(path, lookup string) string {
	a := strings.Split(path, string(os.PathSeparator))
	b := strings.Split(lookup, string(os.PathSeparator))
	i := len(a) - 1
	if len(b) > 0 {
		for {
			if i < 0 {
				break
			}
			if a[i] == b[0] {
				r := filepath.Join(append(a[:i], b...)...)
				if strings.HasPrefix(path, string(os.PathSeparator)) {
					r = fmt.Sprintf("%c%s", os.PathSeparator, r)
				}
				return r
			}
			i--
		}
	}

	return filepath.Join(path, lookup)
}

// TrimCrossPrefix
// a/b/c/d  q/b/d/c  =>  c
func TrimCrossPrefix(a, b string) string {
	items := strings.Split(b, "/")
	if len(items) == 0 {
		return ""
	}
	for i := len(items); i >= 0; i-- {
		c := strings.Join(items[0:i], "/")
		if index := strings.Index(a, c); index > -1 {
			return strings.TrimPrefix(b, c)
		}
	}
	return ""
}

// SizeOf Calc a file or directory size
func SizeOf(path string) (size uint64, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size += uint64(info.Size())
		return nil
	})
	if err != nil {
		return
	}
	return
}
