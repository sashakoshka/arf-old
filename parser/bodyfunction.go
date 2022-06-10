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
                        err = parser.parseBodyFunctionArgumentFor(section)
                        if err != nil { return }
                }
                
                inHead = parser.token.Kind != lexer.TokenKindSeparator

                parser.nextLine()
                if parser.endOfFile() || parser.line.Indent == 0 { return }
        }

        // function body
        block, err := parser.parseBodyFunctionBlock(0)
        section.root = block
        return
}

/* parseBodyFunctionArgumentFor parses a function argument for the specified
 * function.
 */
func (parser *Parser) parseBodyFunctionArgumentFor (
        section *Function,
) (
        err error,
) {
        switch parser.token.StringValue {
        case "@":
                section.self.name,
                section.self.what,
                _, err =  parser.parseDeclaration()
                if err != nil { return err }
                
                if section.self.what.points == nil {
                        parser.printError (
                                parser.token.Column,
                                "method reciever must be a",
                                "pointer")
                        break
                }
                
                if section.self.what.mutable {
                        parser.printError (
                                parser.token.Column,
                                "method reciever cannot be",
                                "mutable")
                        break
                }
                
                section.isMember = true
                break

        case ">":
                input := &Data {}
                input.name,
                input.what,
                _, err =  parser.parseDeclaration()
                if err != nil { return err }

                if input.what.mutable {
                        parser.printError (
                                parser.token.Column,
                                "function arguments cannot be",
                                "mutable")
                        break
                }

                section.inputs[input.name] = input
                if parser.endOfLine() { break}
                
                input.value,
                _, err = parser.parseDefaultValues(1)
                if err != nil { return err }
                break
        
        case "<":
                output := &Data {}
                output.name,
                output.what,
                _, err =  parser.parseDeclaration()
                if err != nil { return err }

                if output.what.mutable {
                        parser.printWarning (
                                parser.token.Column,
                                "you don't need to mark return",
                                "values as mutable, they will",
                                "be anyways")
                        break
                }
                
                output.what.mutable = true

                section.outputs[output.name] = output
                if parser.endOfLine() { break}
                
                output.value,
                _, err = parser.parseDefaultValues(1)
                if err != nil { return err }
                break

        default:
                parser.printError (
                        parser.token.Column,
                        "unknown argument type symbol '" +
                        parser.token.StringValue + "',",
                        "use either '@', '>', or '<'")
                break
        }

        return nil
}

/* parseBodyFunctionBlock parses a block of function calls. This is done
 * recursively, so it will also parse sub-blocks.
 */
func (parser *Parser) parseBodyFunctionBlock (
        parentIndent int,
) (
        block *Block,
        err error,
) {
        block = &Block {
                datas: make(map[string] *Data),
        }

        if (parser.line.Indent > 4) {
                parser.printWarning (
                        parser.token.Column,
                        "indentation level of",
                        parser.line.Indent,
                        "is difficult to read.",
                        "consider breaking up this function.")
        }

        for {
                if parser.line.Indent <= parentIndent {
                        break
                        
                } else if parser.line.Indent == parentIndent + 1 {
                        // call
                        var statement *Statement
                        var worked bool
                        statement, worked, err = parser.parseBodyFunctionCall (
                                parentIndent + 1,
                                block)
                        if err != nil || !worked { return }

                        block.items = append (
                                block.items,
                                BlockOrStatement {
                                        statement: statement,
                                },
                        )
                        
                        if parser.endOfLine() {
                                parser.nextLine()
                        }
                        
                } else if parser.line.Indent == parentIndent + 2 {
                        // block
                        var childBlock *Block
                        childBlock, err = parser.parseBodyFunctionBlock (
                                parentIndent + 1)
                        if err != nil { return }

                        block.items = append (block.items, BlockOrStatement {
                                block: childBlock,
                        })
                        
                } else {
                        fmt.Println(parentIndent, parser.line.Indent)
                        parser.printError(0, errTooMuchIndent)
                        
                }

                if parser.endOfFile() { return }
        }

        return
}

/* parseBodyFunctionCall parses a function call of a function body. This is done
 * recursively, and it may eat up more lines than one.
 */
