package gomemql_btree

type matchType int

// 数据匹配模式
const (
	matchType_Unknown = iota
	matchType_Equal
	matchType_Great
	matchType_GreatEqual
	matchType_Less
	matchType_LessEqual
)
