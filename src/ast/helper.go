package ast

import (
	"bytes"

	ktype "github.com/KhushPatibandha/Kolon/src/kType"
)

// ------------------------------------------------------------------------------------------------------------------
// Function Parameters
// ------------------------------------------------------------------------------------------------------------------
type FunctionParameter struct {
	ParameterName *Identifier
	ParameterType *ktype.Type
}

func (fp *FunctionParameter) TokenValue() string { return fp.ParameterName.Token.Value }
func (fp *FunctionParameter) String() string {
	var out bytes.Buffer
	out.WriteString(fp.ParameterName.String())
	out.WriteString(": ")
	out.WriteString(fp.ParameterType.String())
	return out.String()
}

func (fp *FunctionParameter) Equals(other *FunctionParameter) bool {
	if other == nil {
		return false
	}
	return fp.ParameterName.Equals(other.ParameterName) && fp.ParameterType.Equals(other.ParameterType)
}
