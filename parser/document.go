package parser

import (
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
        isMember bool
        // TODO: convert to data
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

type Identifier struct {
        trail []string
}

type Type struct {
        name    Identifier
        points  *Type
        items   uint64
        mutable bool
}

type BlockOrStatement struct {
        block     *Block
        statement *Statement
}

type Block struct {
        datas map[string] *Data
        items []BlockOrStatement
}

type Statement struct {
        command   Identifier
        arguments []Argument

        external        bool
        externalCommand string
}

type Dereference struct {
        dereferences *Argument
        offset       uint64
}

type ArgumentKind int

const (
        ArgumentKindNone ArgumentKind = iota
        ArgumentKindStatement
        ArgumentKindIdentifier
        ArgumentKindDefinition
        ArgumentKindDereference
        ArgumentKindString
        ArgumentKindRune
        ArgumentKindInteger
        ArgumentKindSignedInteger
        ArgumentKindFloat
)

type Argument struct {
        kind ArgumentKind
        
        statementValue     *Statement
        identifierValue    *Identifier
        dereferenceValue   *Dereference
        stringValue        string
        runeValue          rune
        integerValue       uint64
        signedIntegerValue int64
        floatValue         float64
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

