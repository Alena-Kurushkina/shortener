package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/Alena-Kurushkina/shortener/pkg/mltchecker"
)

func main(){
	mychecks := mltchecker.NewMultichecker()

    multichecker.Main(
        mychecks...,
    )
}