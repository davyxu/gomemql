package gomemql

type MatchType int

// 数据匹配模式
const (
	MatchType_Unknown = iota
	MatchType_Equal
	MatchType_NotEqual
	MatchType_Great
	MatchType_GreatEqual
	MatchType_Less
	MatchType_LessEqual
)

var sign2MatchType = map[string]MatchType{
	"==": MatchType_Equal,
	"!=": MatchType_NotEqual,
	">":  MatchType_Great,
	">=": MatchType_GreatEqual,
	"<":  MatchType_Less,
	"<=": MatchType_LessEqual,
}

func getMatchTypeBySign(sign string) MatchType {

	if v, ok := sign2MatchType[sign]; ok {
		return v
	}

	return MatchType_Unknown
}
