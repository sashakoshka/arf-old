package parser

import (
        "os"
        "fmt"
        "path"
        "bufio"
        "errors"
        "strings"
        "io/ioutil"
        "github.com/sashakoshka/arf/lexer"
        "github.com/sashakoshka/arf/lineFile"
        "github.com/sashakoshka/arf/validate"
)

var (
        errEmptyModule   = errors.New("there are no files in this module")
        errEmptyFile     = errors.New("file is devoid of content")
        errBadIndent     = errors.New("this line should not be indented")
        errTooMuchIndent = errors.New("this line is indented too far")
        errSurpriseEOF   = errors.New("file terminated unexpectedly")
        errSurpriseEOL   = errors.New("line terminated unexpectedly")
        errNotArf        = errors.New("not an arf file, expected :arf")
)

/* Parser is a magic machine that turns a path into a parsed AST. Neato!
 */
type Parser struct {
        directory  string
        file       *lineFile.LineFile

        lines      []*lexer.Line
        lineIndex  int
        line       *lexer.Line

        token      *lexer.Token
        tokenIndex int
        
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
        skim bool,
) (
        module     *Module,
        warnCount  int,
        errorCount int,
        err        error,
) {
        moduleDir  := path.Dir(modulePath)
        moduleBase := path.Base(modulePath)
        fmt.Println("...", "parsing module \"" + moduleBase + "\"")

        parser := &Parser {
                directory: moduleDir,
                module:    &Module {
                        name:      moduleBase,
                        functions: make(map[string] *Function),
                        typedefs:  make(map[string] *Typedef),
                        datas:     make(map[string] *Data),
                },
        }

        if !validate.ValidateName(moduleBase) {
                err = errors.New (
                        "\"" + moduleBase + "\" is not a valid module name")
                parser.printGeneralFatal(err)
                return nil, 0, 1, err
        }

        candidates, err := ioutil.ReadDir(parser.directory)
        if err != nil {
                parser.printFatal(err)
                return parser.module, 0, 1, err
        }

        foundFile := false
        for _, candidate := range candidates {
                if candidate.IsDir() { continue }
                filePath := moduleDir + "/" + candidate.Name()
                if getModuleName(filePath) != parser.module.name { continue }

                fmt.Println("(i)", "found file", filePath)
                foundFile = true

                // attempt to parse the file. if any part fails, go on to the
                // next one.
                err = parser.parseFile(filePath, skim)
                if err != nil { continue }
        }

        if !foundFile {
                parser.printGeneralFatal(errEmptyModule)
                return nil, 0, 1, errEmptyModule
        }

        fmt.Println(".//", "module parsed")
        return parser.module, parser.warnCount, parser.errorCount, nil
}

/* parseFile parses a specific file into the module.
 */
func (parser *Parser) parseFile (filePath string, skim bool) (err error) {
        // open file safely
        parser.file, err = lineFile.Open(filePath, parser.module.name)
        if err != nil {
                parser.printFatal(err)
                return
        }
        
        lines, nWarn, nError, err := lexer.Tokenize (
                parser.file,
                parser.module.name)

        if err != nil { return err }
        if len(lines) == 0 {
                parser.printFatal(errEmptyFile)
                return errEmptyFile
        }

        parser.lines = lines
        parser.warnCount  += nWarn
        parser.errorCount += nError
        
        parser.lineIndex = 0
        parser.line = parser.lines[parser.lineIndex]
        parser.tokenIndex = 0
        parser.token = parser.line.Tokens[parser.tokenIndex]

        // parse metadata
        err = parser.parseMeta()
        if err != nil {
                parser.printFatal(err)
                return err
        }

        // parse body
        err = parser.parseBody(skim)
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

/* nextLine advances the parser to the next line, and returns whether or not the
 * end of the file was reached.
 */
func (parser *Parser) nextLine () (done bool) {
        parser.lineIndex ++
        parser.tokenIndex = 0
        if parser.endOfFile() {
                parser.line = nil
                return true
        }
        parser.line = parser.lines[parser.lineIndex]
        parser.token = parser.line.Tokens[parser.tokenIndex]
        return false
}

func (parser *Parser) endOfFile () (eof bool) {
        return parser.lineIndex >= len(parser.lines)
}

/* expect takes in a number of token kinds. It advances the parser, and returns
 * true if it matches what is expected. Otherwise, it prints an error and
 * returns false. If there are no kinds supplied, it will return true only on
 * end of line.
 */
func (parser *Parser) expect (kinds ...lexer.TokenKind) (match bool) {
        currentKind := lexer.TokenKindNone

        if !parser.endOfLine() {
                currentKind = parser.token.Kind
        }

        if len(kinds) == 0 {
                if parser.endOfLine() {
                        return true
                } else {
                        parser.printError (
                                parser.token.Column,
                                "unexpected token, expected end of line")
                        return false
                }
        }

        if !parser.endOfFile() {
                for _, kind := range kinds {
                        if currentKind == kind {
                                return true
                        }
                }
        }

        errText := "unexpected "
        errColumn := parser.token.Column

        if parser.endOfFile() {
                errText += "end of file"
                errColumn = parser.lines[len(parser.lines) - 1].GetLength()
        } else if parser.endOfLine() {
                errText += "end of line"
                errColumn = parser.line.EndColumn
        } else {
                errText += currentKind.ToString() + " token"
        }
        errText += "."

        // print out what was expected, if there are less than 6 expected items.
        if (len(kinds) < 6) {
                errText += " expected "
                if len(kinds) > 1 {
                        for _, kind := range kinds[:len(kinds) - 1] {
                                errText += kind.ToString() + ", "
                        }

                        errText += "or "
                }
                
                errText += kinds[len(kinds) - 1].ToString()
        }

        parser.printError(errColumn, errText)
        return false
}

/* nextToken advances the parser to the next token, and returns whether or not
 * the end of the line was reached.
 */
func (parser *Parser) nextToken () (done bool) {
        if parser.endOfFile() { return true }

        parser.tokenIndex ++
        if parser.endOfLine() {
                parser.token = &lexer.Token { Kind: lexer.TokenKindNone }
                return true
        }
        parser.token = parser.line.Tokens[parser.tokenIndex]
        return false
}

func (parser *Parser) endOfLine () (eol bool) {
        if parser.line == nil { return false }
        return parser.tokenIndex >= len(parser.line.Tokens)
}

func (parser *Parser) getCurrentRealRow () (row int) {
        var line *lexer.Line

        if parser.endOfFile() {
                line = parser.lines[len(parser.lines) - 1]
        } else {
                line = parser.line
        }

        return line.Row
}

func (parser *Parser) printWarning (column int, cause ...interface {}) {
        parser.warnCount ++
        parser.file.PrintWarning(column, parser.getCurrentRealRow(), cause...)
}

func (parser *Parser) printError (column int, cause ...interface {}) {
        parser.errorCount ++
        parser.file.PrintError(column, parser.getCurrentRealRow(), cause...)
}

func (parser *Parser) printFatal (err error) {
        parser.errorCount ++
        parser.file.PrintFatal(err)
}

func (parser *Parser) printGeneralFatal (err error) {
        parser.errorCount ++
        fmt.Println (
                "\033[31mXXX\033[0m",
                "\033[90min\033[0m",
                parser.module.name)
        fmt.Println("   ", err)
}

/* embedPosition
 * Returns a position object reflecting the current position of the parser that
 * can be embedded into a struct.
 */
func (parser *Parser) embedPosition () (position Position) {
        return Position {
                Column: parser.token.Column,
                Row:    parser.line.Row,
                File:   parser.file,
        }
}
