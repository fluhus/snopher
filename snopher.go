package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

func main() {
	fmt.Println("Hi")
	pkg := "foo"
	f := "foo.h"
	dll := "./foo.so"
	err := processPkg(pkg, f, dll, "foo2.py")
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}

var (
	indent  = "    "
	libName = "_lib"
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

	goFuncRe = regexp.MustCompile("^func (\\w+)\\((.*)\\)\\s*(.*)$")
	cFuncRe  = regexp.MustCompile("^(\\w+) (\\w+)\\((.*)\\)$")
)

func processPkg(pkg, hfile, dllfile, pyfile string) error {
	funcs, err := parseHFile(hfile)
	if err != nil {
		return err
	}
	err = addGodocInfo(funcs, pkg)
	if err != nil {
		return err
	}
	src, err := generatePythonSource(funcs, dllfile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pyfile, src, 0644)
}

func generatePythonSource(funcs []*function, dllfile string) ([]byte, error) {
	t := template.Must(template.New("python-src").Funcs(template.FuncMap{
		"lib":        func() string { return libName },
		"indent":     func() string { return indent },
		"ctype":      func(s string) string { return cTypes[s] },
		"ctypes":     paramCTypes,
		"pytype":     func(s string) string { return pyTypes[s] },
		"paramnames": paramNames,
	}).Parse(srcTemplate))

	buf := &bytes.Buffer{}
	err := t.Execute(buf, map[string]interface{}{
		"dll":   dllfile,
		"funcs": funcs,
	})
	return buf.Bytes(), err
}

func parseHFile(file string) ([]*function, error) {
	lines, err := readRawLines(file)
	if err != nil {
		return nil, err
	}
	tokens, err := tokenizeLines(lines)
	if err != nil {
		return nil, err
	}
	funcs, err := parseTokens(tokens)
	if err != nil {
		return nil, err
	}
	return funcs, nil
}

func addGodocInfo(funcs []*function, pkg string) error {
	for _, f := range funcs {
		sig, err := goSignature(pkg, f.Name)
		if err != nil {
			return err
		}
		gofunc, err := parseGoSignature(sig)
		if err != nil {
			return err
		}
		f.Params = gofunc.Params
		f.Typ = gofunc.Typ
	}
	return nil
}

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
			paramsStr[i] = append(p, paramsStr[i+1][1])
		}
		if len(paramsStr[i]) != 2 {
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

	return &function{name, outType, params, strings.Join(comment, "\n"+indent)}, nil
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
		if t.Typ == "comment" {
			commentLines = append(commentLines, t.Data)
			continue
		}
		f, err := parseCSignature(t.Data, commentLines)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
		commentLines = nil
	}
	return result, nil
}

func paramCTypes(params []*param) string {
	var types []string
	for _, p := range params {
		types = append(types, cTypes[p.Typ])
	}
	return strings.Join(types, ", ")
}

func paramNames(params []*param) string {
	var names []string
	for _, p := range params {
		if p.Typ == "string" {
			names = append(names, "to_go_string("+p.Name+")")
		} else {
			names = append(names, p.Name)
		}
	}
	return strings.Join(names, ", ")
}

type token struct {
	Typ  string
	Data string
}

type param struct {
	Name string
	Typ  string
}

func (p *param) String() string {
	return fmt.Sprint(*p)
}

type function struct {
	Name    string
	Typ     string
	Params  []*param
	Comment string
}

var srcTemplate = `
{{- /* BOILERPLATE */ -}}

import ctypes


class GoString(ctypes.Structure):
    _fields_ = [('p', ctypes.c_char_p), ('n', ctypes.c_int)]


def to_go_string(s):
    enc = s.encode()
    return GoString(enc, len(enc))


def from_go_string(s):
    return s.p[:s.n].decode()


{{lib}} = ctypes.CDLL({{printf "%q" .dll}})

{{/* FUNCTION TYPE INITIALIZATION */ -}}
{{range .funcs}}
{{lib}}.{{.Name}}.argtypes = [{{ctypes .Params}}]
{{lib}}.{{.Name}}.restype = {{ctype .Typ}}
{{end}}

{{- /* FUNCTIONS */ -}}

{{""}}
{{range .funcs}}

{{- /* FUNCTION SIGNATURE */}}
def {{.Name}}(
{{- range $i, $p := .Params -}}
{{if gt $i 0}}, {{end -}}
{{.Name}}: {{pytype .Typ}}
{{- end -}}
) -> {{pytype .Typ}}:

{{- /* FUNCTION CONTENT */}}
{{if .Comment}}{{indent}}"""{{.Comment}}"""
{{end -}}
{{indent -}}
return {{if eq .Typ "string"}}from_go_string({{end -}}
{{lib}}.{{.Name}}({{paramnames .Params}})
{{- if eq .Typ "string"}}){{end}}

{{end}}
`
