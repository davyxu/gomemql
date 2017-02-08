package gomemql_btree

import (
	"reflect"

	"github.com/google/btree"
)

type condition struct {
	field *Field
	t     matchType
	value btree.Item
}

type mergeData struct {
	count  int // data在查询field中重复的次数
	record interface{}
}

// 查询对象
type Query struct {
	cons           []condition
	mergeDataByPtr map[uintptr]*mergeData
	resultCallback func(interface{})

	done bool
}

func (self *Query) addResult(el *element) {

	for _, v := range el.list {

		ptr := reflect.ValueOf(v.item).Pointer()

		var md *mergeData
		if exists, ok := self.mergeDataByPtr[ptr]; ok {
			md = exists
		} else {
			md = &mergeData{
				record: v.record,
			}
			self.mergeDataByPtr[ptr] = md
		}

		md.count++

		// 求叉集
		if md.count >= len(self.cons) && self.resultCallback != nil {

			self.resultCallback(v.record)

		}
	}

}

// 记录集合中等于value的所有值
func (self *Query) Equal(f *Field, value btree.Item) *Query {

	self.newCondition(f, matchType_Equal, value)

	return self
}

// 记录集合中小于value的所有值
func (self *Query) Less(f *Field, value btree.Item) *Query {

	self.newCondition(f, matchType_Less, value)

	return self
}

// 记录集合中小于等于value的所有值
func (self *Query) LessEqual(f *Field, value btree.Item) *Query {

	self.newCondition(f, matchType_LessEqual, value)

	return self
}

// 记录集合中大于value的所有值
func (self *Query) Great(f *Field, value btree.Item) *Query {

	self.newCondition(f, matchType_Great, value)

	return self
}

// 记录集合中大于等于value的所有值
func (self *Query) GreatEqual(f *Field, value btree.Item) *Query {

	self.newCondition(f, matchType_GreatEqual, value)

	return self
}

func (self *Query) newCondition(f *Field, t matchType, value btree.Item) {

	c := condition{
		field: f,
		t:     t,
		value: value,
	}

	self.cons = append(self.cons, c)

}

// 开始匹配
func (self *Query) Start() {
	if self.done {
		return
	}

	// 根据条件匹配
	for _, con := range self.cons {
		con.field.match(self, con.t, con.value)
	}

	self.done = true
}

// 创建一个查询, 给定一个结果返回回调
func NewQuery(callback func(interface{})) *Query {
	return &Query{
		resultCallback: callback,
		mergeDataByPtr: make(map[uintptr]*mergeData),
	}
}
