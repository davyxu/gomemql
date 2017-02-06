# gomemql

基于内存结构的多条件组合查询器

# 用途

* 内存表格数据查询

* 游戏触发器条件查询
成就表: 定义成就类型, 事件类型, 玩家等级等静态表格数据
通过本系统查出符合条件的集合, 再检查动态数据, 例如: 玩家拥有物品等

# map实现版
https://github.com/davyxu/gomemql/mapimp

## 特性

* 原生golang编写,无cgo, 无第三方引用, 不依赖sqlite

* 多字段任意组合查询

* 支持构建字段搜索索引, 提高不等匹配(!=, <,>...)查询性能

* 建立缓冲后,查询性能为O(1). 非缓冲的不等匹配字段性能为O(N*M), 等于匹配性能为O(1)

## 适用范围

类sql方式的复杂组合查询, gc不敏感项目

## 支持功能的等效SQL语法
select * from tableData where condition1 and condition2...

## 例子
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
	NewQuery(tab).Less("Level", int32(50)).Equal("Name", "hello").Result(func(v interface{}) bool {

		t.Log(v)

		return true
	})

	t.Log()

	// Got  &{3 20 hello}

	// ====================例子3====================
	// 直接访问结果,无缓存, 效率高, 但不能处理SortBy和Limit

	NewQuery(tab).Result(func(v interface{}) bool {
		t.Log(v)
		return true
	})

	/*
		Got All 6 records
	*/

```
# btree实现版

https://github.com/davyxu/gomemql/btreeimp

依赖
https://github.com/google/btree

## 特性

* 结构更简单

* 更低的GC, 更好的内存控制

* 原生golang编写,无cgo, 无第三方引用, 不依赖sqlite

* 可以按相等/不相等条件, 字段组合分离任意组合查询

* 查询性能为O(logN), 不等匹配为O(N*M)

## 适用范围
gc敏感, 性能稍好的逻辑

## 支持功能的等效SQL语法
select * from tableData where condition1 and condition2...

## 例子
```golang

	type tableDef struct {
		Id    int32
		Level int32
		Name  string
		Tag   int32
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


	var tabData = []*tableDef{
		&tableDef{Id: 6, Level: 20, Name: "kitty", Tag: 1},
		&tableDef{Id: 1, Level: 50, Name: "hello", Tag: 2},
		&tableDef{Id: 4, Level: 20, Name: "kitty", Tag: 2},
		&tableDef{Id: 5, Level: 10, Name: "power", Tag: 2},
		&tableDef{Id: 3, Level: 20, Name: "hello", Tag: 1},
		&tableDef{Id: 2, Level: 10, Name: "kitty", Tag: 1},
	}

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

```


# 其他版本

C#版参见: https://github.com/davyxu/MemQLSharp

# 备注

感觉不错请star, 谢谢!

博客: http://www.cppblog.com/sunicdavy

知乎: http://www.zhihu.com/people/sunicdavy

邮箱: sunicdavy@qq.com
