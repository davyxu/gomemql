package gomemql_map

import "testing"

type tableDef struct {
	Id    int32
	Level int32
	Name  string
}

func genTestTable() *Table {
	// 数据源
	tabData := []*tableDef{
		&tableDef{Id: 6, Level: 20, Name: "kitty"},
		&tableDef{Id: 1, Level: 50, Name: "hello"},
		&tableDef{Id: 4, Level: 20, Name: "kitty"},
		&tableDef{Id: 5, Level: 10, Name: "power"},
		&tableDef{Id: 3, Level: 20, Name: "hello"},
		&tableDef{Id: 2, Level: 20, Name: "kitty"},
	}

	// 创建数据表
	tab := NewTable(new(tableDef))
	for _, r := range tabData {
		tab.AddRecord(r)
	}

	return tab
}

func Test2Condition(t *testing.T) {

	tab := genTestTable()

	// ====================例子1====================
	// 2条件匹配查询

	NewQuery(tab).Less("Level", int32(50)).Equal("Name", "hello").Result(func(v interface{}) bool {

		t.Log(v)

		return true
	})

	// Got  &{3 20 hello}

}

func TestShowAll(t *testing.T) {

	tab := genTestTable()

	// 直接访问结果,无缓存, 效率高, 但不能处理SortBy和Limit

	NewQuery(tab).Result(func(v interface{}) bool {
		t.Log(v)
		return true
	})

	/*
		Got All 6 records
	*/

}

func TestGenIndex(t *testing.T) {

	tab := genTestTable()
	if err := tab.GenFieldIndexNotEqual("Id", 1, 6); err != nil {
		t.Log(err)
		t.FailNow()
	}

	// 索引创建
	NewQuery(tab).NotEqual("Id", int32(3)).Result(func(v interface{}) bool {

		t.Log(v)

		return true
	})

	/*
				Got
		        &{4 20 kitty}
		        &{5 10 power}
		        &{6 20 kitty}
	*/
}
