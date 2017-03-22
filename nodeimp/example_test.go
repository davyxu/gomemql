package gomemql

import (
	"testing"
)

type tableDef struct {
	Id    int32
	Level int32
	Name  string
}

var tabData = []*tableDef{
	&tableDef{Id: 6, Level: 20, Name: "kitty"},
	&tableDef{Id: 1, Level: 50, Name: "hello"},
	&tableDef{Id: 4, Level: 20, Name: "kitty"},
	&tableDef{Id: 5, Level: 10, Name: "power"},
	&tableDef{Id: 3, Level: 20, Name: "hello"},
	&tableDef{Id: 2, Level: 20, Name: "kitty"},
}

func TestHelloWorld(t *testing.T) {

	tab := NewTable()

	for _, v := range tabData {
		tab.AddRecord(v.Name, v)
	}

	// 匹配Name为hello
	NewQuery(tab).Equal("hello").Result(func(v interface{}) bool {

		t.Log(v)

		return true
	})

}

//func TestTemp(t *testing.T) {

//	tab := NewTable()

//	for i := 1; i < 2; i++ {
//		tab.AddRecord(int32(1), int32(2), i*100)
//	}

//	// 构建第二个字段(Id), 从1~6的索引
//	tab.GenIndexGreatEqual(0, 1, 2)
//	tab.GenIndexLessEqual(1, 1, 2)

//	//tab.Print()

//	NewQuery(tab).GreatEqual(int32(1)).LessEqual(int32(1)).Result(func(v interface{}) bool {

//		t.Log(v)

//		return true
//	})

//}

func Test2ConditionWithIndex(t *testing.T) {

	tab := NewTable()

	for _, v := range tabData {
		tab.AddRecord(v.Id, v.Name, v)
	}

	// 构建第二个字段(Id), 从1~6的索引
	tab.GenIndexGreatEqual(0, 1, 6)

	NewQuery(tab).GreatEqual(int32(4)).Equal("kitty").Result(func(v interface{}) bool {

		t.Log(v)

		return true
	})

}

func TestVariantFieldLen(t *testing.T) {

	tab := NewTable()

	tab.AddRecord(1, "Genji")

	tab.AddRecord(1, 2, "Zenyatta")

	NewQuery(tab).Equal(1).Result(func(v interface{}) bool {

		t.Log(1, v)

		return true
	})

	NewQuery(tab).Equal(1).Equal(2).Result(func(v interface{}) bool {

		t.Log(2, v)

		return true
	})

}
