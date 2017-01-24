# gomemql

基于内存结构的多条件组合查询器

# 用途

* 内存表格数据查询

* 游戏触发器条件查询
成就表: 定义成就类型, 事件类型, 玩家等级等静态表格数据
通过本系统查出符合条件的集合, 再检查动态数据, 例如: 玩家拥有物品等

# 特性

* 原生golang编写,无cgo, 无第三方引用

* 组合查询方便

# 支持功能的等效SQL语法
select * from memory where condition1 and condition2... limit count

```golang
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

	// ====================例子1====================
	// 2条件匹配查询
	for _, v := range NewQuery(tab).Where("Level", "<", int32(50)).Where("Name", "==", "hello").Result() {

		t.Log(v)
	}

	t.Log()

	// Got  &{3 20 hello}

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
	t.Log()
	// ====================例子3====================
	// 直接访问结果,无缓存, 效率高, 但不能处理SortBy和Limit

	NewQuery(tab).VisitRawResult(func(v interface{}) bool {
		t.Log(v)
		return true
	})

	/*
		Got All 6 records
	*/

```

# TODO
* 支持构建索引, 便于提高不等匹配(!=, <,>...)查询性能

# 备注

感觉不错请star, 谢谢!

博客: http://www.cppblog.com/sunicdavy

知乎: http://www.zhihu.com/people/sunicdavy

邮箱: sunicdavy@qq.com
