package yamlapigo

import (
    "fmt"
    "net/http"
)

type Router struct {
    NotFoundHandler http.Handler
    routes []*Route
}

type RouteMatch struct {
    Route    *Route
    Handler  http.Handler
    Vars     map[string]string
    MatchErr error
}

func NewRouter() *Router {
    return &Router{}
}

func NotFoundDefault(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "API Not Found")
}

func (r *Router) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
    var handler http.Handler
    var match RouteMatch
    if r.Match(req, &match) {
        handler = match.Handler
    } else {
        handler = http.HandlerFunc(NotFoundDefault)
    }
    handler.ServeHTTP(w, req)
}

func (r *Router) Match (req *http.Request, match *RouteMatch ) bool {
    for _, route := range r.routes {
        if route.Match(req, match) {
            return true
        }
    }
    if r.NotFoundHandler != nil {
        match.Handler = r.NotFoundHandler
        return true
    }
    return false
}

func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request), apitype string) *Route {
    return r.NewRoute().RouteConf(path, apitype).HandlerFunc(f)
}

func (r *Router) NewRoute() *Route {
    route := &Route{}
    r.routes = append(r.routes, route)
    return route
}

func (r *Router) RouteConf(tpl string, apitype string) *Route {
    return r.NewRoute().RouteConf(tpl,apitype)
}
