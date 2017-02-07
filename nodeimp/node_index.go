package gomemql

import (
	"bytes"
	"fmt"
)

// 根据匹配方式的值集
type matchKey struct {
	t   matchType
	key interface{}
}

// 索引节点
type indexNode struct {
	baseNode

	equalMapper map[interface{}]node // 等于集合

	etcMapper map[matchKey][]node // 不等集合
}

func (self *indexNode) Name() string {
	return "IndexNode"
}

// 添加一个
func (self *indexNode) Add(r *record) {

	key := r.GetField(self.index)

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

func (self *indexNode) Print(b *bytes.Buffer) {

	self.WriteLineWithIndent(b, fmt.Sprintf("[%s] index: %d", self.Name(), self.index))

	for k, v := range self.equalMapper {
		self.WriteLineWithIndent(b, fmt.Sprintf("%v ==", k))

		v.Print(b)
	}

	for k, list := range self.etcMapper {

		for _, v := range list {
			self.WriteLineWithIndent(b, fmt.Sprintf("%v ==", k.key, k.t))
			v.Print(b)
		}

	}

}

func (self *indexNode) IterateNodeByIndex(parent *indexNode, index int, callback func(*indexNode)) {

	if index == self.index {
		callback(parent)
		return
	}

	for _, v := range self.equalMapper {

		idx, ok := v.(*indexNode)
		if !ok {

			return
		}

		idx.IterateNodeByIndex(idx, index, callback)

	}

}

func (self *indexNode) AddIndex(t matchType, key interface{}, n node) {

	if n == nil {
		return
	}

	if t == matchType_Equal {

		self.equalMapper[key] = n

	} else {

		list, _ := self.etcMapper[matchKey{t, key}]

		list = append(list, n)

		self.etcMapper[matchKey{t, key}] = list

	}

}

func (self *indexNode) EqualMatch(key interface{}) node {

	if v, ok := self.equalMapper[key]; ok {

		return v
	}

	return nil
}

func (self *indexNode) match(t matchType, key interface{}) []node {

	if v, ok := self.etcMapper[matchKey{t, key}]; ok {

		return v
	}

	return nil
}

func (self *indexNode) Execute(fl *Query, callback func(interface{}) bool) bool {

	f := fl.get(self.index)

	switch f.t {
	case matchType_Equal:
		if v := self.EqualMatch(f.value); v != nil {

			if !v.Execute(fl, callback) {
				return false

			}
		}

	default:

		if list := self.match(f.t, f.value); list != nil {
			for _, v := range list {
				if !v.Execute(fl, callback) {
					return false
				}
			}
		} else {
			fmt.Println("DO NOT USE INDEX")

			// 暴力匹配
			for k, v := range self.equalMapper {

				if compare(f.t, k, f.value) {

					if !v.Execute(fl, callback) {
						return false
					}

				}

			}
		}

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

func newIndexNode(index int) node {
	return &indexNode{
		baseNode: baseNode{
			index: index,
		},
		equalMapper: make(map[interface{}]node),
		etcMapper:   make(map[matchKey][]node),
	}
}
