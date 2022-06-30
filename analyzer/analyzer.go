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
        // 1. analyze type definitions
        // 2. analyze data sections
        // 3. analyze functions

        // when doing these things, check if another module is being used. if it
        // is, skim-parse the module and add it to the cache. then, recursively
        // analyze and resolve the referenced item and add it to the cache.

        return
}
