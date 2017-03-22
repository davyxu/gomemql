package gomemql

import (
	"reflect"
)

type record struct {
	Data interface{} // [字段数*2]interface{}

	index int

	dataValue reflect.Value //  map做对比时, 这里置空
}

func (self *record) setField(value interface{}) {

	// 反射设置数组下标值
	self.dataValue.Index(self.index).Set(reflect.ValueOf(value))
	self.Data = self.dataValue.Interface()
	self.index++
}

func interfaceParamGetter(v interface{}) {
}

var interfaceType = reflect.TypeOf(interfaceParamGetter).In(0)

func newRecord(fieldCount int) *record {

	arrayType := reflect.ArrayOf(fieldCount*2, interfaceType)

	return &record{
		dataValue: reflect.New(arrayType).Elem(),
	}
}
