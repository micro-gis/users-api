package string_util

import "strings"

func IsEmptyString(str string) bool {
	return strings.TrimSpace(str) == ""
}
