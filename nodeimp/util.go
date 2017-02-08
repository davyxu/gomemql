package gomemql

import "bytes"

func writeLineWithIndent(n *indexNode, b *bytes.Buffer, str string) {

	for i := 0; i < n.index*4; i++ {
		b.WriteString(" ")
	}

	b.WriteString(str)

	b.WriteString("\n")

}
