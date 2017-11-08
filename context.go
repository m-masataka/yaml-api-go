package yamlapigo

import (
    "net/http"
    "sync"
)

var (
    mutex sync.RWMutex
    data = make(map[*http.Request]map[interface{}]interface{})
)

func ContextSet(r *http.Request, key interface{}, val interface{}) {
    mutex.Lock()
    if data[r] == nil {
        data[r] = make(map[interface{}]interface{})
    }
    data[r][key] = val
    mutex.Unlock()
}

func ContextGet(r *http.Request, key interface{} ) interface{} {
    var val interface{}
    mutex.RLock()
    if data[r] != nil {
        val = data[r][key]
        mutex.RUnlock()
        return val
    }
    mutex.RUnlock()
    return val
}
