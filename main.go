package main

import (
        "os"
        "fmt"
        "github.com/sashakoshka/arf/parser"
)

func main () {
        module, nWarn, nError, err := parser.Parse("tests/extranum")
        if err != nil { os.Exit(1) }
        module.Dump()
        fmt.Println("(i)", nWarn, "warnings and", nError, "errors")
        
        // lines, nWarn, nError, err := lexer.Tokenize("tests/extranum.arf", "extranum")
        // if err != nil { os.Exit(1) }
        // for _, line := range lines {
                // line.Dump()
        // }
        // fmt.Println("(i)", nWarn, "warnings and", nError, "errors")
}
