package analyzer

import "github.com/sashakoshka/arf/parser"

type SectionSpecifier struct {
        module string
        name   string
}

type SemanticTree struct {
        typedefs  map[SectionSpecifier] *Typedef
        datas     map[SectionSpecifier] *Data
        functions map[SectionSpecifier] *Function
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

        points *Type
        items  uint64
}

type Data struct {
        
}

type Function struct {
        
}

func (tree *SemanticTree) GetTypedef (
        module string,
        name   string,
) (
        typedef *Typedef,
) {
        return tree.typedefs[SectionSpecifier { module: module, name: name }]
}

func (tree *SemanticTree) GetData (
        module string,
        name   string,
) (
        data *Data,
) {
        return tree.datas[SectionSpecifier { module: module, name: name }]
}

func (tree *SemanticTree) GetFunction (
        module string,
        name   string,
) (
        function *Function,
) {
        return tree.functions[SectionSpecifier { module: module, name: name }]
}

func (tree *SemanticTree) setTypedef (
        module string,
        name   string,
        typedef *Typedef,
) {
        tree.typedefs[SectionSpecifier { module: module, name: name }] =
                typedef
}

func (tree *SemanticTree) setData (
        module string,
        name   string,
        data *Data,
) {
        tree.datas[SectionSpecifier { module: module, name: name }] =
                data
}

func (tree *SemanticTree) setFunction (
        module string,
        name   string,
        function *Function,
) {
        tree.functions[SectionSpecifier { module: module, name: name }] =
                function
}
