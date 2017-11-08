package yamlapigo

func MatchRegexp( path string, url string ) (bool, []string, []string){
    match := true
    brace_start := []byte(`{`)
    brace_end := []byte(`}`)
    slash := []byte(`/`)
    keys := []string{}
    values := []string{}
    u := 0
    p := 0
    for {
        if url[u] == path[p] {
        } else {
            if path[p] != brace_start[0] {
                match = false
                break
            } else {
                key := ""
                for {
                    p ++
                    if path[p] == brace_end[0] {
                        keys = append(keys, key)
                        p++
                        break
                    }
                    if !(p<len(path)){
                        keys = append(keys, key)
                        break
                    }
                    key = key + string(path[p])
                }
                value := ""
                for {
                    if url[u] == slash[0] {
                        values = append(values, value)
                        break
                    }
                    value = value + string(url[u])
                    u ++
                    if !(u < len(url)) {
                        values = append(values, value)
                        break
                    }
                }
            }
        }
        u ++
        p ++
        if p < len(path) && u < len(url) {
        }else if !(p < len(path)) && !(u < len(url)){
            break
        } else {
            match = false
            break
        }
    }
    return match, keys, values
}


