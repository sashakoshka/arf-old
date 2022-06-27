package parser

import "github.com/sashakoshka/arf/lexer"

/* parseMeta parses the metadata header of an arf file. This contains the module
 * name, and other miscellaneous fields such as author and license. Returns an
 * error if the file cannot be parsed further.
 */
func (parser *Parser) parseMeta () (err error) {
        for {
                if parser.line.Indent != 0 {
                        parser.printError(0, errBadIndent)
                        return nil
                }

                if !parser.expect (
                        lexer.TokenKindName,
                        lexer.TokenKindSeparator,
                ) { continue }

                if parser.token.Kind == lexer.TokenKindSeparator {
                        parser.nextLine()
                        return
                }

                key := parser.token.Value

                parser.nextToken()
                if !parser.expect (
                        lexer.TokenKindName,
                        lexer.TokenKindString,
                ) { continue }

                switch key {
                case "module":
                         // no-op, we already know this info
                        break
                case "author":
                        parser.module.author = parser.token.Value.(string)
                        break
                case "license":
                        parser.module.license = parser.token.Value.(string)
                        break
                case "require":
                        parser.module.imports = append (
                                parser.module.imports,
                                parser.token.Value.(string),
                        )
                        break
                default:
                        parser.printWarning(0, "uknown header directive")
                }

                // the rest of the line should be empty
                parser.nextToken()
                parser.expect()
                
                done := parser.nextLine()
                if done {
                        return errSurpriseEOF
                }
        }
        
        return
}
