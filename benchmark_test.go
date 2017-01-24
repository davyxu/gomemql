package gomemql

import "testing"

func TestBenchmark(t *testing.T) {

	var tabData []*tableDef

	// 预估每个索引包含的记录
	for i := 0; i < 200; i++ {
		tabData = append(tabData, &tableDef{
			Id:    int32(i + 1),
			Level: int32(i * 10),
			Name:  "kitty",
		})
	}

	tab := NewTable(new(tableDef))
	for _, r := range tabData {
		tab.AddRecord(r)
	}

	// 并发查询量
	for i := 0; i < 1000; i++ {
		NewQuery(tab).Where("Id", ">", int32(50)).Where("Level", "==", int32(500)).VisitRawResult(nil)
	}
}
