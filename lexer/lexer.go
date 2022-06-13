package lexer

import (
        "fmt"
        "math"
        "errors"
        "strings"
        "strconv"
        "github.com/sashakoshka/arf/lineFile"
        "github.com/sashakoshka/arf/validate"
)

/* TokenKind is an enum represzenting what type a token is.
 */
type TokenKind int

const (
        TokenKindNone TokenKind = iota
        TokenKindSeparator
        TokenKindDirection
        TokenKindPermission

        TokenKindInteger
        TokenKindSignedInteger
        TokenKindFloat
        TokenKindString
        TokenKindRune

        TokenKindName
        TokenKindSymbol

        TokenKindColon
        TokenKindDot
        
        TokenKindLBracket
        TokenKindRBracket
        TokenKindLBrace
        TokenKindRBrace
)

/* Lexer holds information about a current lexing operation. This struct is only
 * used within Tokenize().
 */
type Lexer struct {
        file       *lineFile.LineFile
        lines      []*Line
        lineNumber int
        line       *Line

        warnCount  int
        errorCount int
}

/* Line represents a line of a file. Its primary purpose is to store tokens.
 */
type Line struct {
        runes []rune
        index int

        // Row is the position of the line in the file. Again, this is for error
        // reporting.
        Row       int
        Column    int
        EndColumn int

        // Indent is the indentation level of the line. 
        Indent int

        // Tokens is an array of all tokens extracted from the line.
        Tokens []*Token
}

type Token struct {
        Kind        TokenKind
        StringValue string
        Value       interface {}
        Column      int
}

func Tokenize (
        file       *lineFile.LineFile,
        moduleName string,
) (
        lines []*Line,
        warnCount  int,
        errorCount int,
        err        error,
) {
        lexer := &Lexer { file: file }

        done := false
        for {
                done, err = lexer.tokenizeLine()
                if done || err != nil { break }
        }

        // TODO: set all runes in lines to nil to free memory
        return lexer.lines, lexer.warnCount, lexer.errorCount, err
}

func (lexer *Lexer) tokenizeLine () (done bool, err error) {
        for {
                var ignore bool
                done, ignore, err = lexer.nextLine()
                if done || err != nil { return }
                if !ignore { break }
        }

        line := lexer.line

        if line.Indent == 0 && strings.TrimSpace(string(line.runes)) == ":arf" {
                // magic bytes. ignore it
                return false, nil
        }

        for line.notEnd() {
                // make a crude guess at the token based on the first rune
                ch := line.ch()

                lowercase := ch >= 'a' && ch <= 'z'
                uppercase := ch >= 'A' && ch <= 'Z'
                number    := ch >= '0' && ch <= '9'
                
                if number {
                        lexer.tokenizeNumber(false)
                } else if lowercase || uppercase {
                        lexer.tokenizeMulti()
                } else {
                        switch ch {
                        case '"':
                                lexer.tokenizeString('"')
                                line.nextRune()
                                break
                        case '\'':
                                lexer.tokenizeString('\'')
                                line.nextRune()
                                break
                        case ':':
                                line.add(TokenKindColon, ":", ":")
                                line.nextRune()
                                break
                        case '.':
                                line.add(TokenKindDot, ".", ".")
                                line.nextRune()
                                break
                        case '[':
                                line.add(TokenKindLBracket, "[", "[")
                                line.nextRune()
                                break
                        case ']':
                                line.add(TokenKindRBracket, "]", "]")
                                line.nextRune()
                                break
                        case '{':
                                line.add(TokenKindLBrace, "{", "{")
                                line.nextRune()
                                break
                        case '}':
                                line.add(TokenKindRBrace, "}", "}")
                                line.nextRune()
                                break
                        default:
                                lexer.tokenizeSymbol()
                                break
                        }
                }

                lexer.skipWhitespace()
        }

        if len(line.Tokens) > 0 {
                lexer.lines = append(lexer.lines, line)
        }
        return
}

func (line *Line) ch () (ch rune) {
        if !line.notEnd() { return '\000' }
        return line.runes[line.index]
}

