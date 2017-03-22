package gomemql

// 用户输入的字段列表转换为记录
type record struct {
	fields []interface{}

	tab *Table
}

func (self *record) GetField(index int) interface{} {
	if index >= len(self.fields) {
		return nil
	}

	return self.fields[index]
}

func (self *record) IsTerminate(index int) bool {
	return index >= len(self.fields)-1
}

func (self *record) NewNode(index int) *indexNode {

	return newIndexNode(self.tab.genID(), index)
}

func newRecord(tab *Table, fields []interface{}) *record {
	return &record{
		tab:    tab,
		fields: fields,
	}

}
