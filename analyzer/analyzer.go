package analyzer

import "github.com/sashakoshka/arf/parser"

type Analyzer struct {
        warnCount  int
        errorCount int
}

func Analyze (
        module *parser.Module,
) (
        tree       SemanticTree,
        warnCount  int,
        errorCount int,
) {
        analyzer := Analyzer {
                
        }

        // TODO
        // 1. analyze type definitions
        // 2. analyze data sections
        // 3. analyze functions

        // when doing these things, check if another module is being used. if it
        // is, skim-parse the module and add it to the cache. then, recursively
        // analyze and resolve the referenced item and add it to the cache.

        return tree, analyzer.warnCount, analyzer.errorCount
}

func (analyzer *Analyzer) PrintWarning (
        reporter parser.ErrorReporter,
        cause ...interface { },
) {
        reporter.PrintWarning(cause)
}

func (analyzer *Analyzer) PrintError (
        reporter parser.ErrorReporter,
        cause ...interface { },
) {
        reporter.PrintError(cause)
}

func (analyzer *Analyzer) PrintFatal (
        reporter parser.ErrorReporter,
        cause ...interface { },
) {
        reporter.PrintFatal(cause)
}