func (parser *Parser) parseBodyFunctionCall (
        parentIndent int,
        // specifically for defining variables in
        parent *Block,
) (
        statement *Statement,
        worked bool,
        err error,
) {       
        statement = &Statement {}

        match := parser.expect (
                lexer.TokenKindLBracket,
                lexer.TokenKindName,
                lexer.TokenKindString,
                lexer.TokenKindSymbol)
        if !match {
                // we have no idea what the users intent with that was, so try
                // to move on to the next statement.
                parser.skipBodyFunctionCall(parentIndent, false)
                return nil, false, nil
        }

        // if the first token found was a bracket, this statement is wrapped in
        // brackets and we have to do some things differently
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
                        err = parser.skipBodyFunctionCall (
                                parentIndent, bracketed)
                        return nil, false, err
                }
        }

        if parser.token.Kind == lexer.TokenKindString {
                // this statement calls a function of arbitrary name
                statement.external = true
                statement.externalCommand = parser.token.StringValue
                parser.nextToken()
        } else if parser.token.Kind == lexer.TokenKindSymbol {
                // this statement is an operator
                statement.command = Identifier {
                        trail: []string { parser.token.StringValue },
                }
                parser.nextToken()
        } else {
                // this statement calls a reachable function
                trail, worked, err := parser.parseIdentifier()
                        if err != nil { return nil, false, err }
                if !worked {
                        parser.skipBodyFunctionCall(parentIndent, bracketed)
                        return nil, false, nil
                }

                statement.command = Identifier { trail: trail }
        }

        // get statement arguments
        complete := false
        for !complete {
                match = parser.expect (
                        lexer.TokenKindNone,
                        lexer.TokenKindLBracket,
                        lexer.TokenKindRBracket,
                        lexer.TokenKindLBrace,
                        lexer.TokenKindName,
                        lexer.TokenKindString,
                        lexer.TokenKindRune,
                        lexer.TokenKindInteger,
                        lexer.TokenKindSignedInteger,
                        lexer.TokenKindFloat)
                if !match {
                        err = parser.skipBodyFunctionCall (
                                parentIndent, bracketed)
                        return nil, false, err
                }

                if (parser.token.Kind == lexer.TokenKindNone) {
                        // if we have brackets, we can continue to parse the
                        // statement on the next line. if we don't, we are done
                        // parsing this statement.
                        if bracketed {
                                parser.nextLine()
                        } else {
                                complete = true
                        }
                        continue
                } else if (parser.token.Kind == lexer.TokenKindRBracket) {
                        complete = true
                        parser.nextToken()
                        continue
                }

                argument, worked, err := parser.parseArgument (
                        parentIndent, parent)
                if err != nil { return nil, false, err }
                if !worked { continue }

                statement.arguments = append (
                        statement.arguments,
                        argument)
        }
        
        return statement, true, nil
}

