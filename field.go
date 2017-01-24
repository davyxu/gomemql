package gomemql

type tableField struct {

	// 根据字段里的各种数值创建的索引
	mapper map[interface{}]*RecordList
}

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

func (self *tableField) All(ml *Query) {

	for _, v := range self.mapper {
		addListToResult(ml, v)
	}

}

func (self *tableField) Match(ml *Query, t MatchType, data interface{}) {

	switch t {
	case MatchType_Equal:

		if v, ok := self.mapper[data]; ok {

			addListToResult(ml, v)
		}

	case MatchType_NotEqual:

		for k, v := range self.mapper {

			if k != data {
				addListToResult(ml, v)
			}
		}

	default:
		vdata := data.(int32)

		for k, v := range self.mapper {

			key := k.(int32)

			switch t {
			case MatchType_Great:
				if key > vdata {
					addListToResult(ml, v)
				}
			case MatchType_GreatEqual:
				if key >= vdata {
					addListToResult(ml, v)
				}
			case MatchType_Less:
				if key < vdata {
					addListToResult(ml, v)
				}
			case MatchType_LessEqual:
				if key <= vdata {
					addListToResult(ml, v)
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
