package cache

type CacheItem {
        module *Module
        skimmed bool
}

var cache = make(map[string] CacheItem)

func AddCache (path string, module *Module, skimmed bool) {
        cache()
}

func GetCache (path string) (module *Module) {
        return cache
}
