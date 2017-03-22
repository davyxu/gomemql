package gomemql

import (
	"reflect"
)

// 存储静态的记录集合
type Table struct {
	recordByMultiKey map[interface{}][]interface{}

	fieldKind []reflect.Kind
}

func (self *Table) FieldCount() int {
	return len(self.fieldKind)
}

func (self *Table) matchFieldKind(index int, kind reflect.Kind) {
	if self.fieldKind[index] != kind {
		panic("Record field type not match")
	}
}

// 添加记录, 多个字段以逗号分割, 最后一个为整个记录集引用
func (self *Table) AddRecord(tag interface{}, params ...interface{}) {

	record := newRecord(len(params))

	needBuildFieldKind := self.fieldKind == nil

	for index, v := range params {
		record.setField(int(matchType_Equal))
		record.setField(v)

		valueKind := reflect.TypeOf(v).Kind()

		if needBuildFieldKind {
			self.fieldKind = append(self.fieldKind, valueKind)
		} else {
			self.matchFieldKind(index, valueKind)
		}
	}

	if !needBuildFieldKind && len(self.fieldKind) != len(params) {
		panic("Record param count should equal")
	}

	list, _ := self.recordByMultiKey[record.Data]

	list = append(list, tag)

	self.recordByMultiKey[record.Data] = list

}

func (self *Table) match(record interface{}) []interface{} {

	if v, ok := self.recordByMultiKey[record]; ok {
		return v
	}

	return nil
}

func NewTable() *Table {
	return &Table{
		recordByMultiKey: make(map[interface{}][]interface{}),
	}

}
