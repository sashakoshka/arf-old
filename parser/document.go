package parser

import (
        "fmt"
        "errors"
)

type Module struct {
        name    string
        author  string
        license string
        imports []string

        functions map[string] *Function
        typedefs  map[string] *Typedef
        datas     map[string] *Data
}

type Function struct {
        self struct {
                name string
                what Type
        }
        
        name     string
        inputs   map[string] *Data
        outputs  map[string] *Data
        root     *Block

        modeInternal Mode
        modeExternal Mode
}

type Type struct {
        name   string
        points bool
        items  uint64
}

type Block struct {
        datas map[string] *Data
        calls []*Call
}

type Call struct {
        command string
        arguments []interface {}
}

type Data struct {
        name  string
        what  Type
        value []interface {}

        modeInternal Mode
        modeExternal Mode
}

type Mode int

const (
        ModeDeny Mode = iota
        ModeRead
        ModeWrite
)

type Typedef struct {
        name     string
        inherits Type

        members  []*Data

        modeInternal Mode
        modeExternal Mode
}

func decodePermission (value string) (internal Mode, external Mode) {
        if len(value) < 1 { return }
        switch value[0] {
                case 'n': internal = ModeDeny;  break
                case 'r': internal = ModeRead;  break
                case 'w': internal = ModeWrite; break
        }

        if len(value) < 2 { return }
        switch value[1] {
                case 'n': external = ModeDeny;  break
                case 'r': external = ModeRead;  break
                case 'w': external = ModeWrite; break
        }

        return
}

func (module *Module) Dump () {
        fmt.Println(":arf")
        fmt.Println("module", module.name)
        fmt.Println("author", "\"" + module.author + "\"")
        fmt.Println("license", "\"" + module.license + "\"")
        
        for _, item := range module.imports {
                fmt.Println("require", "\"" + item + "\"")
        }

        fmt.Println("---")

        for item, _ := range module.functions {
                fmt.Println("func", item)
        }

        for item, section := range module.typedefs {
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

                fmt.Print(" ", item + ":")
                if section.inherits.points { fmt.Print("{") }
                fmt.Print(section.inherits.name)
                if section.inherits.points {
                        fmt.Print(" ", section.inherits.items, "}")
                }
                fmt.Println()

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

        fmt.Print(" ", data.name + ":")
        if data.what.points { fmt.Print("{") }
        fmt.Print(data.what.name)
        if data.what.points {
                fmt.Print(" ", data.what.items, "}")
        }
        fmt.Println()
        
        for _, value := range data.value {
                printIndent(indent + 1)
                fmt.Println (value)
        }
}

func printIndent (level int) {
        for level > 0 {
                level --
                fmt.Print("        ")
        }
} 

func (module *Module) addData (data *Data) (err error) {
        if data == nil { return }
        _, exists := module.datas[data.name]
        if exists {
                return errors.New (
                        "data section " + data.name + "already exists")
        }

        module.datas[data.name] = data
        return nil
}

func (module *Module) addTypedef (typedef *Typedef) (err error) {
        if typedef == nil { return }
        _, exists := module.typedefs[typedef.name]
        if exists {
                return errors.New (
                        "data section " + typedef.name + "already exists")
        }

        module.typedefs[typedef.name] = typedef
        return nil
}

func (module *Module) addFunction (function *Function) (err error) {
        if function == nil { return }
        _, exists := module.functions[function.name]
        if exists {
                return errors.New (
                        "data section " + function.name + "already exists")
        }

        module.functions[function.name] = function
        return nil
}

