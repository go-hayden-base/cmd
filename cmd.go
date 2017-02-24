package cmd

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"io"

	"github.com/go-hayden-base/fs"
	"github.com/go-hayden-base/str"
)

func Cmd(cmd string) *exec.Cmd {
	cmd = regexp.MustCompile(`\s{2,}`).ReplaceAllString(cmd, " ")
	args := strings.Split(cmd, " ")
	mainCommond := args[0]
	if len(args) > 1 {
		return exec.Command(mainCommond, args[1:]...)
	}
	return exec.Command(mainCommond)
}

func Exec(cmd string) ([]byte, error) {
	return Cmd(cmd).Output()
}

func ExecOutputPerLine(cmd string, f func(line string)) error {
	aCmd := Cmd(cmd)
	if f == nil {
		_, e := aCmd.Output()
		return e
	}

	stdout, e := aCmd.StdoutPipe()
	if e != nil {
		return e
	}

	aCmd.Start()
	aReader := bufio.NewReader(stdout)
	for {
		line, e := aReader.ReadString('\n')
		if e != nil {
			if e != io.EOF {
				f(e.Error())
			}
			break
		}
		f(line)
	}
	aCmd.Wait()
	return nil
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

func ExecForEnumerable(cmd string) str.IEnumerableString {
	ec := new(enumerableCmd)
	ec.Cmd = cmd
	return ec
}
