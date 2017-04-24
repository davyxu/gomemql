package gomemql

import (
	"bytes"
	"fmt"
	"reflect"
)

type recordValue struct {
	tagList []interface{}
}

func (self *recordValue) TagString() string {
	var buff bytes.Buffer

	for _, tag := range self.tagList {
		buff.WriteString(fmt.Sprintln(tag))
		buff.WriteString(" ")
	}

	return buff.String()
}

// 存储静态的记录集合
type Table struct {
	recordByMultiKey map[interface{}]*recordValue

	fieldKind []reflect.Kind
}

func (self *Table) FieldCount() int {
	return len(self.fieldKind)
}

func (self *Table) matchFieldKind(index int, kind reflect.Kind) {

	if index >= len(self.fieldKind) {
		panic("Out of bound of field " + fmt.Sprintln(index, len(self.fieldKind)))
	}

	if self.fieldKind[index] != kind {
		panic("Record field type not match")
	}
}

// 添加记录, 多个字段以逗号分割, 最后一个为整个记录集引用
func (self *Table) AddRecord(tag interface{}, params ...interface{}) {

	record := newRecordByValues(params)

	needBuildFieldKind := self.fieldKind == nil

	for index, v := range params {

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

	rv, ok := self.recordByMultiKey[record.raw]
	if !ok {
		rv = &recordValue{}
		self.recordByMultiKey[record.raw] = rv
	}

	rv.tagList = append(rv.tagList, tag)

}

func (self *Table) GenIndexGreat(fieldIndex int, begin, end int32) {
	self.genIndex(fieldIndex, matchType_Great, begin, end)
}

func (self *Table) GenIndexGreatEqual(fieldIndex int, begin, end int32) {
	self.genIndex(fieldIndex, matchType_GreatEqual, begin, end)
}

func (self *Table) GenIndexNotEqual(fieldIndex int, begin, end int32) {
	self.genIndex(fieldIndex, matchType_NotEqual, begin, end)
}

func (self *Table) GenIndexLess(fieldIndex int, begin, end int32) {
	self.genIndex(fieldIndex, matchType_Less, begin, end)
}

func (self *Table) GenIndexLessEqual(fieldIndex int, begin, end int32) {
	self.genIndex(fieldIndex, matchType_LessEqual, begin, end)
}

func (self *Table) genIndex(fieldIndex int, t matchType, begin, end int32) {

	createdMap := map[interface{}]*recordValue{}

	for key, rv := range self.recordByMultiKey {

		// 目标值
		targetRecord := newRecordFromRaw(key)
		_, fieldValue := targetRecord.Value(fieldIndex)

		// 遍历实际访问的数值
		for fromkey := begin; fromkey <= end; fromkey++ {

			var insertIndexRecord bool
			switch t {
			case matchType_Great:
				insertIndexRecord = fieldValue.(int32) > fromkey
			case matchType_GreatEqual:
				insertIndexRecord = fieldValue.(int32) >= fromkey
			case matchType_LessEqual:
				insertIndexRecord = fieldValue.(int32) <= fromkey
			case matchType_Less:
				insertIndexRecord = fieldValue.(int32) < fromkey
			case matchType_NotEqual:
				insertIndexRecord = fieldValue.(int32) != fromkey
			case matchType_Equal:

			default:

			}

			if insertIndexRecord {
				targetRecord.SetValue(fieldIndex, t, fromkey)

				exists, ok := createdMap[targetRecord.raw]

				if !ok {
					exists = &recordValue{}
					createdMap[targetRecord.raw] = exists
				}

				// 符合条件的结果集合
				exists.tagList = append(exists.tagList, rv.tagList...)
			}

		}

	}

	// 将创建的结合合并到最终索引中
	for k, v := range createdMap {
		self.recordByMultiKey[k] = v
	}

}

func (self *Table) String() string {

	var buff bytes.Buffer

	for raw, value := range self.recordByMultiKey {
		buff.WriteString(fmt.Sprintf("%v -> %v", raw, value.TagString()))

	}

	return buff.String()
}
func (self *Table) match(record interface{}) []interface{} {

	if v, ok := self.recordByMultiKey[record]; ok {
		return v.tagList
	}

	return nil
}

func NewTable() *Table {
	return &Table{
		recordByMultiKey: make(map[interface{}]*recordValue),
	}

}
