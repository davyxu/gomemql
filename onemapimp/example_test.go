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
	{Id: 6, Level: 20, Name: "kitty"},
	{Id: 1, Level: 50, Name: "hello"},
	{Id: 4, Level: 20, Name: "kitty"},
	{Id: 5, Level: 10, Name: "power"},
	{Id: 3, Level: 20, Name: "hello"},
	{Id: 2, Level: 20, Name: "kitty"},
}

func TestHelloWorld(t *testing.T) {

	tab := NewTable()

	for _, v := range tabData {

		tab.AddRecord(v, v.Name, v.Level)

	}

	result := NewQuery(tab).Equal("kitty").Equal(int32(20)).Result()
	for _, r := range result {
		t.Log(r)
	}

	t.Log(tab.String())
}

func TestGreatIndex(t *testing.T) {

	tab := NewTable()

	for _, v := range tabData {

		tab.AddRecord(v, v.Name, v.Id)

	}

	tab.GenIndexGreat(1, 1, 6)

	result := NewQuery(tab).Equal("hello").Great(int32(2)).Result()
	for _, r := range result {
		t.Log(r)
	}

}

func TestMultiIndex(t *testing.T) {

	tab := NewTable()

	tab.AddRecord("a", int32(1), int32(3))
	tab.AddRecord("b", int32(1), int32(3))

	tab.GenIndexGreatEqual(0, 1, 3)
	tab.GenIndexLessEqual(1, 1, 3)

	t.Log(tab.String())

	result := NewQuery(tab).GreatEqual(int32(1)).LessEqual(int32(1)).Result()
	for _, r := range result {
		t.Log(r)
	}

}
