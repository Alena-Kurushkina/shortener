package exitanalyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
    Name: "osexitcheck",
    Doc:  `
check for os.Exit call in main function of main package

It returns warning if it finds os.Exit() call in main function of main package.

For example:
  package main

  import (
	"os"
  )
  
  func main(){
    os.Exit(3) 
  }

`,
    Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspectMainMethod:=func(node ast.Node){
		ast.Inspect(node, func (mnode ast.Node) bool {
			if call,ok:=mnode.(*ast.CallExpr); ok {
				if fun,ok:=call.Fun.(*ast.SelectorExpr); ok{
					if x,ok:=fun.X.(*ast.Ident); ok{
						if x.Name=="os" && fun.Sel.Name=="Exit"{
							pass.Reportf(fun.Sel.NamePos, "call os.Exit")
						}
					}
				}
			}
			return true
		})
	}

	inspectMainPackage:=func(node ast.Node){
		ast.Inspect(node, func (mnode ast.Node) bool {
			// check if main method
			if mmethod,ok:=mnode.(*ast.FuncDecl); ok {
				if mmethod.Name.Name == "main" {
					inspectMainMethod(mnode)
					// main method is found, no sense to search futher
					return false
				}							
			}						
			return true
		})
	}

    for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(filename, ".go") {
			continue
		}
        ast.Inspect(file, func(node ast.Node) bool {
			//check if main package
			if f,ok:=node.(*ast.File); ok{
				if f.Name.Name=="main" {
					inspectMainPackage(node)				
					// main package is found, no sense to search futher
					return false					
				}
			}            
            return true
        })
    }
	
    return nil, nil
} 