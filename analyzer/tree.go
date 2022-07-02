package analyzer

type SemanticTree struct {
        typedefs  map[string] *Typedef
        datas     map[string] *Data
        functions map[string] *Function
}

type Typedef struct {
        module    string
        name      string
        inherits *Type

        members []*Data

        modeInternal Mode
        modeExternal Mode
}

type Type struct {
        typedef *Typedef
        mutable  bool

        points bool
        items  uint64
}

type Data struct {
        
}

type Function struct {
        
}

type Mode int

const (
        ModeDeny Mode = iota
        ModeRead
        ModeWrite
)
