package cf

import (
	"regexp"
)

func init() {
	variableRegex = regexp.MustCompile("\\$\\{(.+)\\}")
}

var variableRegex *regexp.Regexp

func variableReference(v interface{}) (string, bool) {
	if vt, ok := v.(string); ok {
		if vmatch := variableRegex.FindSubmatch([]byte(vt)); vmatch != nil {
			return string(vmatch[1]), true
		}
	}
	return "", false
}
