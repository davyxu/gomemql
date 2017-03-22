package gomemql

type Query struct {
	child []*condition
	tab   *Table
}

// 构建1个查询字段, 期望 表中值==value
func (self *Query) Equal(value interface{}) *Query {

	self.newCondition(matchType_Equal, value)

	return self
}

// 构建1个查询字段, 期望 表中值!=value
func (self *Query) NotEqual(value interface{}) *Query {

	self.newCondition(matchType_NotEqual, value)

	return self
}

// 构建1个查询字段, 期望 表中值<value
func (self *Query) Less(value interface{}) *Query {

	self.newCondition(matchType_Less, value)

	return self
}

// 构建1个查询字段, 期望 表中值<=value
func (self *Query) LessEqual(value interface{}) *Query {

	self.newCondition(matchType_LessEqual, value)

	return self
}

// 构建1个查询字段, 期望 表中值>value
func (self *Query) Great(value interface{}) *Query {

	self.newCondition(matchType_Great, value)

	return self
}

// 构建1个查询字段, 期望 表中值>=value
func (self *Query) GreatEqual(value interface{}) *Query {

	self.newCondition(matchType_GreatEqual, value)

	return self
}

// 查询返回的结果
func (self *Query) Result(callback func(interface{}) bool) {

	self.tab.root.Execute(self, callback)
}

func (self *Query) newCondition(t matchType, value interface{}) {

	f := &condition{
		index: len(self.child),
		value: value,
		t:     t,
	}

	self.child = append(self.child, f)
}

func (self *Query) get(index int) *condition {

	if index >= len(self.child) {
		return nil
	}

	return self.child[index]
}

func (self *Query) isTerminate(index int) bool {

	return index >= len(self.child)
}

// 开始一个新的查询
func NewQuery(tab *Table) *Query {
	return &Query{
		tab: tab,
	}
}
