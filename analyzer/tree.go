package analyzer

import "github.com/sashakoshka/arf/parser"

type SemanticTree struct {
        typedefs  map[string] *Typedef
        datas     map[string] *Data
        functions map[string] *Function
}

type Typedef struct {
        parser.Permissions
        parser.Name

        module    string
        inherits *Type

        members []*Data
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
