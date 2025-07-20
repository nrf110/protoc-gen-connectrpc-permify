package util

import (
	"strings"
)

const indent = 4

func Indent(nestingLevel int) string {
	return strings.Repeat(" ", indent*nestingLevel)
}
