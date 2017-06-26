package gomemql

import (
	"errors"
	"reflect"
	"sort"
)

type Query struct {
	qr  *record
	tab *Table

	conditions int

	result  []interface{}
	matched bool
}

func (self *Query) ConditionCount() int {
	return self.conditions
}

// 构建1个查询字段, 期望 表中值==value
func (self *Query) Equal(value interface{}) *Query {

	self.newCondition(matchType_Equal, value)

	return self
}

func (self *Query) Great(value interface{}) *Query {

	self.newCondition(matchType_Great, value)

	return self
}

func (self *Query) GreatEqual(value interface{}) *Query {

	self.newCondition(matchType_GreatEqual, value)

	return self
}

func (self *Query) LessEqual(value interface{}) *Query {

	self.newCondition(matchType_LessEqual, value)

	return self
}

func (self *Query) Less(value interface{}) *Query {

	self.newCondition(matchType_Less, value)

	return self
}

func (self *Query) NotEqual(value interface{}) *Query {

	self.newCondition(matchType_NotEqual, value)

	return self
}

func (self *Query) newCondition(t matchType, value interface{}) {

	self.tab.matchFieldKind(self.conditions, reflect.TypeOf(value).Kind())

	self.qr.SetValue(self.conditions, t, value)

	self.conditions++
}

var ErrRequireValueOfSlice = errors.New("require value of slice")

// 给定一个切片地址, 转换为切片的类型
func (self *Query) ConverTo(targetPtr interface{}) error {

	rt := reflect.TypeOf(targetPtr).Elem()

	if rt.Kind() != reflect.Slice {
		return ErrRequireValueOfSlice
	}

	elementType := rt.Elem()

	ret := reflect.New(reflect.SliceOf(elementType)).Elem()

	for _, v := range self.result {

		targetElement := reflect.ValueOf(v).Convert(elementType)

		ret = reflect.Append(ret, targetElement)
	}

	reflect.ValueOf(targetPtr).Elem().Set(ret)

	return nil
}

// 根据tag里的结构体字段名对应的value去重, 并返回到result中
func (self *Query) DistinctByField(tagFieldName string, sortCallback func(a, b interface{}) bool) *Query {

	var sortFunc func(a, b *Group) bool

	if sortCallback != nil {
		sortFunc = func(a, b *Group) bool {
			return sortCallback(a.Key, b.Key)
		}
	}

	result := self.GroupByField(tagFieldName, sortFunc)

	self.result = make([]interface{}, len(result))

	for index, v := range result {
		self.result[index] = v.Key
	}

	return self
}

type Group struct {
	Key  interface{}
	List []interface{}
}

func (self *Query) GroupByField(tagFieldName string, sortCallback func(a, b *Group) bool) []*Group {

	self.match()

	groupByField := make(map[interface{}]*Group)

	for _, r := range self.result {

		rv := reflect.Indirect(reflect.ValueOf(r))

		fieldValue := rv.FieldByName(tagFieldName)

		if !fieldValue.IsValid() {
			continue
		}

		group, _ := groupByField[fieldValue.Interface()]

		if group == nil {
			group = &Group{Key: fieldValue.Interface()}
			groupByField[fieldValue.Interface()] = group
		}

		group.List = append(group.List, r)

	}

	var grouplist []*Group

	for _, group := range groupByField {
		grouplist = append(grouplist, group)
	}

	if sortCallback != nil {
		sort.Slice(grouplist, func(i, j int) bool {

			return sortCallback(grouplist[i], grouplist[j])
		})
	}

	return grouplist
}

func (self *Query) match() {
	if !self.matched && self.conditions > 0 {
		self.result = self.tab.match(self.qr.raw)
		self.matched = true
	}
}

// 查询返回的结果
func (self *Query) Result() []interface{} {

	self.match()

	return self.result
}

func (self *Query) All() *Query {

	for _, v := range self.tab.recordByMultiKey {
		self.result = append(self.result, v.tagList...)
	}

	return self
}

// 开始一个新的查询
func NewQuery(tab *Table) *Query {
	return &Query{
		tab: tab,
		qr:  newRecord(tab.FieldCount()),
	}
}
