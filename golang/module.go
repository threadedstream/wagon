package golang

import "github.com/traefik/yaegi/interp"

type GoModule struct {
	*interp.Interpreter
	Prog *interp.Program
	Src  string
}

func NewGoModule() *GoModule {
	return &GoModule{}
}

func (*GoModule) IsModule() bool {
	return true
}
