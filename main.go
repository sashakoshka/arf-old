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

        var totalWarnings int
        var totalErrors   int
        
        module,
        parserWarnings,
        parserErrors,
        err := parser.Parse(os.Args[1], false)
        
        totalWarnings += parserWarnings
        totalErrors   += parserErrors
        if err != nil { os.Exit(1) }

        tree, analyzerWarnings, analyzerErrors := analyzer.Analyze(module)
        totalWarnings += analyzerWarnings
        totalErrors   += analyzerErrors

        tree.Dump()
        
        fmt.Println("(i)", totalWarnings, "warnings and", totalErrors, "errors")
}
