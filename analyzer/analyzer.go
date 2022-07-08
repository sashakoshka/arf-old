package analyzer

import "github.com/sashakoshka/arf/parser"

type Analyzer struct {
        moduleName string

        tree *SemanticTree
        
        warnCount  int
        errorCount int
}

func Analyze (
        module *parser.Module,
) (
        tree      *SemanticTree,
        warnCount  int,
        errorCount int,
) {
        analyzer := Analyzer {
                moduleName: module.GetName(),
                tree:      &SemanticTree {
                        typedefs:  make(map[string] *Typedef),
                        datas:     make(map[string] *Data),
                        functions: make(map[string] *Function),
                },
        }

        // TODO
        // 1. analyze type definitions

        _,
        moduleTypedefs,
        _ := module.GetSections()

        typedefChecklist := Checklist { }
        for _, moduleTypedef := range moduleTypedefs {
                if typedefChecklist.IsDone(moduleTypedef.GetName()) {
                        continue
                }

                analyzer.analyzeTypedef(analyzer.moduleName, moduleTypedef)
        }
        
        // 2. analyze data sections
        // 3. analyze functions

        // when doing these things, check if another module is being used. if it
        // is, skim-parse the module and add it to the cache. then, recursively
        // analyze and resolve the referenced item and add it to the cache.

        return tree, analyzer.warnCount, analyzer.errorCount
}

func (analyzer *Analyzer) analyzeTypedef (
        moduleName     string,
        moduleTypedef *parser.Typedef,
) {
        typedef := Typedef {
                module: moduleName,
        }

        typedef.SetName(moduleTypedef.GetName())
        typedef.SetInternalPermission(moduleTypedef.GetInternalPermission())
        typedef.SetExternalPermission(moduleTypedef.GetExternalPermission())

        // TODO: analyze inherits
        // TODO: analyze members

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
