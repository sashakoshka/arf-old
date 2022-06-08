package parser

import (
        "fmt"
)

func (module *Module) Dump () {
        fmt.Println(":arf")
        fmt.Println("module", module.name)
        fmt.Println("author", "\"" + module.author + "\"")
        fmt.Println("license", "\"" + module.license + "\"")
        
        for _, item := range module.imports {
                fmt.Println("require", "\"" + item + "\"")
        }

        fmt.Println("---")

        for _, section := range module.functions {
                section.Dump()
        }

        for _, section := range module.typedefs {
                fmt.Print("type ")

                switch section.modeInternal {
                        case ModeDeny:  fmt.Print("n")
                        case ModeRead:  fmt.Print("r")
                        case ModeWrite: fmt.Print("w")
                }
                
                switch section.modeExternal {
                        case ModeDeny:  fmt.Print("n")
                        case ModeRead:  fmt.Print("r")
                        case ModeWrite: fmt.Print("w")
                }

                fmt.Println (
                        "", section.name +
                        ":" + section.inherits.ToString())

                for _, member := range section.members {
                        member.Dump(1)
                }
        }

        for _, section := range module.datas {
                section.Dump(0)
        }
}

func (data *Data) Dump (indent int) {
        printIndent(indent)
        if indent == 0 { fmt.Print("data ") }

        switch data.modeInternal {
                case ModeDeny:  fmt.Print("n")
                case ModeRead:  fmt.Print("r")
                case ModeWrite: fmt.Print("w")
        }
        
        switch data.modeExternal {
                case ModeDeny:  fmt.Print("n")
                case ModeRead:  fmt.Print("r")
                case ModeWrite: fmt.Print("w")
        }

        fmt.Println("", data.name + ":" + data.what.ToString())
        
        for _, value := range data.value {
                printIndent(indent + 1)
                fmt.Println (value)
        }
}

func (function *Function) Dump () {
        fmt.Print("func ")

        switch function.modeInternal {
                case ModeDeny:  fmt.Print("n")
                case ModeRead:  fmt.Print("r")
                case ModeWrite: fmt.Print("w")
        }
        
        switch function.modeExternal {
                case ModeDeny:  fmt.Print("n")
                case ModeRead:  fmt.Print("r")
                case ModeWrite: fmt.Print("w")
        }

        fmt.Println("", function.name)

        if function.isMember {
                fmt.Println (
                        "        @",
                        function.self.name + ":" +
                        function.self.what.ToString())
        }

        for _, input := range function.inputs {
                fmt.Println (
                        "        >",
                        input.name + ":" +
                        input.what.ToString())
                        
                for _, value := range input.value {
                        printIndent(2)
                        fmt.Println (value)
                }
        }

        for _, output := range function.outputs {
                fmt.Println (
                        "        <",
                        output.name + ":" +
                        output.what.ToString())
                        
                for _, value := range output.value {
                        printIndent(2)
                        fmt.Println (value)
                }
        }
        
        fmt.Println("        ---")

        if function.root != nil {
                function.root.Dump(1)
        }
}

func (block *Block) Dump (indent int) {
        for _, item := range block.items {
                if item.block != nil {
                        item.block.Dump(indent + 1)
                } else if item.statement != nil {
                        item.statement.Dump(indent)
                        fmt.Println()
                }
        }
}

func (statement *Statement) Dump (indent int) {
        printIndent(indent)
        fmt.Print("[")
        fmt.Print(statement.command.ToString())

        for _, argument := range statement.arguments {
                fmt.Println()
                if argument.kind == ArgumentKindStatement {
                        argument.statementValue.Dump(indent + 1)
                        continue
                }
                
                printIndent(indent + 1)
                switch argument.kind {
                case ArgumentKindIdentifier:
                        fmt.Print(argument.identifierValue.ToString())
                        break
                case ArgumentKindDefinition:
                        fmt.Print(argument.definitionValue.ToString())
                        break
                case ArgumentKindInteger:
                        fmt.Print(argument.integerValue)
                        break
                case ArgumentKindSignedInteger:
                        fmt.Print(argument.signedIntegerValue)
                        break
                case ArgumentKindFloat:
                        fmt.Print(argument.floatValue)
                        break
                case ArgumentKindString:
                        fmt.Print("\"", argument.stringValue, "\"")
                        break
                case ArgumentKindRune:
                        fmt.Print("'", string(argument.runeValue), "'")
                        break
                }
        }
        fmt.Print("]")
}

func (what *Type) ToString () (description string) {
        if what.points != nil {
                description += "{" + what.points.ToString()
                if what.items > 1 {
                        description += fmt.Sprint(" ", what.items)
                }
                description += "}"
        } else {
                description = what.name.ToString()
        }

        if what.mutable {
                description += ":mut"
        }
        
        return
}

func (identifier *Identifier) ToString () (description string) {
        if len(identifier.trail) < 1 { return "EMPTY.IDENTIFIER"}
        description = identifier.trail[0]
        if len(identifier.trail) < 2 { return }
        
        for _, item := range identifier.trail[1:] {
                description += "." + item
        }

        return
}

func (definition *Definition) ToString () (description string) {
        description += definition.name.ToString() + ":"
        description += definition.what.ToString()
        return
}

func printIndent (level int) {
        for level > 0 {
                level --
                fmt.Print("        ")
        }
}
