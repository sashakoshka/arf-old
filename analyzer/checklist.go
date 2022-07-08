package analyzer

type Checklist struct {
        data map[string] interface { }
}

func (checklist *Checklist) IsDone (item string) (done bool) {
        if checklist.data == nil { return }
        _, done = checklist.data[item]
        return
}

func (checklist *Checklist) Do (item string) () {
        if checklist.data == nil {
                checklist.data = make(map[string] interface { })
        }
        checklist.data[item] = nil
}

func (checklist *Checklist) Undo (item string) () {
        if checklist.data == nil { return }
        delete(checklist.data, item)
}
