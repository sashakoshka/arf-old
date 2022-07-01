package parser

import "errors"
import "github.com/sashakoshka/arf/lineFile"

type ErrorReporter interface {
        PrintWarning (cause ...interface { })
        PrintError   (cause ...interface { })
        PrintFatal   (cause ...interface { })
}

type Position struct {
        row    int
        column int
        file   *lineFile.LineFile
}

type Module struct {
        where Position
        path  string

        name    string
        author  string
        license string
        imports []string

        functions map[string] *Function
        typedefs  map[string] *Typedef
        datas     map[string] *Data
}

type Function struct {
        where Position

        isMember bool
        self     string
        selfType string
        
        name     string
        inputs   []string
        outputs  []string
        root     *Block

        modeInternal Mode
        modeExternal Mode

        external bool
}

type Identifier struct {
        where Position
        
        trail []string
}

type Type struct {
        where Position
        
        name    Identifier
        points  *Type
        items   uint64
        mutable bool
}

type BlockOrStatement struct {
        where Position
        
        block     *Block
        statement *Statement
}

type Block struct {
        where Position
        
        variables map[string] *Variable
        items     []BlockOrStatement
}

type Statement struct {
        where Position
        
        command   Identifier
        arguments []Argument

        external        bool
        externalCommand string

        returnsTo []*Identifier
}

type Dereference struct {
        where Position
        
        dereferences *Argument
        offset       uint64
}

type ArgumentKind int

const (
        ArgumentKindNone ArgumentKind = iota
        ArgumentKindStatement
        ArgumentKindIdentifier
        ArgumentKindDereference
        ArgumentKindString
        ArgumentKindRune
        ArgumentKindInteger
        ArgumentKindSignedInteger
        ArgumentKindFloat
)

type Argument struct {
        where Position
        
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

type Variable struct {
        where Position

        name  string
        what  Type
        value []interface { }
}

type Data struct {
        where Position
        
        name  string
        what  Type
        value []interface { }

        modeInternal Mode
        modeExternal Mode

        external bool
}

type Mode int

const (
        ModeDeny Mode = iota
        ModeRead
        ModeWrite
)

type Typedef struct {
        where Position
        
        name     string
        inherits Type

        members  []*Data

        modeInternal Mode
        modeExternal Mode
}

/* addData adds a data section to a module
 */
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

/* addVariable adds a variable to a block
 */
func (block *Block) addVariable (variable *Variable) (worked bool) {
        if variable == nil { return }
        _, exists := block.variables[variable.name]
        if exists {
                return false
        }

        block.variables[variable.name] = variable
        return true
}

/* addTypedef adds a type section to a module
 */
func (module *Module) addTypedef (typedef *Typedef) (err error) {
        if typedef == nil { return }
        _, exists := module.typedefs[typedef.name]
        if exists {
                return errors.New (
                        "data section " + typedef.name + " already exists")
        }

        module.typedefs[typedef.name] = typedef
        return nil
}

/* addFunction adds a function section to a module
 */
func (module *Module) addFunction (function *Function) (err error) {
        if function == nil { return }
        _, exists := module.functions[function.name]
        if exists {
                return errors.New (
                        "data section " + function.name + " already exists")
        }

        module.functions[function.name] = function
        return nil
}

func (where *Position) PrintWarning (cause ...interface {}) {
        where.file.PrintWarning(where.column, where.row, cause...)
}

func (where *Position) PrintError (cause ...interface {}) {
        where.file.PrintError(where.column, where.row, cause...)
}

func (where *Position) PrintFatal (cause ...interface {}) {
        where.file.PrintFatal(cause)
}

/* GetMetadata returns the metadata fields of the module
 */
func (module *Module) GetMetadata () (
        name    string,
        author  string,
        license string,
        imports []string,
) {
        return module.name,
                module.author,
                module.license,
                module.imports
}

/* GetSections returns all sections within the module.
 */
func (module *Module) GetSections () (
        functions map[string] *Function,
        typedefs  map[string] *Typedef,
        datas     map[string] *Data,
) {
        return module.functions,
                module.typedefs,
                module.datas
}

/* GetPath returns the module's path on the filesystem.
 */
func (module *Module) GetPath () (path string) {
        return module.path
}
