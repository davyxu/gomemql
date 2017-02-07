package gomemql

import (
	"bytes"
	"fmt"
)

// 记录节点, 保存父级符合条件的记录集合
type recordNode struct {
	baseNode

	record []interface{}
}

func (self *recordNode) Name() string {
	return "RecordNode"
}

func (self *recordNode) Print(b *bytes.Buffer) {

	self.WriteLineWithIndent(b, fmt.Sprintf("[%s] index: %d", self.Name(), self.index))

	for _, v := range self.record {
		self.WriteLineWithIndent(b, fmt.Sprintf("%v", v))

	}

}
func (self *recordNode) Execute(fl *Query, callback func(interface{}) bool) bool {

	if callback == nil {
		return true
	}

	for _, v := range self.record {
		if !callback(v) {
			return false
		}
	}

	return true

}

func (self *recordNode) Add(r *record) {

	value := r.GetField(self.index)

	self.record = append(self.record, value)

}

func newRecordNode(index int) node {

	return &recordNode{
		baseNode: baseNode{
			index: index,
		},
	}

}
