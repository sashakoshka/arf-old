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

        // if the type is braced, we have a pointer
        if parser.token.Kind == lexer.TokenKindLBrace {
                parser.nextToken()
                if !parser.expect(lexer.TokenKindName) { return }

                expectBrace = true
                what.points = true
                what.items = 1
        }

        // get the identifier of this declaration's type
        trail, worked, err := parser.parseIdentifier()
        if !worked || err != nil { return }

        what.name = Identifier { trail: trail }

        // if the type is a pointer, get its right brace
        if expectBrace {
                if !parser.expect (
                        lexer.TokenKindRBrace,
                        lexer.TokenKindInt,
                ) { return }

                // get the count, if there is one
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

/* parseIdentifier parses an identifier of the form name.name.name
 */
func (parser *Parser) parseIdentifier () (
        trail []string,
        worked bool,
        err error,
) {
        for {
                if !parser.expect(lexer.TokenKindName) {
                        worked = false
                        return
                }
                
                trail = append(trail, parser.token.StringValue)
                
                parser.nextToken()

                if 
                        parser.endOfLine() ||
                        parser.token.Kind != lexer.TokenKindDot {
                
                        worked = true
                        return
                }
                parser.nextToken()
        }

}
