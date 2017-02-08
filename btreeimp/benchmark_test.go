package gomemql_btree

import (
	"testing"

	"github.com/google/btree"
)

type tableDef_Id tableDef

func (self *tableDef_Id) Less(than btree.Item) bool {

	other := than.(*tableDef_Id)

	if self.Id != other.Id {
		return self.Id < other.Id
	}

	return false
}

func BenchmarkTest(b *testing.B) {

	var tabData []*tableDef

	// 预估每个索引包含的记录
	for i := 0; i < 1000; i++ {
		tabData = append(tabData, &tableDef{
			Id:    int32(i + 1),
			Level: int32(i * 10),
			Name:  "kitty",
		})
	}

	f1 := NewField()

	for _, v := range tabData {

		f1.AddRecord((*tableDef_Id)(v), v)
	}

	f2 := NewField()

	for _, v := range tabData {

		f2.AddRecord((*tableDef_Level)(v), v)
	}

	b.ResetTimer()
	// 并发查询量
	for i := 0; i < 3000; i++ {

		NewQuery(func(el interface{}) {

		}).Great(f1, &tableDef_Id{
			Id: 50,
		},
		).Equal(f2, &tableDef_Level{
			Level: 500,
		},
		).Start()

	}
}
