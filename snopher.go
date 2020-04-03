package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("Hi")
	f := "/home/amitmit/Desktop/gopyt/foo.h"
	a, err := readRawLines(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	b, err := tokenizeLines(a)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	c, _ := parseTokens(b)
	for _, x := range c {
		fmt.Println(x)
	}
}

var (
	indent  = "    "
	libName = "_mylib"
	cTypes  = map[string]string{
		"string":  "GoString",
		"int":     "ctypes.c_longlong",
		"int64":   "ctypes.c_longlong",
		"int32":   "ctypes.c_int",
		"int16":   "ctypes.c_short",
		"uint":    "ctypes.c_ulonglong",
		"uint64":  "ctypes.c_ulonglong",
		"uint32":  "ctypes.c_uint",
		"uint16":  "ctypes.c_ushort",
		"float64": "ctypes.c_double",
		"float32": "ctypes.c_float",
		"bool":    "ctypes.c_bool",
		"":        "None",
	}
	pyTypes = map[string]string{
		"string":  "str",
		"int":     "int",
		"int64":   "int",
		"int32":   "int",
		"int16":   "int",
		"uint":    "int",
		"uint64":  "int",
		"uint32":  "int",
		"uint16":  "int",
		"float64": "float",
		"float32": "float",
		"bool":    "bool",
		"":        "None",
	}

	goFuncRe = regexp.MustCompile("^func (\\w+)\\((.*)\\)(.*)$")
	cFuncRe  = regexp.MustCompile("^(\\w+) (\\w+)\\((.*)\\)$")
)

func readRawLines(file string) ([]string, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(src), "\n")
	for _, line := range lines {
		lines = lines[1:]
		if strings.HasPrefix(line, "extern \"C\" {") {
			break
		}
	}
	var result []string
	for _, line := range lines {
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "extern ") {
			result = append(result, strings.Trim(line, "\r"))
		}
	}
	return result, nil
}

func tokenizeLines(lines []string) ([]*token, error) {
	var tokens []*token
	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			tokens = append(tokens, &token{"comment", strings.Trim(line[2:], " \t")})
		} else if strings.HasPrefix(line, "extern ") {
			tokens = append(tokens, &token{"function", line[7 : len(line)-1]})
		} else {
			return nil, fmt.Errorf("unrecognized line value: %v", line)
		}
	}
	return tokens, nil
}

func goSignature(pkg, fun string) (string, error) {
	cmd := exec.Command("go", "doc", "-u", pkg, fun)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v: %v", err, string(out))
	}
	// TODO(amit): Handle cases where function is not the first line.
	return strings.Split(string(out), "\n")[0], nil
}

func parseGoSignature(sig string) (*function, error) {
	// Break down into main parts.
	match := goFuncRe.FindStringSubmatch(sig)
	if match == nil {
		return nil, fmt.Errorf("signature %q does not match expected pattern",
			sig)
	}
	name, paramsRaw, outType := match[1], match[2], match[3]

	// Parse params.
	paramsStr := splitParams(paramsRaw)
	// Handle params where several names share the same type.
	for i := len(paramsStr) - 1; i >= 0; i-- {
		p := paramsStr[i]
		if len(p) == 1 {
			p = append(p, paramsStr[i+1][1])
		}
		if len(p) != 2 {
			return nil, fmt.Errorf("bad number of parts (%v) in param: %v",
				len(p), p)
		}
	}

	var params []*param
	for _, p := range paramsStr {
		params = append(params, &param{p[0], p[1]})
	}

	return &function{name, outType, params, ""}, nil
}

func parseCSignature(sig string, comment []string) (*function, error) {
	// Break down into main parts.
	match := cFuncRe.FindStringSubmatch(sig)
	if match == nil {
		return nil, fmt.Errorf("signature %q does not match expected pattern",
			sig)
	}
	outType, name, paramsRaw := match[1], match[2], match[3]
	paramsStr := splitParams(paramsRaw)

	var params []*param
	for _, p := range paramsStr {
		// TODO(amit): Check param length before accessing.
		params = append(params, &param{p[0], p[1]})
	}

	return &function{name, outType, params, strings.Join(comment, "\n")}, nil
}

func splitParams(params string) [][]string {
	if params == "" {
		return nil
	}
	parts := strings.Split(params, ", ")
	var result [][]string
	for _, p := range parts {
		result = append(result, strings.Split(p, " "))
	}
	return result
}

func parseTokens(tokens []*token) ([]*function, error) {
	var result []*function
	var commentLines []string
	for _, t := range tokens {
		if t.typ == "comment" {
			commentLines = append(commentLines, t.data)
			continue
		}
		f, err := parseCSignature(t.data, commentLines)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
		commentLines = nil
	}
	return result, nil
}

/*
def parse_tokens(tokens: List[Tuple]) -> List[Func]:
    result = []
    for i, token in enumerate(tokens):
        if token[0] != 'function':
            continue
        if i > 0 and tokens[i - 1][0] == 'comment':
            result.append(parse_c_signature(token[1], tokens[i - 1][1]))
        else:
            result.append(parse_c_signature(token[1], []))
    return result
*/

type token struct {
	typ  string
	data string
}

type param struct {
	name string
	typ  string
}

type function struct {
	name    string
	typ     string
	params  []*param
	comment string
}

var srcTemplate = `import ctypes


class GoString(ctypes.Structure):
    _fields_ = [('p', ctypes.c_char_p), ('n', ctypes.c_int)]


def go_string(s):
    enc = s.encode()
    return GoString(enc, len(enc))


def from_go_string(s):
    return s.p[:s.n].decode()


{{.lib}} = ctypes.CDLL({{.dll}})
`
