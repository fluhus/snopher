// Generates the code sections in index.md.
package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	index, err := os.ReadFile("index.md")
	die(err)
	index = generateGo(index)
	index = generatePython(index)
	os.WriteFile("index.md", index, 0o644)
}

// Inserts go code from the source files into the markdown text.
func generateGo(index []byte) []byte {
	sectionRE := regexp.MustCompile("(?s)<!-- gen:[^\n]+\\.go -->\n+```go\n.*?\n```")
	fileRE := regexp.MustCompile("<!-- gen:([^\n]+\\.go) -->")
	return sectionRE.ReplaceAllFunc(index, func(b []byte) []byte {
		f := fileRE.FindAllSubmatch(b, -1)[0][1]
		src, err := readGoFile(string(f))
		die(err)
		return []byte(fmt.Sprintf("<!-- gen:%s -->\n\n```go\n%s\n```",
			f, src))
	})
}

// Inserts python code from the source files into the markdown text.
func generatePython(index []byte) []byte {
	sectionRE := regexp.MustCompile("(?s)<!-- gen:[^\n]+\\.py -->\n+```python\n.*?\n```")
	fileRE := regexp.MustCompile("<!-- gen:([^\n]+\\.py) -->")
	return sectionRE.ReplaceAllFunc(index, func(b []byte) []byte {
		f := fileRE.FindAllSubmatch(b, -1)[0][1]
		src, err := readPyFile(string(f))
		die(err)
		return []byte(fmt.Sprintf("<!-- gen:%s -->\n\n```python\n%s\n```",
			f, src))
	})
}

// Reads a go file and filters out the boilerplate.
func readGoFile(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	data = regexp.MustCompile("(//.*)?package main\n+").ReplaceAll(data, nil)
	data = regexp.MustCompile("func main\\(\\) \\{\\}").ReplaceAll(data, nil)
	data = regexp.MustCompile("\n\n+").ReplaceAll(data, []byte("\n\n"))
	data = regexp.MustCompile("^\n+").ReplaceAll(data, nil)
	data = regexp.MustCompile("\n+$").ReplaceAll(data, nil)
	return data, nil
}

// Reads a python file and filters out the boilerplate.
func readPyFile(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	data = regexp.MustCompile("from [^\n]+ import [^\n]+").ReplaceAll(data, nil)
	data = regexp.MustCompile("import [^\n]+").ReplaceAll(data, nil)
	data = regexp.MustCompile("\n\n+").ReplaceAll(data, []byte("\n\n"))
	data = regexp.MustCompile("^\n+").ReplaceAll(data, nil)
	data = regexp.MustCompile("\n+$").ReplaceAll(data, nil)
	return data, nil
}

// Panics if error.
func die(err error) {
	if err != nil {
		panic(err)
	}
}
