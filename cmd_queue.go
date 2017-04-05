package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"regexp"
	"strings"

	"github.com/go-hayden-base/foundation"
)

type CmdExecer interface {
	Cmd(c string) CmdExecer
	Cd(d string) CmdExecer
	Dir(d string) CmdExecer
	Map(mapFunc func(line string) string) CmdExecer
	Reduce(reduceFunc func(lines []string) []string) CmdExecer
	Fall() CmdExecer
	Output() ([]string, error)
}

type tCmdQueue struct {
	cmds        []*tCmdQueueTask
	currentTask *tCmdQueueTask
	dir         string
}

type tCmdQueueTask struct {
	cmd        string
	dir        string
	fall       bool
	mapFunc    func(line string) string
	reduceFunc func(lines []string) []string
}

func NewCmdExecer(dir string) CmdExecer {
	aQueue := new(tCmdQueue)
	aQueue.dir = dir
	return aQueue
}

func (s *tCmdQueue) Cmd(c string) CmdExecer {
	if s.cmds == nil {
		s.cmds = make([]*tCmdQueueTask, 0, 5)
	}
	aTask := new(tCmdQueueTask)
	aTask.cmd = c
	aTask.dir = s.dir
	s.cmds = append(s.cmds, aTask)
	s.currentTask = aTask
	return s
}

func (s *tCmdQueue) Cd(d string) CmdExecer {
	if s.currentTask != nil {
		s.currentTask.dir = filepath.Join(s.currentTask.dir, d)
	}
	return s
}

func (s *tCmdQueue) Dir(d string) CmdExecer {
	if s.currentTask != nil {
		s.currentTask.dir = d
	}
	return s
}

func (s *tCmdQueue) Map(mapFunc func(line string) string) CmdExecer {
	if s.currentTask != nil {
		s.currentTask.mapFunc = mapFunc
	}
	return s
}

func (s *tCmdQueue) Reduce(reduceFunc func(lines []string) []string) CmdExecer {
	if s.currentTask != nil {
		s.currentTask.reduceFunc = reduceFunc
	}
	return s
}

func (s *tCmdQueue) Fall() CmdExecer {
	if s.currentTask != nil {
		s.currentTask.fall = true
	}
	return s
}

func (s *tCmdQueue) Output() ([]string, error) {
	l := len(s.cmds)
	if l == 0 {
		return nil, errors.New("No commond to execute!")
	}
	var output []string
	hasFall := false
	for idx, aCmdTask := range s.cmds {
		cmdStr := aCmdTask.cmd
		if hasFall && output != nil {
			args := make([]interface{}, len(output), len(output))
			for idx, item := range output {
				args[idx] = item
			}
			cmdStr = fmt.Sprintf(cmdStr, args...)
		}
		hasFall = aCmdTask.fall
		aCmd := Cmd(cmdStr)
		if aCmdTask.dir != "" {
			aCmd.Dir = aCmdTask.dir
		}
		b, err := aCmd.Output()
		if err != nil {
			return nil, funcFmtError(cmdStr, err)
		}
		if idx < l-1 && !aCmdTask.fall {
			output = nil
			continue
		}
		output = make([]string, 0, 5)
		foundation.NewEnumerableBytes(b).Enumerate(func(item interface{}, err error, stop *bool) {
			if err != nil {
				return
			}
			str, ok := item.(string)
			if !ok {
				return
			}
			if aCmdTask.mapFunc != nil {
				str = aCmdTask.mapFunc(str)
			}
			if str != "" {
				output = append(output, str)
			}
		})
		if aCmdTask.reduceFunc != nil {
			output = aCmdTask.reduceFunc(output)
		}
	}
	return output, nil
}

func funcFmtError(cmdStr string, err error) error {
	return errors.New("Exec '" + cmdStr + "' failed: " + err.Error())
}

func IsCmdError(cmdStr string, err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	idx := strings.Index(s, "' failed: ")
	if idx < 0 {
		return false
	}
	s = s[:idx]
	if !strings.HasPrefix(s, "Exec '") {
		return false
	}
	s = s[6:]
	reg := regexp.MustCompile(`\s{2,}`)
	s = reg.ReplaceAllString(s, " ")
	cmdStr = reg.ReplaceAllString(cmdStr, " ")
	return s == cmdStr || strings.HasPrefix(s, cmdStr)
}
