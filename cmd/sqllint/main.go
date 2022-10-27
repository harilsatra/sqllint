package main

import (
	"github.com/harilsatra/sqllint/sqlcheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(sqlcheck.Analyzer)
}
