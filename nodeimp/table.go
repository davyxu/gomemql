package gomemql

import (
	"bytes"
	"fmt"
)

// 存储静态的记录集合
type Table struct {
	root *indexNode

	nodeIdAcc int
}

func (self *Table) genID() int {
	self.nodeIdAcc++
	return self.nodeIdAcc
}

// 添加记录, 多个字段以逗号分割, 最后一个为整个记录集引用
func (self *Table) AddRecord(fields ...interface{}) {

	record := newRecord(self, fields)

	self.root.Add(record)
}

func (self *Table) Print() {

	var b bytes.Buffer

	self.root.Print(&b)

	fmt.Println(b.String())
}

// 在指定字段上构建不等索引, 范围[begin,end]
func (self *Table) GenIndexNotEqual(index int, begin, end int32) *Table {

	self.genIndex(index, matchType_NotEqual, begin, end)

	return self
}

// 在指定字段上构建小于索引, 范围[begin,end]
func (self *Table) GenIndexLess(index int, begin, end int32) *Table {

	self.genIndex(index, matchType_Less, begin, end)

	return self
}

// 在指定字段上构建小于等于索引, 范围[begin,end]
func (self *Table) GenIndexLessEqual(index int, begin, end int32) *Table {

	self.genIndex(index, matchType_LessEqual, begin, end)

	return self
}

// 在指定字段上构建大于索引, 范围[begin,end]
func (self *Table) GenIndexGreat(index int, begin, end int32) *Table {

	self.genIndex(index, matchType_Great, begin, end)

	return self
}

// 在指定字段上构建大于等于索引, 范围[begin,end]
func (self *Table) GenIndexGreatEqual(index int, begin, end int32) *Table {

	self.genIndex(index, matchType_GreatEqual, begin, end)

	return self
}

func (self *Table) genIndex(index int, t matchType, begin, end int32) {

	rootIndexNode := self.root

	rootIndexNode.IterateNodeByIndex(rootIndexNode, index, func(parentNode *indexNode) {

		var i, j int32

		// 遍历实际访问的数值
		for i = begin; i <= end; i++ {

			switch t {
			case matchType_NotEqual:

				for j = begin; j <= end && j != i; j++ {

					self.addIndexNode(t, parentNode, i, j)

				}
			case matchType_Less:

				for j = begin; j < i; j++ {

					self.addIndexNode(t, parentNode, i, j)

				}
			case matchType_LessEqual:

				for j = begin; j <= i; j++ {

					self.addIndexNode(t, parentNode, i, j)

				}
			case matchType_Great:

				for j = i + 1; j <= end; j++ {

					self.addIndexNode(t, parentNode, i, j)

				}
			case matchType_GreatEqual:

				for j = i; j <= end; j++ {

					self.addIndexNode(t, parentNode, i, j)
				}

			case matchType_Equal:
				panic("no need to create index of equal")
			}

		}

	})

}

func (self *Table) addIndexNode(t matchType, parentNode *indexNode, i, j int32) {
	nn := parentNode.EqualMatch(j)

	if nn == nil {
		nn = newIndexNode(self.genID(), parentNode.index)
		parentNode.AddIndex(matchType_Equal, j, nn)
	}

	parentNode.AddIndex(t, i, nn)
}

func NewTable() *Table {
	return &Table{
		root: newIndexNode(0, 0),
	}

}
