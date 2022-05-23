package Model

import (
	"strings"
)

func (e *Element) HtmlTag(content string) string {
	up := strings.ToUpper(content)
	first := strings.Index(up, e.Start) + len(e.Start)
	last := strings.Index(up[first:], e.End)
	if first != -1 && last != -1 {
		return content[first : first+last]
	}
	return ""
}
