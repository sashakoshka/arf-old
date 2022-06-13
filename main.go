package main

import (
        "os"
        "fmt"
        "github.com/sashakoshka/arf/parser"
        "github.com/sashakoshka/arf/analyzer"
)

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
        module.Dump()

        analyzerWarnings, analyzerErrors, err := analyzer.Analyze(module)
        totalWarnings += analyzerWarnings
        totalErrors   += analyzerErrors
        
        fmt.Println("(i)", totalWarnings, "warnings and", totalErrors, "errors")
        
        // lines, nWarn, nError, err := lexer.Tokenize("tests/extranum.arf", "extranum")
        // if err != nil { os.Exit(1) }
        // for _, line := range lines {
                // line.Dump()
        // }
        // fmt.Println("(i)", nWarn, "warnings and", nError, "errors")
}
