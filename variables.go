package cf

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

func init() {
	variableRegex = regexp.MustCompile("^\\$\\{([a-zA-Z0-9_.]+)\\}$")
	inlineVariableRegex = regexp.MustCompile("\\$\\{([a-zA-Z0-9_.]+)\\}")
}

var variableRegex *regexp.Regexp
var inlineVariableRegex *regexp.Regexp

func variableReference(v interface{}) (string, bool) {
	if vt, ok := v.(string); ok {
		if vmatch := variableRegex.FindSubmatch([]byte(strings.TrimSpace(vt))); vmatch != nil {
			return string(vmatch[1]), true
		}
	}
	return "", false
}

func inlineVariablesFound(in string) bool {
	vmatch := inlineVariableRegex.FindSubmatch([]byte(strings.TrimSpace(in)))
	return vmatch != nil
}

func replaceInlineVariables(in string, opt *Options) (string, error) {
	out := strings.TrimSpace(in)
	vmatch := inlineVariableRegex.FindSubmatchIndex([]byte(out))
	for vmatch != nil {
		vname := out[vmatch[2]:vmatch[3]]
		vvalue, resolved := opt.resolveVariable(vname)
		if !resolved {
			return "", errors.Errorf("variable ${%s} not found", vname)
		}
		vvstring, ok := vvalue.(string)
		if !ok {
			return "", errors.Errorf("variable ${%s} not string value", vname)
		}
		out = out[:vmatch[0]]+vvstring+out[vmatch[1]:]

		vmatch = inlineVariableRegex.FindSubmatchIndex([]byte(out))
	}
	return out, nil
}
