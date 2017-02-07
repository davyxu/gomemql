package gomemql

// 用户输入的字段列表转换为记录
type record struct {
	fields []interface{}
}

func (self *record) GetField(index int) interface{} {
	if index >= len(self.fields) {
		return nil
	}

	return self.fields[index]
}

func (self *record) NewNode(index int) node {

	if index >= len(self.fields)-1 {
		return newRecordNode(index)
	}

	return newIndexNode(index)
}

func newRecord(fields []interface{}) *record {
	return &record{
		fields: fields,
	}

}
