package gomemql

type unequalData struct {
	matchTypeList [matchType_MAX]*RecordList
}

type tableField struct {

	// 根据字段里的各种数值创建的索引
	equalMapper map[interface{}]*RecordList

	etcMapper map[interface{}]*unequalData
}

// 添加数据到字段, 索引, 引用数据所在结构体
func (self *tableField) Add(data, refRecord interface{}) {

	var value *RecordList

	if exists, ok := self.equalMapper[data]; ok {
		value = exists
	} else {
		value = newRecordList()

		self.equalMapper[data] = value
	}

	value.Add(refRecord)
}

func (self *tableField) addIndexData(t matchType, key int32, list *RecordList) {

	var ud *unequalData

	if self.etcMapper == nil {
		self.etcMapper = make(map[interface{}]*unequalData)
	}

	if v, ok := self.etcMapper[key]; ok {
		ud = v
	} else {
		ud = &unequalData{}
		self.etcMapper[key] = ud
	}

	ud.matchTypeList[t] = list

}

func (self *tableField) KeyCount() int {
	return len(self.equalMapper)
}

// 向结果集添加数据
func (self *tableField) All(q *Query) {

	for _, v := range self.equalMapper {
		if !q.addRecord(v) {
			return
		}
	}

}

func (self *tableField) getByKey(key interface{}) *RecordList {
	if v, ok := self.equalMapper[key]; ok {

		return v
	}

	return nil
}

// 向结果集添加符合条件的数据
func (self *tableField) Match(q *Query, t matchType, data interface{}) {

	switch t {
	case matchType_Equal:

		if v := self.getByKey(data); v != nil {
			q.addRecord(v)
		}

	case matchType_NotEqual:

		if !self.matchByIndex(q, t, data) {
			for k, v := range self.equalMapper {

				if k != data {
					if !q.addRecord(v) {
						break
					}
				}
			}
		}

	default:

		// 使用索引过的数据
		if self.matchByIndex(q, t, data) {
			return
		}

		// 暴力匹配
		for k, v := range self.equalMapper {

			if compare(t, k, data) {
				if !q.addRecord(v) {
					break
				}
			}

		}

	}

}

func compare(t matchType, tabData, userExpect interface{}) bool {

	switch tabDataT := tabData.(type) {
	case int32:
		userExpectT := userExpect.(int32)

		switch t {
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

func (self *tableField) matchByIndex(q *Query, t matchType, data interface{}) bool {

	if self.etcMapper == nil {
		return false
	}

	// 这个数值对应的各种操作符映射数据
	if v, ok := self.etcMapper[data]; ok {

		// 找出这个操作符的缩影
		typeList := v.matchTypeList[t]

		if typeList == nil {
			panic("match type index not built: " + getSignByMatchType(t))
		}

		q.addRecord(typeList)
	}

	return true
}

func newTableField() *tableField {
	return &tableField{
		equalMapper: make(map[interface{}]*RecordList),
	}
}
