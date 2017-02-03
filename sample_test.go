package gomemql

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
	for _, v := range NewQuery(tab).Where("Level", "<", int32(50)).Where("Name", "==", "hello").Result() {

		t.Log(v)
	}

	// Got  &{3 20 hello}

}

func TestSortLimit(t *testing.T) {

	tab := genTestTable()

	// ====================例子2====================
	// 1条件, 排序和数量限制
	for _, v := range NewQuery(tab).Where("Level", "==", int32(20)).SortBy(func(x, y interface{}) bool {
		a := x.(*tableDef)
		b := y.(*tableDef)

		if a.Id != b.Id {
			return a.Id < b.Id
		}

		return false
	}).Limit(3).Result() {

		t.Log(v)
	}

	/*
		Got
		&{3 20 hello}
		&{4 20 kitty}
		&{6 20 kitty}
	*/
}

func TestShowAll(t *testing.T) {

	tab := genTestTable()

	// 直接访问结果,无缓存, 效率高, 但不能处理SortBy和Limit

	NewQuery(tab).VisitRawResult(func(v interface{}) bool {
		t.Log(v)
		return true
	})

	/*
		Got All 6 records
	*/

}

func TestGenIndex(t *testing.T) {

	tab := genTestTable()
	if err := tab.GenFieldIndex("Id", "!=", 1, 6); err != nil {
		t.Log(err)
		t.FailNow()
	}

	// 索引创建
	for _, v := range NewQuery(tab).Where("Id", "!=", int32(3)).Result() {

		t.Log(v)
	}

	/*
			Got
	        &{4 20 kitty}
	        &{5 10 power}
	        &{6 20 kitty}
	*/
}
