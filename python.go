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
		"argtypes":   argTypes,
		"pytype":     func(s string) string { return pyTypes[s] },
		"funcinputs": funcInputs,
	}).Parse(srcTemplate))

	buf := &bytes.Buffer{}
	err := t.Execute(buf, map[string]interface{}{
		"dll":   dllfile,
		"funcs": funcs,
	})
	return buf.Bytes(), err
}

func (p *param) ArgType() string {
	if strings.HasPrefix(p.Typ, "[]") {
		_, ok := cTypes[p.Typ[2:]]
		if !ok {
			panic("unsupported type: " + p.Typ)
		}
		return "GoSlice"
	}
	t, ok := cTypes[p.Typ]
	if !ok {
		panic("unsupported type: " + p.Typ)
	}
	return t
}

func (p *param) ResType() string {
	t, ok := cTypes[p.Typ]
	if !ok {
		panic("unsupported type: " + p.Typ)
	}
	return t
}

func (p *param) PyType() string {
	if strings.HasPrefix(p.Typ, "[]") {
		t, ok := pyTypes[p.Typ[2:]]
		if !ok {
			panic("unsupported type: " + p.Typ)
		}
		return "List[" + t + "]"
	}
	t, ok := pyTypes[p.Typ]
	if !ok {
		panic("unsupported type: " + p.Typ)
	}
	return t
}

func (p *param) FuncInput() string {
	// TODO(amit): Add a condition for string slices.
	if strings.HasPrefix(p.Typ, "[]") {
		return "go_slice(" + p.Name + ", " + cTypes[p.Typ[2:]] + ")"
	}
	if p.Typ == "string" {
		return "go_string(" + p.Name + ")"
	}
	return p.Name
}

func argTypes(params []*param) string {
	var types []string
	for _, p := range params {
		types = append(types, p.ArgType())
	}
	return strings.Join(types, ", ")
}

func funcInputs(params []*param) string {
	var names []string
	for _, p := range params {
		names = append(names, p.FuncInput())
	}
	return strings.Join(names, ", ")
}

var srcTemplate = `
{{- /* BOILERPLATE */ -}}

import ctypes
from typing import List


{{lib}}: ctypes.CDLL = None


class GoString(ctypes.Structure):
{{indent}}_fields_ = [('p', ctypes.c_char_p), ('n', ctypes.c_int)]


class GoSlice(ctypes.Structure):
{{indent}}_fields_ = [('data', ctypes.c_void_p),
{{indent}}            ('len', ctypes.c_longlong),
{{indent}}            ('cap', ctypes.c_longlong)]


def go_string(s: str) -> GoString:
{{indent}}enc = s.encode()
{{indent}}return GoString(enc, len(enc))


def from_go_string(s: GoString) -> str:
{{indent}}return s.p[:s.n].decode()


def go_slice(arr: List, ctype) -> GoSlice:
{{indent}}p = ctypes.cast((ctype * len(arr))(*arr), ctypes.c_void_p)
{{indent}}return GoSlice(p, len(arr), len(arr))


{{/* INIT FUNCTION */ -}}

def init(dll_path: str) -> None:
{{indent}}global {{lib}}
{{indent}}{{lib}} = ctypes.CDLL(dll_path)

{{/* FUNCTION TYPE INITIALIZATION */ -}}

{{range .funcs -}}
{{indent}}{{lib}}.{{.Name}}.argtypes = [{{argtypes .Params}}]
{{indent}}{{lib}}.{{.Name}}.restype = {{ctype .Typ}}

{{end}}


{{- /* FUNCTIONS */ -}}

{{range .funcs}}

{{- /* FUNCTION SIGNATURE */}}
def {{.Name}}(
{{- range $i, $p := .Params -}}
{{if gt $i 0}}, {{end -}}
{{.Name}}: {{.PyType}}
{{- end -}}
) -> {{pytype .Typ}}:

{{- /* FUNCTION CONTENT */}}
{{if .Comment}}{{indent}}"""{{.Comment}}"""
{{end -}}
{{indent -}}
return {{if eq .Typ "string"}}from_go_string({{end -}}
{{lib}}.{{.Name}}({{funcinputs .Params}})
{{- if eq .Typ "string"}}){{end}}

{{end -}}
`
