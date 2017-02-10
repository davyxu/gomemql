package gomemql_btree

import (
	"testing"

	"github.com/google/btree"
)

type tableDef struct {
	Id    int32
	Level int32
	Name  string
	Tag   int32
}

type tableDef_Name tableDef

func (self *tableDef_Name) Less(than btree.Item) bool {

	other := than.(*tableDef_Name)

	if self.Name != other.Name {
		return self.Name < other.Name
	}

	return false
}

type tableDef_2Field tableDef

func (self *tableDef_2Field) Less(than btree.Item) bool {

	other := than.(*tableDef_2Field)

	if self.Tag != other.Tag {
		return self.Tag < other.Tag
	}

	if self.Name != other.Name {
		return self.Name < other.Name
	}

	return false
}

type tableDef_Level tableDef

func (self *tableDef_Level) Less(than btree.Item) bool {

	other := than.(*tableDef_Level)

	if self.Level != other.Level {
		return self.Level < other.Level
	}

	return false
}

// 数据源
var tabData = []*tableDef{
	&tableDef{Id: 6, Level: 20, Name: "kitty", Tag: 1},
	&tableDef{Id: 1, Level: 50, Name: "hello", Tag: 2},
	&tableDef{Id: 4, Level: 20, Name: "kitty", Tag: 2},
	&tableDef{Id: 5, Level: 10, Name: "power", Tag: 2},
	&tableDef{Id: 3, Level: 20, Name: "hello", Tag: 1},
	&tableDef{Id: 2, Level: 10, Name: "kitty", Tag: 1},
}

// 最简单例子
func TestHelloWorld(t *testing.T) {

	f := NewField()

	// 将Name记录添加到字段中
	for _, v := range tabData {

		// 将原类型, 转换为Name类型, 以匹配btree.Item接口
		f.AddRecord((*tableDef_Name)(v), v)
	}

	var result []int32

	NewQuery(func(el interface{}) {

		record := el.(*tableDef)

		result = append(result, record.Id)

		t.Log(el)

		// 对Name字段进行匹配
	}).Equal(f, &tableDef_Name{
		Name: "hello",
	},
	).Start()

	/*
		&{1 50 hello 2}
		&{3 20 hello 1}
	*/

	if len(result) != 2 || result[0] != 1 || result[1] != 3 {
		t.FailNow()
	}

}

func Test2Condition(t *testing.T) {

	f1 := NewField()

	for _, v := range tabData {

		f1.AddRecord((*tableDef_2Field)(v), v)
	}

	f2 := NewField()

	for _, v := range tabData {

		f2.AddRecord((*tableDef_Level)(v), v)
	}

	var result []int32

	NewQuery(func(el interface{}) {

		record := el.(*tableDef)

		result = append(result, record.Id)

		t.Log(el)

		// 两个结构体字段同时匹配
	}).Equal(f1, &tableDef_2Field{
		Name: "kitty",
		Tag:  1,
	},
	// 小于匹配
	).Less(f2, &tableDef_Level{
		Level: 20,
	},
	).Start()

	/*
		&{2 10 kitty 1g}
	*/

	if len(result) != 1 || result[0] != 2 {
		t.FailNow()
	}

}
