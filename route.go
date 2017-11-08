package yamlapigo

import (
    "net/http"
)

type Route struct {
    handler http.Handler
    path    string
    apitype string
}

func (r *Route) RouteConf(tpl string, apitype string) *Route {
    r.path    = tpl
    r.apitype = apitype
    return r
}

func MatchVars(s string, req *http.Request) bool {
    match, keys, values := MatchRegexp(s, req.URL.Path)
    for i, key := range keys {
        ContextSet(req, key, values[i])
    }
    return match
}

func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
    switch r.apitype {
    case "vars":
        if MatchVars(r.path, req) {
            match.Route = r
            match.Handler = r.handler
            return true
        }
    default:
    }
    if r.path == req.URL.Path {
        match.Route = r
        match.Handler = r.handler
        return true
    }else{
        return false
    }
}

func (r *Route) Handler(handler http.Handler) *Route {
    r.handler = handler
    return r
}

func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
    return r.Handler(http.HandlerFunc(f))
}
