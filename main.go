package main

import (
        "os"
        "fmt"
        "github.com/sashakoshka/arf/document"
)

func main () {
        module, nWarn, nError, err := document.Parse("tests/hello")
        if err != nil { os.Exit(1) }
        module.Dump()
        fmt.Println("(i)", nWarn, "warnings and", nError, "errors")
        
        // lines, nWarn, nError, err := lexer.Tokenize("tests/hello.arf", "several")
        // if err != nil { os.Exit(1) }
        // for _, line := range lines {
                // line.Dump()
        // }
        // fmt.Println("(i)", nWarn, "warnings and", nError, "errors")
}
