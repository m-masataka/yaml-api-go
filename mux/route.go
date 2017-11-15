package mux

import (
    "net/http"
)

type Route struct {
    handler http.Handler
    path    string
    apitype string
    methods []string
}

func (r *Route) RouteConf(tpl string, apitype string) *Route {
    r.path    = tpl
    r.apitype = apitype
    return r
}

func MatchVars(s string, req *http.Request) bool {
    match, keys, values := MatchVarsRegexp(s, req.URL.Path)
    for i, key := range keys {
        ContextSet(req, key, values[i])
    }
    return match
}

func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
    if !r.MethodMatch(req) {
        match.MethodErr = false
    } else {
        match.MethodErr = true
    }
    switch r.apitype {
    case "vars":
        if MatchVars(r.path, req) {
            match.Route = r
            match.Handler = r.handler
            return true
        }
    default:
        if r.path == req.URL.Path {
            match.Route = r
            match.Handler = r.handler
            return true
        }else{
            return false
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

func (r *Route) Methods(methods []string) {
    for _, method := range methods {
        r.methods = append(r.methods, method)
    }
}

func (r *Route) MethodMatch(req *http.Request) bool {
    for _, method := range r.methods {
        if req.Method == method {
            return true
        }
    }
    return false
}
