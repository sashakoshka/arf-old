package analyzer

import "fmt"
import "github.com/sashakoshka/arf/parser"

func (tree *SemanticTree) Dump () {
        fmt.Println("typedefs")
        for _, typedef := range tree.typedefs {
                typedef.Dump()
        }
}

func (typedef *Typedef) Dump () {
        printIndent(1)
        fmt.Println(typedef.module + "." + typedef.GetName(), "is an")

        switch typedef.GetInternalPermission() {
                case parser.ModeDeny:  fmt.Print("n")
                case parser.ModeRead:  fmt.Print("r")
                case parser.ModeWrite: fmt.Print("w")
        }
        
        switch typedef.GetExternalPermission() {
                case parser.ModeDeny:  fmt.Print("n")
                case parser.ModeRead:  fmt.Print("r")
                case parser.ModeWrite: fmt.Print("w")
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
        description += what.typedef.GetName()
        return
}

func printIndent (level int) {
        for level > 0 {
                level --
                fmt.Print("        ")
        }
}
