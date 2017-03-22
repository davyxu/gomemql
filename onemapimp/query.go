package gomemql

import (
	"reflect"
)

type Query struct {
	qr  *record
	tab *Table

	conditions int
}

// 构建1个查询字段, 期望 表中值==value
func (self *Query) Equal(value interface{}) *Query {

	self.newCondition(matchType_Equal, value)

	return self
}

func (self *Query) newCondition(t matchType, value interface{}) {

	self.tab.matchFieldKind(self.conditions, reflect.TypeOf(value).Kind())

	self.qr.setField(int(t))
	self.qr.setField(value)

	self.conditions++

}

// 查询返回的结果
func (self *Query) Result() []interface{} {

	return self.tab.match(self.qr.Data)
}

// 开始一个新的查询
func NewQuery(tab *Table) *Query {
	return &Query{
		tab: tab,
		qr:  newRecord(tab.FieldCount()),
	}
}
