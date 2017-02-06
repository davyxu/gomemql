package gomemql

import (
	"errors"
	"reflect"
)

type Table struct {
	fieldByName map[string]*tableField
	fields      []*tableField
}

// 添加一行数据
func (self *Table) AddRecord(record interface{}) {

	vRecord := reflect.Indirect(reflect.ValueOf(record))

	// 根据字段进行索引
	for i := 0; i < len(self.fields); i++ {

		recordField := self.fields[i]

		// 结构体中数值
		key := vRecord.Field(i).Interface()

		// 将数值添加到字段索引中, 同一个值可能有多个, 引用记录集合
		recordField.Add(key, record)

	}
}

func (self *Table) FieldByName(name string) *tableField {

	if v, ok := self.fieldByName[name]; ok {

		return v
	}

	return nil
}

func (self *Table) FieldByIndex(index int) *tableField {
	return self.fields[index]
}

func (self *Table) GenFieldIndexNotEqual(name string, begin, end int32) error {

	return self.genFieldIndex(name, matchType_NotEqual, begin, end)
}

func (self *Table) GenFieldIndexLess(name string, begin, end int32) error {

	return self.genFieldIndex(name, matchType_Less, begin, end)
}

func (self *Table) GenFieldIndexLessEqual(name string, begin, end int32) error {

	return self.genFieldIndex(name, matchType_LessEqual, begin, end)
}

func (self *Table) GenFieldIndexGreat(name string, begin, end int32) error {

	return self.genFieldIndex(name, matchType_Great, begin, end)
}

func (self *Table) GenFieldIndexGreatEqual(name string, begin, end int32) error {

	return self.genFieldIndex(name, matchType_GreatEqual, begin, end)
}

func (self *Table) genFieldIndex(name string, t matchType, begin, end int32) error {

	if self.NumFields() == 0 {
		return nil
	}

	if self.FieldByIndex(0).KeyCount() == 0 {
		return errors.New("require table data to gen index")
	}

	field := self.FieldByName(name)
	if field == nil {
		return errors.New("field not found: " + name)
	}

	var i, j int32

	// 遍历实际访问的数值
	for i = begin; i <= end; i++ {

		indexList := newRecordList()

		switch t {
		case matchType_NotEqual:

			for j = i; j <= end; j++ {
				if j == i {
					continue
				}

				list := field.getByKey(j)
				indexList.AddRange(list)
			}

		case matchType_Great:
			// 大于当前值的所有列表合并
			for j = i + 1; j <= end; j++ {
				list := field.getByKey(j)
				indexList.AddRange(list)
			}

		case matchType_GreatEqual:

			// 大于等于当前值的所有列表合并
			for j = i; j <= end; j++ {
				list := field.getByKey(j)
				indexList.AddRange(list)
			}

		case matchType_Less:

			for j = begin; j < i; j++ {
				list := field.getByKey(j)
				indexList.AddRange(list)
			}

		case matchType_LessEqual:

			for j = begin; j <= i; j++ {
				list := field.getByKey(j)
				indexList.AddRange(list)
			}
		case matchType_Equal:
			panic("no need to create index of equal")
		}

		field.addIndexData(t, i, indexList)

	}

	return nil
}

func (self *Table) NumFields() int {
	return len(self.fields)
}

func NewTable(userStruct interface{}) *Table {

	tStruct := reflect.TypeOf(userStruct)
	if tStruct.Kind() == reflect.Ptr {
		tStruct = tStruct.Elem()
	}

	self := &Table{
		fieldByName: make(map[string]*tableField),
		fields:      make([]*tableField, tStruct.NumField()),
	}

	for i := 0; i < tStruct.NumField(); i++ {

		fd := tStruct.Field(i)

		tf := newTableField()

		self.fields[i] = tf
		self.fieldByName[fd.Name] = tf
	}

	return self
}
