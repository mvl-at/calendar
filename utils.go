package calendar

import (
	"bytes"
	"fmt"
)

type strBuffer struct {
	bytes.Buffer
}

func (s *strBuffer) WriteFmt(format string, a ...interface{}) {
	s.WriteString(fmt.Sprintf(format+"\n", a...))
}
