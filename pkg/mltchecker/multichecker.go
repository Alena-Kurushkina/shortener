package mltchecker

import (
	"strings"

	"github.com/ichiban/thelper"
	"github.com/kyoh86/exportloopref"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/Alena-Kurushkina/shortener/pkg/exitanalyzer"
)

func NewMultichecker() []*analysis.Analyzer {
	mychecks := []*analysis.Analyzer{
		exitanalyzer.ExitCheckAnalyzer,

		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}

	// определяем map подключаемых правил
	checks := map[string]bool{
		"S1003":  true,
		"ST1005": true,
		"QF1010": true,
	}

	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		if strings.HasPrefix(v.Analyzer.Name, "SA") || checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	mychecks = append(mychecks, exportloopref.Analyzer)
	mychecks = append(mychecks, thelper.Analyzer)

	return mychecks
}
