package main

// Contains main function and common utilities.

import (
	"fmt"
	"io/ioutil"
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

type token struct {
	Typ  string
	Data string
}

type param struct {
	Name string
	Typ  string
}

type function struct {
	Name    string
	Typ     string
	Params  []*param
	Comment string
}
