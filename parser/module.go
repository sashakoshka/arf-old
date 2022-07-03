package parser

import "errors"
import "github.com/sashakoshka/arf/lineFile"

type ErrorReporter interface {
        PrintWarning (cause ...interface { })
        PrintError   (cause ...interface { })
        PrintFatal   (cause ...interface { })
}

type Mode int

const (
        ModeDeny Mode = iota
        ModeRead
        ModeWrite
)

type Position struct {
        row    int
        column int
        file   *lineFile.LineFile
}

type Name string

type Permissions struct {
        modeInternal Mode
        modeExternal Mode 
}

type Module struct {
        Position
        Name
        
        path  string

        author  string
        license string
        imports []string

        functions map[string] *Function
        typedefs  map[string] *Typedef
        datas     map[string] *Data
}

type Typedef struct {
        Position
        Name
        Permissions
        Type

        members map[string] *Data
}

type Data struct {
        Position
        Name
        Permissions
        Type
        
        value []interface { }

        external bool
}

type Function struct {
        Position
        Name
        Permissions

        isMember bool
        self     string
        selfType string
        
        inputs   []string
        outputs  []string
        root     *Block

        external bool
}

type Identifier struct {
        Position
        
        trail []string
}

type Type struct {
        Position
        
        name    Identifier
        points  *Type
        items   uint64
        mutable bool
}

type BlockOrStatement struct {
        Position
        
        block     *Block
        statement *Statement
}

type Block struct {
        Position
        
        variables map[string] *Variable
        items     []BlockOrStatement
}

type Statement struct {
        Position
        
        command   Identifier
        arguments []Argument

        external        bool
        externalCommand string

        returnsTo []*Identifier
}

type Dereference struct {
        Position
        
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
        Position
        
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
        Position
        Name
        Type
        
        value []interface { }
}

func (position *Position) SetPosition (newPoisition Position) {
        *position = newPoisition
}

func (name *Name) SetName (newName string) {
        *name = Name(newName)
}

func (name *Name) GetName () (nameString string) {
        return string(*name)
}

func (permissions *Permissions) SetInternalPermission (mode Mode) {
        permissions.modeInternal = mode
} 

func (permissions *Permissions) SetExternalPermission (mode Mode) {
        permissions.modeExternal = mode
}

func (permissions *Permissions) GetInternalPermission () (mode Mode) {
        return permissions.modeInternal
}

func (permissions *Permissions) GetExternalPermission () (mode Mode) {
        return permissions.modeExternal
}

func (what *Type) SetType (newType Type) {
        *what = newType
}

func (what *Type) GetType () (Type) {
        return *what
}

func (what *Type) GetTypeData () (
        name    Identifier,
        points *Type,
        items   uint64,
        mutable bool,
) {
        return what.name,
                what.points,
                what.items,
                what.mutable
}

func (typedef *Typedef) GetMembers () (members map[string] *Data) {
        return typedef.members
}

func (typedef *Typedef) addMember (data *Data) (err error) {
        if data == nil { return }
        _, exists := typedef.members[data.GetName()]
        if exists {
                return errors.New (
                        "member " + data.GetName() + "already exists")
        }

        typedef.members[data.GetName()] = data
        return nil
}

/* addData adds a data section to a module
 */
func (module *Module) addData (data *Data) (err error) {
        if data == nil { return }
        _, exists := module.datas[data.GetName()]
        if exists {
                return errors.New (
                        "data section " + data.GetName() + "already exists")
        }

        module.datas[data.GetName()] = data
        return nil
}

/* addVariable adds a variable to a block
 */
func (block *Block) addVariable (variable *Variable) (worked bool) {
        if variable == nil { return }
        _, exists := block.variables[variable.GetName()]
        if exists {
                return false
        }

        block.variables[variable.GetName()] = variable
        return true
}

/* addTypedef adds a type section to a module
 */
func (module *Module) addTypedef (typedef *Typedef) (err error) {
        if typedef == nil { return }
        _, exists := module.typedefs[typedef.GetName()]
        if exists {
                return errors.New (
                        "data section " + typedef.GetName() + " already exists")
        }

        module.typedefs[typedef.GetName()] = typedef
        return nil
}

/* addFunction adds a function section to a module
 */
func (module *Module) addFunction (function *Function) (err error) {
        if function == nil { return }
        _, exists := module.functions[function.GetName()]
        if exists {
                return errors.New (
                        "data section " + function.GetName() +
                        " already exists")
        }

        module.functions[function.GetName()] = function
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
        author  string,
        license string,
        imports []string,
) {
        return module.author,
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
