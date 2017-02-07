package gomemql

import "bytes"

type node interface {
	Execute(f *Query, callback func(interface{}) bool) bool

	Add(r *record)

	Print(b *bytes.Buffer)

	Name() string
}

type baseNode struct {
	index int
}

func (self *baseNode) WriteLineWithIndent(b *bytes.Buffer, str string) {

	for i := 0; i < self.index*4; i++ {
		b.WriteString(" ")
	}

	b.WriteString(str)

	b.WriteString("\n")

}
