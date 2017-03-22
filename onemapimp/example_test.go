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
}
