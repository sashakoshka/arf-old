package parser

import "fmt"

type Module struct {
        name string
        author string
        license string
        imports []string

        functions map[string] Function
        typedefs  map[string] Typedef
        datas     map[string] Data
}

type Function struct {
        self struct {
                name string
                what Type
        }
        
        name     string
        inputs   map[string] Data
        outputs  map[string] Data
        root     *Block
}

type Type struct {
        name   string
        points bool
}

type Block struct {
        datas map[string] Data
        calls []Call
}

type Call struct {
        command string
        arguments []interface {}
}

type Data struct {
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
        inherits []Type
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

        for item, _ := range module.typedefs {
                fmt.Println("type", item)
        }

        for item, _ := range module.datas {
                fmt.Println("data", item)
        }
}
