package document

import (
        "os"
        "fmt"
        "path"
        "bufio"
        "errors"
        "strings"
        "io/ioutil"
        "github.com/sashakoshka/arf/lexer"
        "github.com/sashakoshka/arf/validate"
)

var (
        errSurpriseEOF = errors.New("file terminated unexpectedly")
        errNotArf      = errors.New("not an arf file, expected :arf")
)

/* Parser is a magic machine that turns a path into a parsed AST. Neato!
 */
type Parser struct {
        directory  string
        fileName   string

        lines      []*lexer.Line
        lineNumber int
        line       *lexer.Line
        
        module     *Module

        warnCount  int
        errorCount int
}

/* Parse takes in a module path, and returns a Module. The file at the end of
 * the path should not, as it is a virtual concept that is searched for in the
 * path's base directory. All files with a matching module feild are parsed into
 * the module that gets returned. It's like golang packages, except we are
 * calling them modules because we aren't insane.
 */
func Parse (
        modulePath string,
) (
        module     *Module,
        warnCount  int,
        errorCount int,
        err        error,
) {
        moduleDir  := path.Dir(modulePath)
        moduleBase := path.Base(modulePath)
        fmt.Println("...", "parsing module \"" + moduleBase + "\"")

        if !validate.ValidateName(moduleBase) {
                err = errors.New (
                        "\"" + moduleBase + "\" is not a valid module name")
                fmt.Println("XXX", err)
                return nil, 0, 0, err
        }

        parser := &Parser {
                directory: moduleDir,
                module:    &Module {
                        name:      moduleBase,
                        functions: make(map[string] Function),
                        typedefs:  make(map[string] Typedef),
                        datas:     make(map[string] Data),
                },
        }

        candidates, err := ioutil.ReadDir(parser.directory)
        if err != nil {
                parser.printFatal(err)
                return parser.module, parser.warnCount, parser.errorCount, err
        }

        for _, candidate := range candidates {
                if candidate.IsDir() { continue }
                filePath := moduleDir + "/" + candidate.Name()
                if getModuleName(filePath) != parser.module.name { continue }

                fmt.Println("(i)", "found file", filePath)

                // attempt to parse the file. if any part fails, go on to the
                // next one.
                err = parser.parseFile(filePath)
                if err != nil { continue }
        }

        fmt.Println(".//", "module parsed")
        return parser.module, parser.warnCount, parser.errorCount, nil
}

/* parseFile parses a specific file into the module.
 */
func (parser *Parser) parseFile (filePath string) (err error) {
        // open file safely
        parser.fileName = filePath
        
        lines, nWarn, nError, err := lexer.Tokenize (
                filePath,
                parser.module.name)

        if err != nil { return err }

        parser.lines = lines
        parser.warnCount  += nWarn
        parser.errorCount += nError

        // parse metadata
        err = parser.parseMeta()
        if err != nil {
                parser.printFatal(err)
                return err
        }

        // parse body
        err = parser.parseBody()
        if err != nil {
                parser.printFatal(err)
                return err
        }

        return nil
}

/* getModuleName takes in a file path (an actual one!) and returns the module
 * name that the file is a part of. If the file is not an arf file, it returns
 * an empty string.
 */
func getModuleName (filePath string) (name string) {
        // open file
        if path.Ext(filePath) != ".arf" { return "" }
        file, err := os.Open(filePath)
        defer file.Close()
        if err != nil { return "" }

        // look for magic bytes
        scanner := bufio.NewScanner(file)
        scanned := scanner.Scan()
        if !scanned                 { return "" }
        if scanner.Err()  != nil    { return "" }
        if scanner.Text() != ":arf" { return "" }

        // search for module field
        for scanner.Scan() {
                if strings.HasPrefix(scanner.Text(), "---") { return "" }
                fields := strings.Fields(scanner.Text())
                if len(fields) == 2 && fields[0] == "module" {
                        return fields[1]
                }
        }

        return ""
}

/* parseMeta parses the metadata header of an arf file. This contains the module
 * name, and other miscellaneous fields such as author and license.
 */
func (parser *Parser) parseMeta () (err error) {
        return
}

func (parser *Parser) parseBody () (err error) {
        return
}

/* nextLine advances the parser to the next line, and returns whether or not the
 * end of the file was reached.
 */
func (parser *Parser) nextLine () (done bool) {
        
        return
}

func (parser *Parser) printWarning (column int, cause... interface {}) {
        parser.warnCount ++
        fmt.Println("!!!", "in", parser.fileName, "of", parser.module.name)
        fmt.Println("   ", parser.lineNumber, ":",  column, parser.line)
        fmt.Println("   ", cause)
}

func (parser *Parser) printError (column int, cause... interface {}) {
        parser.errorCount ++
        fmt.Println("ERR", "in", parser.fileName, "of", parser.module.name)
        fmt.Println("   ", parser.lineNumber, ":",  column, parser.line)
        fmt.Println("   ", cause)
}

func (parser *Parser) printFatal (err error) {
        parser.errorCount ++
        fmt.Println("XXX", "in", parser.fileName, "of", parser.module.name)
        fmt.Println("   ", "could not parse module -", err)
}
