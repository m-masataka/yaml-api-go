package mux

import (
    "net/http"
)

type Route struct {
    handler http.Handler
    path    string
    methods []string
	children   []Route
}

func MatchVars(s string, req *http.Request) (bool, bool) {
    match, next, keys, values := MatchVarsRegexp(s, req.URL.Path)
    for i, key := range keys {
        ContextSet(req, key, values[i])
    }
    return match, next
}

func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
    if !r.MethodMatch(req) {
        match.MethodErr = false
    } else {
        match.MethodErr = true
    }
	m, n := MatchVars(r.path, req)
    if m {
		if n {
			var ret bool
			for i:= 0; i<len(r.children) ; i++ {
				ret = r.children[i].Match(req, match)
				if ret {
					return true
				}
			}
			return ret
		} else {
			match.Handler = r.handler
			match.Route = r
			return true
		}
	}
    return false
}

func (r *Route) Handler(handler http.Handler) *Route {
    r.handler = handler
    return r
}

func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
    return r.Handler(http.HandlerFunc(f))
}

func (r *Route) Methods(methods []string) *Route{
    for _, method := range methods {
        r.methods = append(r.methods, method)
    }
	return r
}

func (r *Route) MethodMatch(req *http.Request) bool {
    for _, method := range r.methods {
        if req.Method == method {
            return true
        }
    }
    return false
}

func (r *Route) RouteConf(tpl string) *Route {
    r.path    = tpl
    return r
}

func (r *Route) PathPrefix(tpl string) *Route {
	r.path = tpl
	return r
}

func (r *Route) Subroute(tpl string, f func(http.ResponseWriter, *http.Request)) *Route {
    route := Route{}
	r.children = append(r.children, route)
	r.children[len(r.children)-1].path = r.path + tpl
	return r.children[len(r.children)-1].Handler(http.HandlerFunc(f))
}
