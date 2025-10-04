package ast

import (
	"bytes"
)

// ------------------------------------------------------------------------------------------------------------------
// Function Parameters
// ------------------------------------------------------------------------------------------------------------------
type FunctionParameter struct {
	ParameterName *Identifier
	ParameterType *Type
}

func (fp *FunctionParameter) TokenValue() string { return fp.ParameterName.Token.Value }
func (fp *FunctionParameter) String() string {
	var out bytes.Buffer
	out.WriteString(fp.ParameterName.String())
	out.WriteString(": ")
	out.WriteString(fp.ParameterType.String())
	return out.String()
}
