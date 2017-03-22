package gomemql

import "fmt"

type condition struct {
	index int
	value interface{}

	t matchType
}

func (self *condition) String() string {

	return fmt.Sprintf("index:%d op:%s value:%v", self.index, getSignByMatchType(self.t), self.value)
}
