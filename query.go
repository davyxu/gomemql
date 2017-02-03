package gomemql

import "reflect"

type mergeData struct {
	count int // data在查询field中重复的次数
	data  interface{}
}

type condition struct {
	field *tableField
	mt    MatchType
	value interface{}
}

// 查询结果
type Query struct {
	mergeDataByPtr map[uintptr]*mergeData
	cons           []*condition
	result         map[uintptr]interface{}
	limit          int
	tab            *Table
	sortor         func(interface{}, interface{}) bool

	done bool
}

// 添加数据, 自动去重, 生成结果
func (self *Query) add(data interface{}) {
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
	if md.count >= len(self.cons) {
		self.result[ptr] = md.data
	}
}

// 条件约束
func (self *Query) Where(fieldName string, matchTypeStr string, value interface{}) *Query {

	matchType := getMatchTypeBySign(matchTypeStr)
	if matchType == MatchType_Unknown {
		panic("unknown match type: " + matchTypeStr)
	}

	con := &condition{
		field: self.tab.FieldByName(fieldName),
		mt:    matchType,
		value: value,
	}

	if con.field == nil {
		panic("field not found: " + fieldName)
	}

	self.cons = append(self.cons, con)

	return self
}

// 约束输出数量
func (self *Query) Limit(count int) *Query {

	if count < 0 {
		panic("count should > 0 ")
	}

	self.limit = count

	return self
}

// 根据排序函数排序
func (self *Query) SortBy(callback func(interface{}, interface{}) bool) *Query {

	self.sortor = callback

	return self
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
				con.field.Match(self, con.mt, con.value)
			}
		}
	}

	self.done = true
}

// 处理数据并且输出
func (self *Query) Result() []interface{} {

	self.do()

	// map转数组
	var ret *RecordList
	ret = newRecordListInitCount(len(self.result))
	var index int
	for _, v := range self.result {
		ret.set(index, v)
		index++
	}

	// 先排序
	if self.sortor != nil {
		ret.Sort(self.sortor)
	}

	if self.limit != -1 && self.limit < len(self.result) {
		// 约束数量
		ret.Resize(self.limit)
	}

	return ret.Raw()
}

// 直接访问原始数据, 不支持Limit,SortBy
func (self *Query) VisitRawResult(callback func(interface{}) bool) {

	self.do()

	if callback == nil {
		return
	}

	for _, v := range self.result {
		if !callback(v) {
			return
		}
	}

}

func NewQuery(tab *Table) *Query {

	return &Query{
		tab:            tab,
		limit:          -1,
		mergeDataByPtr: make(map[uintptr]*mergeData),
		result:         make(map[uintptr]interface{}),
	}
}