func (line *Line) nextRune () {
        line.index ++
}

func (line *Line) notEnd () (keepGoing bool) {
        return line.index < len(line.runes)
}

func (line *Line) add (
        kind        TokenKind,
        stringValue string,
        value       interface {},
) {
        line.Tokens = append (line.Tokens, &Token {
                Kind:        kind,
                StringValue: stringValue,
                Value:       value,
                Column:      line.index + line.Column,
        })
}

func (line *Line) addExisting (token *Token) {
        line.Tokens = append (line.Tokens, token)
}

func (line *Line) GetLength () (length int) {
        return len(line.runes)
}

func (line *Line) Dump () {
        var kind string

        for i := 0; i < line.Indent; i++ { fmt.Print("        ") }
        fmt.Println("line", line.Row)
        for _, token := range line.Tokens {
                switch token.Kind {
                case TokenKindSeparator:     kind = "Separator";     break
                case TokenKindDirection:     kind = "Direction";     break
                case TokenKindPermission:    kind = "Permission";    break
                case TokenKindInteger:       kind = "Integer";       break
                case TokenKindSignedInteger: kind = "SignedInteger"; break
                case TokenKindFloat:         kind = "Float";         break
                case TokenKindString:        kind = "String";        break
                case TokenKindRune:          kind = "Rune";          break
                case TokenKindName:          kind = "Name";          break
                case TokenKindSymbol:        kind = "Symbol";        break
                case TokenKindColon:         kind = "Colon";         break
                case TokenKindDot:           kind = "Dot";           break
                case TokenKindLBracket:      kind = "LBracket";      break
                case TokenKindRBracket:      kind = "RBracket";      break
                case TokenKindLBrace:        kind = "LBrace";        break
                case TokenKindRBrace:        kind = "RBrace";        break
                }

                for i := 0; i < line.Indent; i++ { fmt.Print("        ") }
                fmt.Println("-", kind, token.Value)
        }
}

func (lexer *Lexer) tokenizeNumber (negative bool) {
        line := lexer.line

        radix := 10
        isFloat := false
        floatPosition := 0

        token := Token { Column: line.index + line.Column }

        if line.ch() == '0' {
                line.nextRune()
                switch line.ch() {
                case 'x':
                        // hexidecimal
                        radix = 16
                        line.nextRune()
                        for line.notEnd() {
                                ch := line.ch()
                                notLower := ch < 'a' || ch > 'f'
                                notUpper := ch < 'A' || ch > 'F'
                                notNum   := ch < '0' || ch > '9'
                                dot := ch == '.'
                                if notLower && notUpper && notNum && !dot {
                                        break
                                }

                                line.nextRune()
                                if dot {
                                        isFloat = true
                                        floatPosition = len(token.StringValue)
                                } else {
                                        token.StringValue += string(ch)
                                }
                        }
                        break
                        
                case 'b':
                        // binary
                        radix = 2
                        line.nextRune()
                        for line.notEnd() {
                                ch := line.ch()
                                if ch < '0' || ch > '1' { break }
                                token.StringValue += string(ch)
                                line.nextRune()
                        }
                        break
                        
                default:
                        // octal
                        radix = 8
                        for line.notEnd() {
                                ch := line.ch()
                                notNum := ch < '0' || ch > '7'
                                dot := ch == '.'
                                if notNum && !dot { break }

                                line.nextRune()
                                if dot {
                                        isFloat = true
                                        floatPosition = len(token.StringValue)
                                } else {
                                        token.StringValue += string(ch)
                                }
                        }
                        break
                }
        } else {
                // decimal
                for line.notEnd() {
                        ch := line.ch()
                        notNum := ch < '0' || ch > '9'
                        dot := ch == '.'
                        if notNum && !dot { break }

                        line.nextRune()
                        if dot {
                                isFloat = true
                                floatPosition = len(token.StringValue)
                        } else {
                                token.StringValue += string(ch)
                        }
                }
        }

        parsedNumber, _ := strconv.ParseUint (
                token.StringValue, radix, 64)
        floatValue := float64(parsedNumber)
        if negative {
                token.Kind  = TokenKindSignedInteger
                token.Value = int64(parsedNumber) * -1
                floatValue *= -1
        } else {
                token.Kind  = TokenKindInteger
                token.Value = parsedNumber
        }

        if isFloat {
                token.Kind  = TokenKindFloat
                floatPosition = len(token.StringValue) - floatPosition
                coefficient := math.Pow(float64(radix), float64(floatPosition))
                floatValue /= coefficient
                token.Value = floatValue
        }
        
        line.addExisting(&token)
}

