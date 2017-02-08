package gomemql_btree

import "github.com/google/btree"

type recordContext struct {
	item   btree.Item  // 用于btree判断
	record interface{} // 引用的记录集
}

type element struct {
	list   []recordContext
	header btree.Item
}

func (self *element) add(item btree.Item, record interface{}) {
	self.list = append(self.list, recordContext{
		item:   item,
		record: record,
	})
}

func newElement(item btree.Item, record interface{}) *element {

	self := &element{
		header: item,
	}
	self.add(item, record)

	return self
}

func (self *element) Less(than btree.Item) bool {

	other := than.(*element).header

	return self.header.Less(other)
}

// 数据源的一个字段, 对应实际的一个结构体
type Field struct {
	bt *btree.BTree
}

// 添加一条记录
func (self *Field) AddRecord(data btree.Item, record interface{}) {

	elem := newElement(data, record)

	pre := self.bt.Get(elem)

	if pre != nil {

		el := pre.(*element)
		el.add(data, record)

	} else {

		self.bt.ReplaceOrInsert(elem)

	}

}

func (self *Field) match(q *Query, t matchType, value btree.Item) {

	elem := newElement(value, nil)

	switch t {
	case matchType_Equal:

		result := self.bt.Get(elem)
		if result != nil {
			q.addResult(result.(*element))
		}
	case matchType_Less:

		self.bt.AscendLessThan(elem, func(i btree.Item) bool {

			q.addResult(i.(*element))
			return true
		})
	case matchType_Great:

		self.bt.DescendGreaterThan(elem, func(i btree.Item) bool {

			q.addResult(i.(*element))
			return true
		})
	case matchType_LessEqual:

		self.bt.DescendLessOrEqual(elem, func(i btree.Item) bool {

			q.addResult(i.(*element))
			return true
		})
	case matchType_GreatEqual:

		self.bt.AscendGreaterOrEqual(elem, func(i btree.Item) bool {

			q.addResult(i.(*element))
			return true
		})

	}

}

func NewField() *Field {

	return &Field{
		bt: btree.New(3),
	}

}
