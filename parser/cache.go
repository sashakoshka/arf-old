package parser

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

func GetCache (path string, needFull bool) (item CacheItem, exists bool) {
        item, exists = cache[path]
        return
}

func (item CacheItem) GetModule () (module *Module) {
        return item.module
}
