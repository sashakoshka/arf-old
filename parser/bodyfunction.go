package parser

import (
        "fmt"
        "github.com/sashakoshka/arf/lexer"
)

/* parseBodyFunction parses a function section.
 */
func (parser *Parser) parseBodyFunction () (section *Function, err error) {
        section = &Function {
                inputs:  make(map[string] *Data),
                outputs: make(map[string] *Data),
                root:    &Block {},
        }

        if !parser.expect(lexer.TokenKindPermission) {
                 return nil, parser.skipBodySection()
        }

        section.modeInternal,
        section.modeExternal = decodePermission(parser.token.StringValue)

        parser.nextToken()
        if !parser.expect(lexer.TokenKindName) {
                 return nil, parser.skipBodySection()
        }

        section.name = parser.token.StringValue

        parser.nextToken()
        if !parser.expect() { return nil, parser.skipBodySection() }
                
        parser.nextLine()
        if parser.endOfFile() || parser.line.Indent == 0 { return }

        // function arguments
        inHead := true
        for inHead {
                if !parser.expect (
                        lexer.TokenKindSeparator,
                        lexer.TokenKindSymbol,
                ) { return nil, parser.skipBodySection() }

                if parser.token.Kind == lexer.TokenKindSymbol {
                        switch parser.token.StringValue {
                        case "@":
                                section.self.name,
                                section.self.what,
                                _, err =  parser.parseDeclaration()
                                if err != nil { return nil, err }
                                section.isMember = true
                                break

                        case ">":
                                input := &Data {}
                                input.name,
                                input.what,
                                _, err =  parser.parseDeclaration()
                                if err != nil { return nil, err }

                                input.value,
                                _, err = parser.parseDefaultValues(1)
                                if err != nil { return nil, err }

                                section.inputs[input.name] = input
                                break
                        
                        case "<":
                                output := &Data {}
                                output.name,
                                output.what,
                                _, err =  parser.parseDeclaration()
                                if err != nil { return nil, err }

                                output.value,
                                _, err = parser.parseDefaultValues(1)
                                if err != nil { return nil, err }

                                section.outputs[output.name] = output
                                break

                        default:
                                parser.printError (
                                        parser.token.Column,
                                        "unknown argument type symbol '" +
                                        parser.token.StringValue + "',",
                                        "use either '@', '>', or '<'")
                                break
                        }
                }
                
                inHead = parser.token.Kind != lexer.TokenKindSeparator

                parser.nextLine()
                if parser.endOfFile() || parser.line.Indent == 0 { return }
        }

        // function body
        err = parser.parseBodyFunctionBlock(0, section.root)
        return
}

/* parseBodyFunctionBlock parses a block of function calls. This is done
 * recursively, so it will also parse sub-blocks.
 */
func (parser *Parser) parseBodyFunctionBlock (
        parentIndent int,
        parent *Block,
) (
        err error,
) {
        block := &Block {
                datas: make(map[string] *Data),
        }

        for {
                if parser.line.Indent <= parentIndent {
                        break
                } else if parser.line.Indent == parentIndent + 1 {
                        // call
                        err = parser.parseBodyFunctionCall (
                                parentIndent + 1,
                                block)
                        if parser.endOfFile() || err != nil { return }
                } else if parser.line.Indent == parentIndent + 2 {
                        // block
                        err = parser.parseBodyFunctionBlock (
                                parentIndent + 1,
                                block)
                        if parser.endOfFile() || err != nil { return }
                } else {
                        fmt.Println(parentIndent, parser.line.Indent)
                        parser.printError(0, errTooMuchIndent)
                        
                }
        }
        
        parent.items = append (parent.items, BlockOrStatement {
                block: block,
        })

        return
}

/* parseBodyFunctionCall parses a function call of a function body. This is done
 * recursively, it may eat up more lines than one.
 */
func (parser *Parser) parseBodyFunctionCall (
        parentIndent int,
        parent *Block,
) (
        err error,
) {
        statement := &Statement {}

        match := parser.expect (
                lexer.TokenKindLBracket,
                lexer.TokenKindName,
                lexer.TokenKindString,
                lexer.TokenKindSymbol)
        // we have no idea what the users intent with that was, so try to move
        // on to the next statement.
        if !match {
                parser.nextLine()
                return
        }
        
        bracketed := parser.token.Kind == lexer.TokenKindLBracket
        if bracketed {
                // that wasn't the function name, so try to get the function
                // name again.
                parser.nextToken()
                match = parser.expect (
                        lexer.TokenKindName,
                        lexer.TokenKindString,
                        lexer.TokenKindSymbol)
                if !match {
                        // TODO: skip to matching right bracket, or indentation
                        // drop.
                        parser.nextLine()
                        return
                }
        }

        if parser.token.Kind == lexer.TokenKindString {
                statement.external = true
        }

        // we now have enough information to do this, thank you for coming to my
        // ted talk
        statement.command = parser.token.StringValue
        
        parent.items = append (parent.items, BlockOrStatement {
                statement: statement,
        })
        
        parser.nextLine()
        return
}
