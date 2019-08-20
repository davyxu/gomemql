# gomemql

基于内存结构的多条件组合查询器

# 用途

* 内存表格数据查询

* 游戏触发器条件查询
成就表: 定义成就类型, 事件类型, 玩家等级等静态表格数据
通过本系统查出符合条件的集合, 再检查动态数据, 例如: 玩家拥有物品等


# onemap实现版
https://github.com/davyxu/gomemql/onemapimp

* 原生golang编写,无cgo, 无第三方引用, 不依赖sqlite

* 多字段任意组合查询

* 支持构建字段搜索索引, 提高不等匹配(!=, <,>...)查询性能

* 建立索引后,查询性能为O(字段数). 非缓冲的不等匹配字段性能为O(字段数*字段记录数), 等于匹配性能为O(字段数)

* 记录相关无内存分配, GC友好

# 实现原理

根据每个记录集的N个字段, 建立N层的树状节点

所有记录的字段数必须统一

每个节点根据字段值进行索引

最终节点持有匹配前面节点条件的结果集合

# 例子
```golang
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


```



# 其他版本

C#版参见: https://github.com/davyxu/MemQLSharp

# 备注

感觉不错请star, 谢谢!

知乎: [http://www.zhihu.com/people/sunicdavy](http://www.zhihu.com/people/sunicdavy)
