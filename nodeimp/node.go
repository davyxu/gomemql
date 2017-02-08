package gomemql

import (
	"bytes"
	"fmt"
)

// 根据匹配方式的值集
type etcMatchKey struct {
	t   matchType
	key interface{}
}

// 索引节点
type indexNode struct {
	index int

	// 映射
	equalMapper map[interface{}]*indexNode // 等于集合

	etcMapper map[etcMatchKey][]*indexNode // 不等集合

	// 记录
	record []interface{}
}

func (self *indexNode) Name() string {
	return "IndexNode"
}

// 添加一个
func (self *indexNode) Add(r *record) {

	key := r.GetField(self.index)

	// 构建记录时, 最后一个字段直接添加记录
	if r.IsTerminate(self.index) {

		self.record = append(self.record, key)

	} else {

		n := self.EqualMatch(key)
		if n == nil {

			nextF := r.GetField(self.index + 1)
			if nextF == nil {
				return
			}

			n = r.NewNode(self.index + 1)

			self.AddIndex(matchType_Equal, key, n)
		}

		// 迭代添加每个字段对应的节点
		n.Add(r)

	}

}

func (self *indexNode) Print(b *bytes.Buffer) {

	writeLineWithIndent(self, b, fmt.Sprintf("[%s] index: %d", self.Name(), self.index))

	for k, v := range self.equalMapper {
		writeLineWithIndent(self, b, fmt.Sprintf("%v ==", k))

		v.Print(b)
	}

	for k, list := range self.etcMapper {

		for _, v := range list {
			writeLineWithIndent(self, b, fmt.Sprintf("%v %v", k.key, k.t))
			v.Print(b)
		}

	}

}

// 根据字段index, 找到对应的这个节点的
func (self *indexNode) IterateNodeByIndex(parent *indexNode, index int, callback func(*indexNode)) {

	if index == self.index {
		callback(parent)
		return
	}

	for _, v := range self.equalMapper {

		v.IterateNodeByIndex(v, index, callback)
	}

}

// 不等于添加索引
func (self *indexNode) AddIndex(t matchType, key interface{}, n *indexNode) {

	if n == nil {
		return
	}

	if t == matchType_Equal {

		self.equalMapper[key] = n

	} else {

		list, _ := self.etcMapper[etcMatchKey{t, key}]

		list = append(list, n)

		self.etcMapper[etcMatchKey{t, key}] = list

	}

}

// 等于匹配
func (self *indexNode) EqualMatch(key interface{}) *indexNode {

	if v, ok := self.equalMapper[key]; ok {

		return v
	}

	return nil
}

// 不等于匹配
func (self *indexNode) etcMatch(t matchType, key interface{}) []*indexNode {

	if v, ok := self.etcMapper[etcMatchKey{t, key}]; ok {

		return v
	}

	return nil
}

// 最终结果集调用
func (self *indexNode) invokeRecord(callback func(interface{}) bool) bool {

	if callback != nil {
		for _, v := range self.record {
			if !callback(v) {
				return false
			}
		}
	}

	return true
}

// 匹配值
func (self *indexNode) matchValue(q *Query, c *condition, callback func(interface{}) bool) bool {
	switch c.t {
	case matchType_Equal:
		if v := self.EqualMatch(c.value); v != nil {

			if !v.Execute(q, callback) {
				return false

			}
		}

	default:

		if list := self.etcMatch(c.t, c.value); list != nil {
			for _, v := range list {
				if !v.Execute(q, callback) {
					return false
				}
			}
		} else {
			fmt.Println("DO NOT USE INDEX")

			// 暴力匹配
			for k, v := range self.equalMapper {

				if compare(c.t, k, c.value) {

					if !v.Execute(q, callback) {
						return false
					}

				}

			}
		}

	}

	return true
}

func (self *indexNode) Execute(q *Query, callback func(interface{}) bool) bool {

	c := q.get(self.index)

	if c == nil {

		if !self.invokeRecord(callback) {
			return false
		}

	} else if !self.matchValue(q, c, callback) {

		return false
	}

	return true
}

func compare(t matchType, tabData, userExpect interface{}) bool {

	switch tabDataT := tabData.(type) {
	case int32:
		userExpectT := userExpect.(int32)

		switch t {
		case matchType_NotEqual:
			return tabDataT != userExpectT
		case matchType_Less:
			return tabDataT < userExpectT
		case matchType_LessEqual:
			return tabDataT <= userExpectT
		case matchType_Great:
			return tabDataT > userExpectT
		case matchType_GreatEqual:
			return tabDataT >= userExpectT
		}
	case int64:
		userExpectT := userExpect.(int64)

		switch t {
		case matchType_NotEqual:
			return tabDataT != userExpectT
		case matchType_Less:
			return tabDataT < userExpectT
		case matchType_LessEqual:
			return tabDataT <= userExpectT
		case matchType_Great:
			return tabDataT > userExpectT
		case matchType_GreatEqual:
			return tabDataT >= userExpectT
		}
	}

	return false
}

func newIndexNode(index int) *indexNode {
	return &indexNode{
		index:       index,
		equalMapper: make(map[interface{}]*indexNode),
		etcMapper:   make(map[etcMatchKey][]*indexNode),
	}
}
