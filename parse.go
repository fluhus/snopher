package main

// Deals with parsing the input package.

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

var (
	goFuncRe = regexp.MustCompile("^func (\\w+)\\((.*)\\)\\s*(.*)$")
	cFuncRe  = regexp.MustCompile("^(\\w+) (\\w+)\\((.*)\\)$")
)

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
			if strings.HasPrefix(line, "extern struct ") {
				return nil, fmt.Errorf("multiple return values are not supported (yet): %q", line)
			}
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
		return nil, fmt.Errorf("parseGoSignature: signature %q does not match expected pattern",
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
		return nil, fmt.Errorf("parseCSignature: signature %q does not match expected pattern",
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
