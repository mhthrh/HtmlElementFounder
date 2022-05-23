package Model

import (
	"strings"
)

func (e *Element) HeadCount(content string) int {
	return len(strings.Split(content, e.End))
}
