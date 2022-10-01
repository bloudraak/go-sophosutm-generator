package main

import (
    "strings"

    "github.com/iancoleman/strcase"
)

type Classes []*Class

func (c Classes) Len() int {
    return len(c)
}

func (c Classes) Less(i, j int) bool {
    is := strcase.ToCamel(c[i].Name())
    js := strcase.ToCamel(c[j].Name())
    return strings.Compare(is, js) < 0
}

func (c Classes) Swap(i, j int) {
    c[i], c[j] = c[j], c[i]
}
