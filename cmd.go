package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/go-hayden-base/fs"
	"github.com/go-hayden-base/str"
)

func Exec(cmd string) ([]byte, error) {
	cmd = regexp.MustCompile(`\s{2,}`).ReplaceAllString(cmd, " ")
	args := strings.Split(cmd, " ")
	mainCommond := args[0]
	if len(args) > 1 {
		return exec.Command(mainCommond, args[1:]...).Output()
	}
	return exec.Command(mainCommond).Output()
}

func ExecThenOutput(cmd string, dest string) error {
	if dest == "" || !path.IsAbs(dest) {
		return errors.New("输出路径错误!")
	}
	if b, e := Exec(cmd); e != nil {
		return e
	} else {
		if e := fs.WriteFile(dest, b, true, os.ModePerm); e != nil {
			return e
		}
	}
	return nil
}

type enumerableCmd struct {
	Cmd        string
	err        error
	filterReg  string
	filterFunc func(line string, stop *bool) bool
}

func (s *enumerableCmd) Filter(reg string) str.IEnumerableString {
	s.filterReg = reg
	return s
}

func (s *enumerableCmd) FilterFunc(f func(line string, stop *bool) bool) str.IEnumerableString {
	s.filterFunc = f
	return s
}

func (s *enumerableCmd) Enumerate(f func(line string, err error)) {
	if f == nil {
		return
	}
	if b, e := Exec(s.Cmd); e != nil {
		f("", e)
		return
	} else {
		be := str.NewEnumerableBytes(b)
		be.Filter(s.filterReg).FilterFunc(s.filterFunc).Enumerate(f)
	}
}

func ExecForEnumerable(cmd string) str.IEnumerableString {
	ec := new(enumerableCmd)
	ec.Cmd = cmd
	return ec
}
