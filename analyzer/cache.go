package analyzer

import "github.com/sashakoshka/arf/parser"

type CacheItem struct {
        module *parser.Module
        skimmed bool
}

var cache = make(map[string] CacheItem)

func AddCache (module *parser.Module, skimmed bool) {
        cache[module.GetPath()] = CacheItem {
                module:  module,
                skimmed: skimmed,
        }
}

func GetCache (path string, skimmedOK bool) (module *parser.Module) {
        cacheItem, exists := cache[path]
        if !exists { return }
        if !skimmedOK && cacheItem.skimmed { return }

        module = cacheItem.module
        return 
}
