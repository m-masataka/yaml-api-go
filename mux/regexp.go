package mux

import (
	"regexp"
)

func ForwardPoint( path string, p int) (string, int) {
    brace_end := []byte(`}`)
	reg := ""
	for {
		p ++
		if path[p] == brace_end[0] {
			p ++
			break
		}
		if !(p<len(path)){
			break
		}
		reg = reg + string(path[p])
	}
	return reg, p
}

func MatchVarsRegexp( path string, url string ) (bool, []string, []string){
    match := true
    brace_start := []byte(`{`)
    brace_end := []byte(`}`)
    coron := []byte(`:`)
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
				reg := ""
                for {
                    p ++
					if path[p] == coron[0] {
						keys = append(keys, key)
						reg, p = ForwardPoint(path, p)
						break
					}
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
				if reg != "" {
					re := regexp.MustCompile(reg)
					if !re.MatchString(value) {
						match = false
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


