package gomemql

import "sort"

type RecordList struct {
	data []interface{}

	sortor func(interface{}, interface{}) bool
}

func (self *RecordList) Raw() []interface{} {
	return self.data
}

func (self *RecordList) Add(data interface{}) {
	self.data = append(self.data, data)
}

func (self *RecordList) AddRange(other *RecordList) {
	self.data = append(self.data, other.data...)
}

func (self *RecordList) Len() int {
	return len(self.data)
}

func (self *RecordList) Get(index int) interface{} {
	return self.data[index]
}

func (self *RecordList) Resize(count int) {
	self.data = self.data[0:count]
}

func (self *RecordList) set(index int, v interface{}) {
	self.data[index] = v
}

func (self *RecordList) Swap(i, j int) {
	self.data[i], self.data[j] = self.data[j], self.data[i]
}

func (self *RecordList) Less(i, j int) bool {

	return self.sortor(self.data[i], self.data[j])
}

func (self *RecordList) Sort(callback func(interface{}, interface{}) bool) {
	self.sortor = callback
	sort.Sort(self)
}

func newRecordList() *RecordList {
	return &RecordList{}
}

func newRecordListInitCount(count int) *RecordList {

	return &RecordList{
		data: make([]interface{}, count),
	}

}
