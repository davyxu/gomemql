package gomemql

type tableField struct {

	// 根据字段里的各种数值创建的索引
	mapper map[interface{}]*RecordList
}

// 添加数据到字段, 索引, 引用数据所在结构体
func (self *tableField) Add(data, refRecord interface{}) {

	var value *RecordList

	if exists, ok := self.mapper[data]; ok {
		value = exists
	} else {
		value = newRecordList()

		self.mapper[data] = value
	}

	value.Add(refRecord)
}

func addListToResult(ml *Query, rl *RecordList) {

	for _, v := range rl.data {
		ml.add(v)
	}
}

// 向结果集添加数据
func (self *tableField) All(q *Query) {

	for _, v := range self.mapper {
		addListToResult(q, v)
	}

}

// 向结果集添加符合条件的数据
func (self *tableField) Match(q *Query, t MatchType, data interface{}) {

	switch t {
	case MatchType_Equal:

		if v, ok := self.mapper[data]; ok {

			addListToResult(q, v)
		}

	case MatchType_NotEqual:

		for k, v := range self.mapper {

			if k != data {
				addListToResult(q, v)
			}
		}

	default:
		vdata := data.(int32)

		for k, v := range self.mapper {

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

	return

}

func newTableField() *tableField {
	return &tableField{
		mapper: make(map[interface{}]*RecordList),
	}
}
