package gomemql

import (
	"reflect"
)

type record struct {
	raw interface{} // [字段数*2]interface{}

	index int

	rawValue reflect.Value

	fieldCount int
}

func (self *record) FieldCount() int {
	return self.fieldCount
}

func (self *record) SetValue(fieldIndex int, t matchType, value interface{}) {

	self.rawSetValue(fieldIndex*2, int(t))

	self.rawSetValue(fieldIndex*2+1, value)
}

func (self *record) Value(fieldIndex int) (matchType, interface{}) {

	return matchType(self.rawValue.Index(fieldIndex * 2).Interface().(int)),
		self.rawValue.Index(fieldIndex*2 + 1).Interface()
}

func (self *record) rawSetValue(index int, value interface{}) {

	// 反射设置数组下标值
	self.rawValue.Index(index).Set(reflect.ValueOf(value))
	self.raw = self.rawValue.Interface()
}

func (self *record) initRaw(count int) {
	arrayType := reflect.ArrayOf(count*2, reflect.TypeOf(func(interface{}) {}).In(0))

	self.rawValue = reflect.New(arrayType).Elem()

	self.fieldCount = count
}

func newRecord(fieldCount int) *record {
	self := &record{}

	self.initRaw(fieldCount)

	return self
}

func newRecordFromRaw(raw interface{}) *record {

	// 临时构造一个结构, 只取值
	dummy := &record{
		rawValue: reflect.ValueOf(raw),
	}
	dummy.fieldCount = dummy.rawValue.Len() / 2

	self := &record{}
	self.initRaw(dummy.fieldCount)

	// 因为直接通过ValueOf( raw) 是无法Address的, 所以拷贝一份
	for i := 0; i < dummy.fieldCount; i++ {
		t, v := dummy.Value(i)
		self.SetValue(i, t, v)
	}

	return self
}

func newRecordByValues(values []interface{}) *record {

	self := &record{}

	self.initRaw(len(values))

	for index, v := range values {
		self.SetValue(index, matchType_Equal, v)
	}

	return self
}
