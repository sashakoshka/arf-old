package analyzer

import "github.com/sashakoshka/arf/parser"

type Analyzer struct {
        moduleName string
        module    *parser.Module

        tree *SemanticTree
        
        warnCount  int
        errorCount int
}

func Analyze (
        moduleName string,
) (
        tree      *SemanticTree,
        warnCount  int,
        errorCount int,
) {
        moduleItem, exists := parser.GetModule(moduleName, false)
        if !exists { panic("module " + moduleName + " does not exist") }
        module := moduleItem.GetModule()

        analyzer := Analyzer {
                moduleName: moduleName,
                tree:      &SemanticTree {
                        typedefs:  make(map[SectionSpecifier] *Typedef),
                        datas:     make(map[SectionSpecifier] *Data),
                        functions: make(map[SectionSpecifier] *Function),
                },
        }

        // TODO
        // 1. analyze type definitions

        _,
        moduleTypedefs,
        _ := module.GetSections()

        for typedefName, _ := range moduleTypedefs {
                analyzer.analyzeTypedef(analyzer.moduleName, typedefName)
        }
        
        // 2. analyze data sections
        // 3. analyze functions

        // when doing these things, check if another module is being used. if it
        // is, skim-parse the module and add it to the cache. then, recursively
        // analyze and resolve the referenced item and add it to the cache.

        return tree, analyzer.warnCount, analyzer.errorCount
}

func (analyzer *Analyzer) analyzeTypedef (
        moduleName  string,
        typedefName string,
) (
        typedef *Typedef,
) {
        typedef = analyzer.tree.GetTypedef(moduleName, typedefName)
        if typedef != nil { return }
        
        typedef = &Typedef {
                module: moduleName,
        }

        // TODO: search cache for module
        moduleTypedef := 

        typedef.SetName(moduleTypedef.GetName())
        typedef.SetInternalPermission(moduleTypedef.GetInternalPermission())
        typedef.SetExternalPermission(moduleTypedef.GetExternalPermission())

        // TODO: analyze inherits
        moduleTypedefInherits := moduleTypedef.GetType()
        analyzer.analyzeType(moduleName, moduleTypedefInherits)
        
        // TODO: analyze members

        analyzer.tree.setTypedef(moduleName, moduleTypedef.GetName(), typedef)
        return
}

func (analyzer *Analyzer) analyzeType (
        moduleName string,
        moduleType parser.Type,
) (
        what *Type,
) {
        what = &Type { }

        var name    parser.Identifier
        var points *parser.Type
        
        name, points,
        what.items,
        what.mutable = moduleType.GetTypeData()

        if (points != nil) {
                what.points = analyzer.analyzeType(moduleName, *points)
                return
        }

        // TODO: analyze typedef specified by name. if nil, or name length is
        // greater than 2, fail
        trail := name.GetTrail()
        if (len(trail) > 2) {
                name.PrintError()
                return nil
        }

        return
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