func (lexer *Lexer) tokenizeString (terminator rune) {
        line := lexer.line

        // TODO: create rune literal token kind
        token := Token {
                Kind:   TokenKindString,
                Column: line.index + line.Column,
        }

        line.nextRune()
        for line.notEnd() {
                ch := line.ch()

                if ch == terminator { break }
                if ch == '\\' {
                        err := line.getEscapeSequence(&token)
                        if err != nil { lexer.printError(line.index, err) }
                        continue
                }
                
                token.StringValue += string(ch)
                line.nextRune()
        }

        if terminator == '\'' {
                runes := []rune(token.StringValue)
                if len(runes) == 1 {
                        token.Value = runes[0]
                } else {
                        lexer.printError (
                                line.index - 1,
                                "rune literal must be one rune in size")
                        token.Value = '\000'
                }
        } else {
                token.Value = token.StringValue
        }

        line.addExisting(&token)
}

var escapeCodeMap = map[rune] rune {
        'a':  '\x07',
        'b':  '\x08',
        'f':  '\x0c',
        'n':  '\x0a',
        'r':  '\x0d',
        't':  '\x09',
        'v':  '\x0b',
        '\'': '\'',
        '"':  '"',
        '\\': '\\',
}

func (line *Line) getEscapeSequence (token *Token) (err error) {
        line.nextRune()
        ch := line.ch()

        code, exists := escapeCodeMap[ch]
        if exists {
                // simple escape sequence
                token.StringValue += string(code)
                
        } else if ch >= '0' && ch <= '7' {
                // octal escape sequence
                number := string(ch)
        
                line.nextRune()
                for line.notEnd() && len(number) < 3 {
                        ch = line.ch()
                        if ch < '0' || ch > '7' { break }

                        number += string(ch)
                        line.nextRune()
                }
                
                if len(number) < 3 {
                        return errors.New("octal escape sequence too short")
                }

                parsedNumber, _ := strconv.ParseInt(number, 8, 8)
                token.StringValue += string(parsedNumber)
                
        } else if ch == 'x' || ch == 'u' || ch == 'U' {
                // hexidecimal escape sequence
                want := 2
                if ch == 'u' { want = 4 }
                if ch == 'U' { want = 8 }
        
                number := ""

                line.nextRune()
                for line.notEnd() && len(number) < want {
                        ch = line.ch()
                        notLower := ch < 'a' || ch > 'f'
                        notUpper := ch < 'A' || ch > 'F'
                        notNum   := ch < '0' || ch > '9'
                        if notLower && notUpper && notNum { break }
                        
                        number += string(ch)
                        line.nextRune()
                }
                
                if len(number) < want {
                        return errors.New("hex escape sequence too short")
                }

                parsedNumber, _ := strconv.ParseInt(number, 16, want * 4)
                token.StringValue += string(parsedNumber)
                
        } else {
                return errors.New("invalid escape code \\" + string(ch))
        }

        return nil
}

func (lexer *Lexer) tokenizeSymbol () {
        line := lexer.line

        token := Token {
                Column: line.index + line.Column,
        }
        
        for line.notEnd() {
                ch := line.ch()

                // *breathes in*
                if ch >= 'a' && ch <= 'z' { break }
                if ch >= 'A' && ch <= 'Z' { break }
                if ch >= '0' && ch <= '9' {
                        // we may in fact be parsing a negative number!
                        if (token.StringValue == "-") {
                                lexer.tokenizeNumber(true)
                                return
                        }
                        break
                }
                if ch == ' '  { break }
                if ch == '"'  { break }
                if ch == '\'' { break }
                if ch == ':'  { break }
                if ch == '.'  { break }
                if ch == '['  { break }
                if ch == ']'  { break }
                if ch == '{'  { break }
                if ch == '}'  { break }
                // AHHHHHHHHHHHHHHHHHHHH

                token.StringValue += string(ch)

                line.nextRune()
        }

        token.Value = token.StringValue

        if token.StringValue == "---" {
                token.Kind = TokenKindSeparator
        } else if token.StringValue == "->" {
                token.Kind = TokenKindDirection
        } else {
                token.Kind = TokenKindSymbol
        }

        line.addExisting(&token)
}

