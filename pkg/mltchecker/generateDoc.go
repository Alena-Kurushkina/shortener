//go:build ignore

package mltchecker

import (
	"os"
	"strings"
	"unicode"
)

func GenerateDoc() {
	doc:=strings.Builder{}
	doc.WriteString(`
Util staticlint is a multichecker consisting of:
• printf, shadow and structtag analyzers of golang.org/x/tools/go/analysis/passes package
• custom ExitCheckAnalyzer
• all analyzers of SA type from staticcheck.io package
• S1003, ST1005, QF1010 analyzers of staticcheck.io package
• exportloopref.Analyzer and thelper.Analyzer.

Usage: 
  cmd/staticlint/main.go [path to folder for check]
	
Example:
  ./cmd/staticlint/main.go ./...

Description of analyzers:

`)

	checks:=NewMultichecker()
	for _, a := range checks {
		doc.WriteString(strings.ToUpper(a.Name))
		doc.WriteString("\n\n")
		doc.WriteString(loverFirstLetter(a.Doc))
		doc.WriteString("\n\n")
	}
	docStr:= strings.Replace(doc.String(), "\n", "\n // ", -1)

	file, err := os.OpenFile("./doc.go", os.O_CREATE|os.O_RDWR, 0775)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	file.Write([]byte(docStr))
	file.Write([]byte("\n"))
	file.Write([]byte("package main"))
}

func loverFirstLetter(str string) string {
	runeStr:=[]rune(str)
	if unicode.IsUpper(runeStr[0]){
		runeStr[0] = unicode.ToLower(runeStr[0])
		return string(runeStr)
	}
	return str
}