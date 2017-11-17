package mux

import (
	"regexp"
)

const (
	braceStart = byte(0x7b)
	braceEnd   = byte(0x7d)
	coron      = byte(0x3a)
	slash      = byte(0x2f)
)

func forwardPoint(path string, p int) (string, int) {
	reg := ""
	for {
		p++
		if path[p] == braceEnd {
			p++
			break
		}
		if !(p < len(path)) {
			break
		}
		reg = reg + string(path[p])
	}
	return reg, p
}

// MatchVarsRegexp judge the match of URL and PATH.
func MatchVarsRegexp(path string, url string) (bool, bool, []string, []string) {
	match := true
	next := false
	keys := []string{}
	values := []string{}
	u := 0
	p := 0
	for {
		if url[u] == path[p] {
		} else {
			if path[p] != braceStart {
				match = false
				break
			} else {
				key := ""
				reg := ""
				for {
					p++
					if path[p] == coron {
						keys = append(keys, key)
						reg, p = forwardPoint(path, p)
						break
					}
					if path[p] == braceEnd {
						keys = append(keys, key)
						p++
						break
					}
					if !(p < len(path)) {
						keys = append(keys, key)
						break
					}
					key = key + string(path[p])
				}
				value := ""
				for {
					if url[u] == slash {
						values = append(values, value)
						break
					}
					value = value + string(url[u])
					u++
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
		u++
		p++
		if p < len(path) && u < len(url) {
		} else if !(p < len(path)) && !(u < len(url)) {
			break
		} else if !(p < len(path)) && u < len(url) {
			next = true
			break
		} else {
			match = false
			break
		}
	}
	return match, next, keys, values
}
