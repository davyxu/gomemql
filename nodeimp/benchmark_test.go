package gomemql

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

	tab := NewTable()
	for _, r := range tabData {
		tab.AddRecord(r.Id, r.Level)
	}

	// 构建Id字段>的0~100的索引索引
	tab.GenIndexGreat(0, 0, 100)

	b.ResetTimer()
	// 并发查询量
	for i := 0; i < 3000; i++ {
		NewQuery(tab).Great(int32(50)).Equal(int32(500)).Result(func(interface{}) bool {
			return true
		})
	}
}
