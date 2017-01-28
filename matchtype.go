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
	MatchType_MAX
)

var sign2MatchType = map[string]MatchType{
	"==": MatchType_Equal,
	"!=": MatchType_NotEqual,
	">":  MatchType_Great,
	">=": MatchType_GreatEqual,
	"<":  MatchType_Less,
	"<=": MatchType_LessEqual,
}

var matchType2Sign = map[MatchType]string{}

func getMatchTypeBySign(sign string) MatchType {

	if v, ok := sign2MatchType[sign]; ok {
		return v
	}

	return MatchType_Unknown
}

func getSignByMatchType(t MatchType) string {
	if v, ok := matchType2Sign[t]; ok {
		return v
	}

	return "unknown"
}

func init() {

	for k, v := range sign2MatchType {
		matchType2Sign[v] = k
	}
}
