package main

import (
	"fmt"
	"log"
	"os"
)

type IndentedTextWriter struct {
	f           *os.File
	indentLevel int
	tabString   string
	tabsPending bool
}

func NewIndentedTextWriter(f *os.File, tabString string) *IndentedTextWriter {
	return &IndentedTextWriter{f: f, tabString: tabString, indentLevel: 0, tabsPending: false}
}

func (g *IndentedTextWriter) Indent() {
	g.indentLevel++
}

func (g *IndentedTextWriter) Outdent() {
	if g.indentLevel > 0 {
		g.indentLevel--
	} else {
		g.indentLevel = 0
	}
}

func (g *IndentedTextWriter) PrintTabs() {
	if g.tabsPending {
		for i := 0; i < g.indentLevel; i++ {
			_, err := g.f.WriteString(g.tabString)
			if err != nil {
				log.Fatal(err)
			}
		}
		g.tabsPending = false
	}
}

func (g *IndentedTextWriter) Println(a ...any) {
	g.PrintTabs()
	for _, s := range a {
		_, err := g.f.WriteString(fmt.Sprint(s))
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err := g.f.WriteString("\n")
	if err != nil {
		log.Fatal(err)
	}
	g.tabsPending = true
}

func (g *IndentedTextWriter) Printf(format string, a ...any) {
	g.PrintTabs()
	s := fmt.Sprintf(format, a...)
	_, err := g.f.WriteString(s)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *IndentedTextWriter) Print(a ...any) {
	g.PrintTabs()
	for _, s := range a {
		_, err := g.f.WriteString(fmt.Sprint(s))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *IndentedTextWriter) Close() {
	g.f.Close()
}
