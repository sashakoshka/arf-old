package main

import "os"
import "fmt"
import "github.com/sashakoshka/arf/parser"
import "github.com/sashakoshka/arf/analyzer"

func main () {
        if (len(os.Args) < 2) {
                fmt.Println("specify module path")
                os.Exit(1)
        }
        
        module := parser.GetModule(os.Args[1], false)

        tree, _, _ := analyzer.Analyze(module.GetName())
        // totalWarnings += analyzerWarnings
        // totalErrors   += analyzerErrors

        tree.Dump()

        // TODO: query parser and analyzer module for total errors
        // fmt.Println("(i)", totalWarnings, "warnings and", totalErrors, "errors")
}
