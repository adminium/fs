package fs

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Read(path string) (data []byte, err error) {
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return
}

func Remove(path string) (err error) {
	err = os.RemoveAll(path)
	return
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

func DeepFiles(path string, suffix ...string) (files []string, err error) {
	
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	
	sep := string(os.PathSeparator)
	
	for _, fi := range dir {
		if fi.IsDir() {
			var dirFiles []string
			dirFiles, err = Files(filepath.Join(path, sep, fi.Name()), suffix...)
			if err != nil {
				return
			}
			files = append(files, dirFiles...)
		} else {
			var r []string
			r, err = Files(filepath.Join(path, sep, fi.Name()), suffix...)
			if err != nil {
				return
			}
			files = append(files, r...)
		}
	}
	
	return
}

func Files(path string, suffix ...string) (files []string, err error) {
	
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	
	sep := string(os.PathSeparator)
	
	for _, fi := range dir {
		if !fi.IsDir() {
			if HasSuffix(suffix, fi.Name()) {
				files = append(files, filepath.Join(path, sep, fi.Name()))
			}
		}
	}
	
	return
}

func HasSuffix(suffix []string, fileName string) bool {
	for _, v := range suffix {
		if strings.HasSuffix(fileName, v) {
			return true
		}
	}
	return false
}

func MakeDir(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("mkdir -p %s", path))
		err = cmd.Run()
		if err != nil {
			return
		}
	}
	return os.MkdirAll(path, os.ModePerm)
}

func Write(path string, content []byte) (err error) {
	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return
	}
	return
}

func Append(path string, content []byte) (err error) {
	
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	
	defer func() {
		_ = f.Close()
	}()
	
	_, err = f.Write(content)
	
	return
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
	
	lp := pwd
	for {
		p := filepath.Join(lp, lookup)
		_, err = os.Stat(p)
		if err == nil {
			return p, nil
		}
		if lp == "/" {
			break
		}
		lp = filepath.Join(lp, "../")
	}
	
	return "", fmt.Errorf("lookup path: %s faild", lookup)
}

func LookupFrom(pwd, lookup string) (string, error) {
	var err error
	lp := pwd
	for {
		p := filepath.Join(lp, lookup)
		_, err = os.Stat(p)
		if err == nil {
			return p, nil
		}
		if lp == "/" {
			break
		}
		lp = filepath.Join(lp, "../")
	}
	
	return "", fmt.Errorf("lookup path: %s faild", lookup)
}

func Join(elem ...string) string {
	for i := len(elem) - 1; i > 0; i-- {
		if strings.HasPrefix(elem[i], string(os.PathSeparator)) {
			elem = elem[i:]
			break
		}
	}
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
