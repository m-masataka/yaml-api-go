package yamlapigo

import (
    "net/http"
)

type Route struct {
    handler http.Handler
    path    string
}

func (r *Route) Path(tpl string) *Route {
    r.path = tpl
    return r
}

func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
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