func (parser *Parser) parseArgument (
        parentIndent int,
        parent *Block,
) (
        argument Argument,
        worked bool,
        err error,
) {
        switch parser.token.Kind {
        case lexer.TokenKindLBracket:
                childStatement,
                worked,
                err := parser.parseBodyFunctionCall (parentIndent, parent)
                if err != nil { return argument, false, err }
                if !worked {
                        parser.nextToken()
                        return argument, false, nil
                }
                
                argument.kind = ArgumentKindStatement
                argument.statementValue = childStatement
                break
                
        case lexer.TokenKindLBrace:
                dereference,
                worked, err := parser.parseDereference(parentIndent, parent)
                if err != nil { return argument, false, err }
                if !worked {
                        parser.nextToken()
                        return argument, false, nil
                }

                argument.kind = ArgumentKindDereference
                argument.dereferenceValue = dereference
                break
                                
        case lexer.TokenKindName:
                trail, worked, err := parser.parseIdentifier()
                if err != nil { return argument, false, err }
                if !worked {
                        parser.nextToken()
                        return argument, false, nil
                }

                argument.kind  = ArgumentKindIdentifier
                argument.identifierValue = Identifier {
                        trail: trail,
                }

                // if there is no colon after this, this is not a
                // definition and we don't need to do anything else...
                if (parser.token.Kind != lexer.TokenKindColon) { break }
                // ... but if there is:

                if len(argument.identifierValue.trail) != 1 {
                        parser.printError (
                                parser.token.Column,
                                "cannot use member selection in " +
                                "definition, name cannot have dots " +
                                "in it")
                        return argument, false, nil
                }

                parser.nextToken()
                if !parser.expect (
                        lexer.TokenKindLBrace,
                        lexer.TokenKindName,
                ) {
                        parser.nextToken()
                        return argument, false, nil
                }

                argument.kind = ArgumentKindDefinition
                argument.definitionValue = Definition {
                        name: argument.identifierValue,
                }

                what, worked, err := parser.parseType()
                if err != nil { return argument, false, err }
                if !worked {
                        parser.nextToken()
                        return argument, false, nil
                }
                
                argument.definitionValue.what = what
                break
                
        case lexer.TokenKindString:
                argument.kind = ArgumentKindString
                argument.stringValue = parser.token.StringValue
                parser.nextToken()
                break
                
        case lexer.TokenKindRune:
                argument.kind = ArgumentKindRune
                argument.runeValue = parser.token.Value.(rune)
                parser.nextToken()
                break
                
        case lexer.TokenKindInteger:
                argument.kind = ArgumentKindInteger
                argument.integerValue = parser.token.Value.(uint64)
                parser.nextToken()
                break
                
        case lexer.TokenKindSignedInteger:
                argument.kind = ArgumentKindSignedInteger
                argument.signedIntegerValue = parser.token.Value.(int64)
                parser.nextToken()
                break
                
        case lexer.TokenKindFloat:
                argument.kind = ArgumentKindFloat
                argument.floatValue = parser.token.Value.(float64)
                parser.nextToken()
                break
        }

        return argument, true, nil
}

/* parseDereference parses a dereference of a value of the form {value N} where
 * N is an optional offset.
 */
func (parser *Parser) parseDereference (
        parentIndent int,
        parent *Block,
) (
        dereference Dereference,
        worked bool,
        err error,
) {
        if !parser.expect (lexer.TokenKindLBrace) {
                return dereference, false, nil
        }
        parser.nextToken()

        if (parser.token.Kind == lexer.TokenKindNone) {
                // if we are at the end of the line, just go on to the next one
                parser.nextLine()
        }
        
        if !parser.expect (
                lexer.TokenKindLBracket,
                lexer.TokenKindLBrace,
                lexer.TokenKindName,
                lexer.TokenKindString,
                lexer.TokenKindInteger,
        ) { return dereference, false, nil }
        
        argument, worked, err := parser.parseArgument (
                parentIndent, parent)
        if err != nil || !worked { return dereference, false, err }

        dereference.dereferences = &argument

        if !parser.expect (
                lexer.TokenKindRBrace,
                lexer.TokenKindInteger,
        ) {
                return dereference, false, nil
        }
        
        // get the count, if there is one
        if parser.token.Kind == lexer.TokenKindInteger {
                dereference.offset = parser.token.Value.(uint64)
                parser.nextToken()
                if !parser.expect(lexer.TokenKindRBrace) {
                        return dereference, false, nil
                }
        }
        
        parser.nextToken()

        return dereference, true, nil
}

/* skipBodyFunctionCall skips to the next body function call, or indentation
 * drop.
 */
func (parser *Parser) skipBodyFunctionCall (
        parentIndent int,
        bracketed bool,
) (
        err error,
) {
        depth := 1
        
        if !bracketed {
                depth --;
        }
        
        for {
                if parser.endOfFile() { return }
                if parser.endOfLine() { parser.nextLine() }

                if parser.token.Kind == lexer.TokenKindLBracket { depth ++ }
                if parser.token.Kind == lexer.TokenKindRBracket { depth -- }

                // if we drop out of the block or exit the statement, we are
                // done.
                if parser.line.Indent < parentIndent { break }
                if depth == 0 { break }
        
                parser.nextToken()
        }

        parser.nextLine()
        return
}
