package gomemql_map

import "testing"

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

	tab := NewTable(new(tableDef))
	for _, r := range tabData {
		tab.AddRecord(r)
	}

	// 构建Id字段>的0~100的索引索引
	tab.GenFieldIndexGreat("Id", 0, 100)

	b.ResetTimer()
	// 并发查询量
	for i := 0; i < 3000; i++ {
		NewQuery(tab).Great("Id", int32(50)).Equal("Level", int32(500)).Result(nil)
	}
}
