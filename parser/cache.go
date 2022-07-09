package parser

import "os"

type CacheItem struct {
        wasSkimmed bool
        module    *Module
}

var cache = make(map[string] CacheItem)

func cacheModule (module *Module, wasSkimmed bool) {
        cache[module.GetPath()] = CacheItem {
                wasSkimmed: wasSkimmed,
                module:     module,
        }
}

func GetModule (path string, skim bool) (item CacheItem, exists bool) {
        item, exists = cache[path]
        if exists {
                module, _, _, err := parse(path, skim)
                if err != nil { os.Exit(1) }
                cacheModule(module, skim)
        }
        return
}

func (item CacheItem) GetModule () (module *Module) {
        return item.module
}