func (lexer *Lexer) tokenizeMulti () {
        line := lexer.line

        token := Token {
                Column: line.index + line.Column,
        }
        
        for line.notEnd() {
                ch := line.ch()
                lowercase := ch >= 'a' && ch <= 'z'
                uppercase := ch >= 'A' && ch <= 'Z'
                number    := ch >= '0' && ch <= '9'

                if !lowercase && !uppercase && !number {
                        break
                }

                token.StringValue += string(ch)

                line.nextRune()
        }

        if validate.ValidatePermission(token.StringValue) {
                token.Kind  = TokenKindPermission
                token.Value = token.StringValue
        } else {
                token.Kind  = TokenKindName
                token.Value = token.StringValue
        }

        line.addExisting(&token)
}

func (lexer *Lexer) skipWhitespace () {
        line := lexer.line
        for line.notEnd() && line.ch() == ' ' {
                line.nextRune()
        }
}

/* nextLine advances the lexer to the next line, and returns whether or not the
 * end of the file was reached.
 */
func (lexer *Lexer) nextLine () (done bool, ignore bool, err error) {
        lexer.lineNumber ++
        
        done = lexer.lineNumber >= lexer.file.GetLength() - 1
        ignore = false

        line := &Line {
                Row: lexer.lineNumber,
        }
        
        lexer.line = line

        lineValue := lexer.file.GetLine(lexer.lineNumber)
        for i, ch := range lineValue {
                line.Indent = i
                if ch == '#' { return false, true, nil }
                if ch != ' ' { break }
        }
        line.Column = line.Indent
        line.runes = []rune(strings.TrimSpace(lineValue))
        line.EndColumn = line.Indent + len(line.runes)
        
        if line.Indent % 8 != 0 {
                line.Indent /= 8
                lexer.printError (
                        0,
                        "malformed indentation, use indentation size of 8",
                        "spaces")
                return false, true, nil
        }
        line.Indent /= 8
        return
}

func (lexer *Lexer) printWarning (column int, cause ...interface {}) {
        lexer.warnCount ++
        lexer.file.PrintWarning (
                column + lexer.line.Column,
                lexer.lineNumber, cause...)
}

func (lexer *Lexer) printError (column int, cause ...interface {}) {
        lexer.errorCount ++
        lexer.file.PrintError (
                column + lexer.line.Column,
                lexer.lineNumber, cause...)
}

func (lexer *Lexer) printFatal (err error) {
        lexer.errorCount ++
        lexer.file.PrintFatal("could not tokenize module -", err)
}

func (tokenKind TokenKind) ToString () (description string) {
        switch tokenKind {
                case TokenKindNone:          return "end of line";
                case TokenKindSeparator:     return "separator";
                case TokenKindDirection:     return "direction";
                case TokenKindPermission:    return "permission";
                case TokenKindInteger:       return "integer literal";
                case TokenKindSignedInteger: return "signed integer literal";
                case TokenKindFloat:         return "float literal";
                case TokenKindString:        return "string literal";
                case TokenKindRune:          return "rune literal";
                case TokenKindName:          return "name";
                case TokenKindSymbol:        return "symbol";
                case TokenKindColon:         return "colon";
                case TokenKindDot:           return "dot";
                case TokenKindLBracket:      return "left bracket";
                case TokenKindRBracket:      return "right bracket";
                case TokenKindLBrace:        return "left brace";
                case TokenKindRBrace:        return "right brace";

                default: return "BUG"
        }
}
