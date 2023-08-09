package fs

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Write(path string, content []byte) (err error) {
	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return
	}
	return
}

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

func Read(path string) (data []byte, err error) {
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return
}

func Scan(path string, reader func(line []byte) error) (err error) {
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = reader(scanner.Bytes())
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

func Remove(path string) error {
	return os.RemoveAll(path)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func AllFiles(dir string, suffix ...string) (files []string, err error) {
	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	sep := string(os.PathSeparator)
	for _, fi := range d {
		if fi.IsDir() {
			var dirFiles []string
			dirFiles, err = Files(filepath.Join(dir, sep, fi.Name()), suffix...)
			if err != nil {
				return
			}
			files = append(files, dirFiles...)
		} else {
			var r []string
			r, err = Files(filepath.Join(dir, sep, fi.Name()), suffix...)
			if err != nil {
				return
			}
			files = append(files, r...)
		}
	}
	return
}

func Files(dir string, suffix ...string) (files []string, err error) {
	
	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	
	sep := string(os.PathSeparator)
	
	for _, fi := range d {
		if !fi.IsDir() {
			if HasSuffix(fi.Name(), suffix...) {
				files = append(files, filepath.Join(dir, sep, fi.Name()))
			}
		}
	}
	
	return
}

func HasSuffix(file string, suffix ...string) bool {
	for _, v := range suffix {
		if strings.HasSuffix(file, v) {
			return true
		}
	}
	return false
}

func MakeDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

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

func Lookup(lookup string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		p := filepath.Join(pwd, lookup)
		_, err = os.Stat(p)
		if err == nil {
			return p, nil
		}
		if pwd == "/" {
			break
		}
		pwd = filepath.Join(pwd, "../")
	}
	
	return "", fmt.Errorf("lookup path: %s faild", lookup)
}

func LookupPath(path, lookup string) (string, error) {
	var err error
	for {
		p := filepath.Join(path, lookup)
		_, err = os.Stat(p)
		if err == nil {
			return p, nil
		}
		if path == "/" {
			break
		}
		path = filepath.Join(path, "../")
	}
	
	return "", fmt.Errorf("lookup path: %s faild", lookup)
}

func Join(elem ...string) string {
	return filepath.Join(elem...)
}

// ParseCrossExtra 谨慎使用
func ParseCrossExtra(a, b string) string {
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
