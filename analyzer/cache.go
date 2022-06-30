package analyzer

import "github.com/sashakoshka/arf/parser"

type CacheItem struct {
        // All modules in the cache are assumed to have been skimmed.
        module *parser.Module
}

var cache = make(map[string] CacheItem)

func CacheModule (module *parser.Module, skimmed bool) {
        cache[module.GetPath()] = CacheItem {
                module:  module,
        }
}

func GetCache (path string) (item CacheItem, exists bool) {
        item, exists = cache[path]
        return
}

func (item CacheItem) GetModule () (module *parser.Module) {
        return item.module
}
