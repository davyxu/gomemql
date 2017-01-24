package gomemql

//type unequalData struct{
//	[MatchType_MAX]*RecordList
//}

type tableField struct {

	// 根据字段里的各种数值创建的索引
	equalMapper map[interface{}]*RecordList

	lessMapper map[interface{}]*RecordList
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

func (self *tableField) mapperByType(t MatchType) map[interface{}]*RecordList {

	switch t {
	case MatchType_Less:
		return self.lessMapper
	}

	return nil
}

func (self *tableField) addIndexData(t MatchType, key int32, list *RecordList) {

	indexMap := self.mapperByType(t)

	indexMap[key] = list

}

func (self *tableField) genIndex(t MatchType) {
	switch t {
	case MatchType_Less:
		if self.lessMapper != nil {
			panic("less index already exists")
		}

		self.lessMapper = make(map[interface{}]*RecordList)
	}
}

func (self *tableField) KeyCount() int {
	return len(self.equalMapper)
}

func addListToResult(ml *Query, rl *RecordList) {

	for _, v := range rl.data {
		ml.add(v)
	}
}

// 向结果集添加数据
func (self *tableField) All(q *Query) {

	for _, v := range self.equalMapper {
		addListToResult(q, v)
	}

}

func (self *tableField) getByKey(key interface{}) *RecordList {
	if v, ok := self.equalMapper[key]; ok {

		return v
	}

	return nil
}

// 向结果集添加符合条件的数据
func (self *tableField) Match(q *Query, t MatchType, data interface{}) {

	switch t {
	case MatchType_Equal:

		if v := self.getByKey(data); v != nil {

			addListToResult(q, v)
		}

	case MatchType_NotEqual:

		for k, v := range self.equalMapper {

			if k != data {
				addListToResult(q, v)
			}
		}

	default:

		if self.lessMapper != nil && t == MatchType_Less {

			if v, ok := self.lessMapper[data]; ok {
				addListToResult(q, v)
			}

		} else {
			vdata := data.(int32)
			for k, v := range self.equalMapper {

				key := k.(int32)

				switch t {
				case MatchType_Great:
					if key > vdata {
						addListToResult(q, v)
					}
				case MatchType_GreatEqual:
					if key >= vdata {
						addListToResult(q, v)
					}
				case MatchType_Less:
					if key < vdata {
						addListToResult(q, v)
					}
				case MatchType_LessEqual:
					if key <= vdata {
						addListToResult(q, v)
					}
				}
			}

		}

	}

}

func newTableField() *tableField {
	return &tableField{
		equalMapper: make(map[interface{}]*RecordList),
	}
}
