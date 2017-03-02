package cmd

import fdt "github.com/go-hayden-base/foundation"

type tEnumerableCmd struct {
	Cmd        string
	FilterFunc func(item interface{}) bool
}

func (s *tEnumerableCmd) Filter(f func(item interface{}) bool) fdt.IEnumerable {
	s.FilterFunc = f
	return s
}

func (s *tEnumerableCmd) Enumerate(f func(itme interface{}, err error, stop *bool)) {
	if f == nil {
		return
	}
	stop := false
	b, err := Exec(s.Cmd)
	if err != nil {
		f("", err, &stop)
		return
	}
	eByte := fdt.NewEnumerableBytes(b)
	eByte.Filter(s.FilterFunc).Enumerate(f)
}

func NewEnumerableCmd(cmd string) fdt.IEnumerable {
	aEnumerableCmd := tEnumerableCmd{Cmd: cmd}
	return &aEnumerableCmd
}
