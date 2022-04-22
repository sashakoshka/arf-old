package parser

import (
        "github.com/sashakoshka/arf/lexer"
)

/* parseBody parses the body of an arf file. This contains sections, which have
 * code in them. Returns an error if the file cannot be parsed further.
 */
func (parser *Parser) parseBody () (err error) {
        for !parser.endOfFile() {                
                if parser.line.Indent != 0 {
                        parser.printError(0, errBadIndent)
                        return nil
                }
                
                if !parser.expect(lexer.TokenKindName) { continue }

                switch parser.token.StringValue {
                case "data":
                        parser.nextToken()
                        section, err := parser.parseBodyData(0)
                        if err != nil { return err }
                        err = parser.module.addData(section)
                        if err != nil { parser.printError(5, err) }
                        break
                case "type":
                        parser.nextToken()
                        section, err := parser.parseBodyTypedef()
                        if err != nil { return err }
                        err = parser.module.addTypedef(section)
                        if err != nil { parser.printError(5, err) }
                        break
                case "func":
                        parser.nextToken()
                        section, err := parser.parseBodyFunction()
                        if err != nil { return err }
                        err = parser.module.addFunction(section)
                        if err != nil { parser.printError(5, err) }
                        break
                default:
                        parser.printError (
                                0, "unknown section type \"" +
                                parser.token.StringValue + "\"")
                        err = parser.skipBodySection()
                        if err != nil { return err }
                        break
                }
        }
        return
}

/* parseBodyData parses a data section.
 */
func (parser *Parser) parseBodyData (
        parentIndent int,
) (
        section *Data,
        err error,
) {
        section = &Data {}

        if !parser.expect(lexer.TokenKindPermission) {
                return nil, parser.skipBodySection()
        }

        section.modeInternal,
        section.modeExternal = decodePermission(parser.token.StringValue)

        worked := false
        section.name, section.what, worked, err = parser.parseDeclaration()
        if !worked {
                return nil, parser.skipBodySection()
        }

        section.value, worked, err = parser.parseDefaultValues(parentIndent)
        if err != nil { return nil, err }
        if !worked { return nil, parser.skipBodySection() }

        return
}

/* parseBodyTypedef parses a type definition section.
 */
func (parser *Parser) parseBodyTypedef () (section *Typedef, err error) {
        section = &Typedef {}

        if !parser.expect(lexer.TokenKindPermission) {
                 return nil, parser.skipBodySection()
        }

        section.modeInternal,
        section.modeExternal = decodePermission(parser.token.StringValue)

        worked := false
        section.name, section.inherits, worked, err = parser.parseDeclaration()
        if !worked {
                 return nil, parser.skipBodySection()
        }

        parser.nextToken()
        if !parser.expect() { return nil, parser.skipBodySection() }

        done := parser.nextLine()
        for {
                if done || parser.line.Indent == 0 { return }

                member, err := parser.parseBodyData(1)
                if err != nil { return nil, err }
                if member == nil { return nil, nil }

                section.members = append(section.members, member)
        }
}

/* parseBodyFunction parses a function section.
 */
func (parser *Parser) parseBodyFunction () (section *Function, err error) {
        section = &Function {
                inputs:  make(map[string] *Data),
                outputs: make(map[string] *Data),
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
                
        done := parser.nextLine()
        if done || parser.line.Indent == 0 { return }

        for {
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
                
                if parser.token.Kind == lexer.TokenKindSeparator { break }

                done = parser.nextLine()
                if done || parser.line.Indent == 0 { return }
        }
        
        for {
                
                
                done := parser.nextLine()
                if done || parser.line.Indent == 0 { return }
        }
}

/* skipBodySection ignores the rest of the current section of the body and moves
 * on to the next one.
 */
func (parser *Parser) skipBodySection () (err error ) {
        for {
                done := parser.nextLine()
                if done || parser.line.Indent == 0 { return }
        }
}


/* parseDefaultValues parses the default values of a variable.
 */
func (parser *Parser) parseDefaultValues (
        parentIndent int,
) (
        value []interface {},
        worked bool,
        err error,
) {
        for {
                parser.nextToken()
                if !parser.expect (
                        lexer.TokenKindNone,
                        lexer.TokenKindInt,
                        lexer.TokenKindFloat,
                        lexer.TokenKindString,
                        lexer.TokenKindRune,
                ) { return nil, false, nil }
                if parser.endOfLine() { break }
                value = append(value, parser.token.Value)
        }
        
        for {
                done := parser.nextLine()
                if done || parser.line.Indent <= parentIndent {
                        worked = true
                        return
                }

                for {
                        if !parser.expect (
                                lexer.TokenKindNone,
                                lexer.TokenKindInt,
                                lexer.TokenKindFloat,
                                lexer.TokenKindString,
                                lexer.TokenKindRune,
                        ) { return nil, false, nil }
                        if parser.endOfLine() { break }
                        value = append(value, parser.token.Value)
                        parser.nextToken()
                }
        }
}

/* parseDeclaration parses a variable declaration of the form name:Type or
 * name:{Type N}
 */
func (parser *Parser) parseDeclaration () (
        name string,
        what Type,
        worked bool,
        err error,
) {
        parser.nextToken()
        if !parser.expect(lexer.TokenKindName) { return }

        name = parser.token.StringValue

        parser.nextToken()
        if !parser.expect(lexer.TokenKindColon) { return }

        parser.nextToken()
        if !parser.expect (
                lexer.TokenKindName,
                lexer.TokenKindLBrace,
        ) { return }

        expectBrace := false
        
        if parser.token.Kind == lexer.TokenKindLBrace {
                parser.nextToken()
                if !parser.expect(lexer.TokenKindName) { return }

                expectBrace = true
                what.points = true
                what.items = 1
        }

        what.name = parser.token.StringValue

        if expectBrace {
                parser.nextToken()
                if !parser.expect (
                        lexer.TokenKindRBrace,
                        lexer.TokenKindInt,
                ) { return }

                if parser.token.Kind == lexer.TokenKindInt {
                        what.items = parser.token.Value.(uint64)
                        parser.nextToken()
                        if !parser.expect(lexer.TokenKindRBrace) { return }
                }

                worked = true
        }

        worked = true
        return
}
