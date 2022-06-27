package analyzer

import "github.com/sashakoshka/arf/parser"

type Analyzer struct {
        
}

func Analyze (
        module *parser.Module,
) (
        warnCount int,
        errorCount int,
        err error,
) {
        // analyzer := Analyzer {
                // 
        // }

        // TODO
        // 0. skim over requires. also store them in a map so that if a require
        //    needs to be skimmed again by another file or module then there u
        //    go
        // 1. analyze type definitions
        // 2. analyze data sections
        // 3. analyze functions

        return
}
