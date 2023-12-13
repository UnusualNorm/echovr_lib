package messages

import (
	echovr "github.com/unusualnorm/echovr_lib"
)

var TemplateSymbol uint64 = echovr.GenerateSymbol("Template")

type Template struct {
}

func (m *Template) Symbol() uint64 {
	return TemplateSymbol
}

func (m *Template) Stream(s *echovr.EasyStream) error {
	return echovr.RunErrorFunctions([]func() error{})
}
