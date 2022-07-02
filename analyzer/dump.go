package analyzer

import "fmt"

func (tree *SemanticTree) Dump () {
        fmt.Println("typedefs")
        for _, typedef := range tree.typedefs {
                typedef.Dump()
        }
}

func (typedef *Typedef) Dump () {
        printIndent(1)
        fmt.Println(typedef.module + "." + typedef.name, "is an")

        switch typedef.modeInternal {
                case ModeDeny:  fmt.Print("n")
                case ModeRead:  fmt.Print("r")
                case ModeWrite: fmt.Print("w")
        }
        
        switch typedef.modeExternal {
                case ModeDeny:  fmt.Print("n")
                case ModeRead:  fmt.Print("r")
                case ModeWrite: fmt.Print("w")
        }
        
        fmt.Println(typedef.inherits.ToString())

        for _, member := range typedef.members {
                member.Dump(1)
        }
}

func (data *Data) Dump (indent int ) {
        
}

func (what *Type) ToString () (description string) {
        if what.mutable   { description += "mutable "    }
        if what.points    { description += "pointer to " }
        if what.items > 0 { description += fmt.Sprint(what.items, "") }
        description += what.typedef.name
        return
}

func printIndent (level int) {
        for level > 0 {
                level --
                fmt.Print("        ")
        }
}
