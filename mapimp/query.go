package gomemql

import "reflect"

type mergeData struct {
	count int // data在查询field中重复的次数
	data  interface{}
}

type condition struct {
	field *tableField
	t     matchType
	value interface{}
}

// 查询结果
type Query struct {
	mergeDataByPtr map[uintptr]*mergeData
	cons           []condition
	tab            *Table

	resultCallback func(interface{}) bool

	done bool
}

func (self *Query) addRecord(rl *RecordList) bool {

	for _, v := range rl.data {
		if !self.add(v) {
			return false
		}
	}

	return true
}

// 添加数据, 自动去重, 生成结果
func (self *Query) add(data interface{}) bool {
	ptr := reflect.ValueOf(data).Pointer()

	var md *mergeData
	if exists, ok := self.mergeDataByPtr[ptr]; ok {
		md = exists
	} else {
		md = &mergeData{
			data: data,
		}
		self.mergeDataByPtr[ptr] = md
	}

	md.count++

	// 求叉集
	if md.count >= len(self.cons) && self.resultCallback != nil {

		if !self.resultCallback(md.data) {
			return false
		}
	}

	return true
}

// 记录集合中等于value的所有值
func (self *Query) Equal(fieldName string, value interface{}) *Query {

	self.newCondition(fieldName, matchType_Equal, value)

	return self
}

// 记录集合中不等于value的所有值
func (self *Query) NotEqual(fieldName string, value interface{}) *Query {

	self.newCondition(fieldName, matchType_NotEqual, value)

	return self
}

// 记录集合中小于value的所有值
func (self *Query) Less(fieldName string, value interface{}) *Query {

	self.newCondition(fieldName, matchType_Less, value)

	return self
}

// 记录集合中小于等于value的所有值
func (self *Query) LessEqual(fieldName string, value interface{}) *Query {

	self.newCondition(fieldName, matchType_LessEqual, value)

	return self
}

// 记录集合中大于value的所有值
func (self *Query) Great(fieldName string, value interface{}) *Query {

	self.newCondition(fieldName, matchType_Great, value)

	return self
}

// 记录集合中大于等于value的所有值
func (self *Query) GreatEqual(fieldName string, value interface{}) *Query {

	self.newCondition(fieldName, matchType_GreatEqual, value)

	return self
}

func (self *Query) newCondition(fieldName string, t matchType, value interface{}) {

	con := condition{
		field: self.tab.FieldByName(fieldName),
		t:     t,
		value: value,
	}

	if con.field == nil {
		panic("field not found: " + fieldName)
	}

	self.cons = append(self.cons, con)
}

func (self *Query) do() {

	if self.done {
		return
	}

	// 结构体没有字段
	if self.tab.NumFields() > 0 {

		// 没有任何条件约束
		if len(self.cons) == 0 {

			// 返回所有
			self.tab.FieldByIndex(0).All(self)

		} else {

			// 根据条件匹配
			for _, con := range self.cons {
				con.field.Match(self, con.t, con.value)
			}
		}
	}

	self.done = true
}

// 处理数据并且输出
func (self *Query) Result(callback func(interface{}) bool) {

	self.resultCallback = callback

	self.do()
}

func NewQuery(tab *Table) *Query {

	return &Query{
		tab:            tab,
		mergeDataByPtr: make(map[uintptr]*mergeData),
	}
}
