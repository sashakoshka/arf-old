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

        for {
                if parser.line.Indent <= parentIndent {
                        break
                        
                } else if parser.line.Indent == parentIndent + 1 {
                        // call
                        var statement *Statement
                        statement, err = parser.parseBodyFunctionCall (
                                parentIndent + 1,
                                block)
                        if err != nil { return }

                        if statement != nil {
                                block.items = append (
                                        block.items,
                                        BlockOrStatement {
                                                statement: statement,
                                        },
                                )
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
 * recursively, it may eat up more lines than one.
 */
func (parser *Parser) parseBodyFunctionCall (
        parentIndent int,
        // specifically for defining variables in
        parent *Block,
) (
        statement *Statement,
        err error,
) {
        statement = &Statement {}

        match := parser.expect (
                lexer.TokenKindLBracket,
                lexer.TokenKindName,
                lexer.TokenKindString,
                lexer.TokenKindSymbol)
        // we have no idea what the users intent with that was, so try to move
        // on to the next statement.
        if !match {
                parser.skipBodyFunctionCall(parentIndent, false)
                return nil, nil
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
                        err = parser.skipBodyFunctionCall (
                                parentIndent, bracketed)
                        return nil, err
                }
        }

        if parser.token.Kind == lexer.TokenKindString {
                // this statement calls a function of arbitrary name
                statement.external = true
                statement.externalCommand = parser.token.StringValue
        } else if parser.token.Kind == lexer.TokenKindSymbol {
                // this statement is an operator
                statement.command = Identifier {
                        trail: []string { parser.token.StringValue },
                }
                parser.nextToken()
        } else {
                // this statement calls a reachable function
                trail, worked, err := parser.parseIdentifier()
                        if err != nil { return nil, err }
                if !worked {
                        parser.skipBodyFunctionCall(parentIndent, bracketed)
                        return nil, nil
                }

                statement.command = Identifier { trail: trail }
        }

        complete := false
        for !complete {
                match = parser.expect (
                        lexer.TokenKindNone,
                        lexer.TokenKindLBracket,
                        lexer.TokenKindRBracket,
                        lexer.TokenKindColon,
                        lexer.TokenKindName,
                        lexer.TokenKindString,
                        lexer.TokenKindRune,
                        lexer.TokenKindInteger,
                        lexer.TokenKindSignedInteger,
                        lexer.TokenKindFloat)
                if !match {
                        err = parser.skipBodyFunctionCall (
                                parentIndent, bracketed)
                        return nil, err
                }

                argument := Argument {}

                switch parser.token.Kind {
                case lexer.TokenKindNone:
                        // if we have brackets, we can continue to parse the
                        // statement on the next line. if we don't, we are done
                        // parsing this statement.
                        if bracketed {
                                parser.nextLine()
                        } else {
                                complete = true
                        }
                        continue
                        
                case lexer.TokenKindName:
                        trail, worked, err := parser.parseIdentifier()
                        if err != nil { return nil, err }
                        if !worked {
                                parser.nextToken()
                                continue
                        }

                        argument.kind  = ArgumentKindIdentifier
                        argument.identifierValue = Identifier {
                                trail: trail,
                        }
                        break

                case lexer.TokenKindColon:
                        discardAfterParse := false
                
                        previousArgument := &statement.arguments [
                                len(statement.arguments) - 1]
                        
                        if previousArgument.kind != ArgumentKindIdentifier {
                                parser.printError (
                                        parser.token.Column,
                                        "type specifier may only follow an " +
                                        "identifier")
                                discardAfterParse = true
                        }

                        if len(previousArgument.identifierValue.trail) != 1 {
                                parser.printError (
                                        parser.token.Column,
                                        "cannot use member selection in " +
                                        "definition, name cannot have dots " +
                                        "in it")
                                discardAfterParse = true
                        }

                        parser.nextToken()
                        if !parser.expect (
                                lexer.TokenKindLBrace,
                                lexer.TokenKindName,
                        ) {
                                parser.nextToken()
                                continue
                        }

                        newArgument := Argument {
                                kind: ArgumentKindDefinition,
                                definitionValue: Definition {
                                        name: previousArgument.identifierValue,
                                },
                        }

                        what, worked, err := parser.parseType()
                        if err != nil { return nil, err }
                        if !worked {
                                parser.nextToken()
                                continue
                        }
                        
                        newArgument.definitionValue.what = what

                        if discardAfterParse {
                                statement.arguments = statement.arguments [
                                        :len(statement.arguments) - 1]
                        } else {
                                *previousArgument = newArgument
                        }
                        continue
                        
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
                        
                case lexer.TokenKindLBracket:
                        childStatement, err := parser.parseBodyFunctionCall (
                                parentIndent, parent)
                        if err != nil { return nil, err }
                        argument.kind = ArgumentKindStatement
                        argument.statementValue = childStatement
                        break
                        
                case lexer.TokenKindRBracket:
                        complete = true
                        continue
                        
                case lexer.TokenKindLBrace:
                        // TODO: get pointer notation
                        parser.nextToken()
                        continue
                }

                statement.arguments = append(statement.arguments, argument)
        }
        
        parser.nextLine()
        return
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
        // if the function isn't bracketed, we can just go on to the next line
        // without any worries
        if !bracketed {
                parser.nextLine()
                return
        }

        depth := 1
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
