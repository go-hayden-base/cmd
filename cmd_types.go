package cmd

import "github.com/go-hayden-base/str"

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
