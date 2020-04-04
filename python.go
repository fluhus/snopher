package main

// Deals with creating python source.

import (
	"bytes"
	"strings"
	"text/template"
)

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
)

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
