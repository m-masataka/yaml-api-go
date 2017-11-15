package mux

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
    MethodErr bool
}

func NewRouter() *Router {
    return &Router{}
}

func NotFoundDefault(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "API Not Found\n")
}

func MethodErrFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
    //fmt.Fprintf(w, "Method not match\n")
}

func (r *Router) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
    var handler http.Handler
    var match RouteMatch
    if r.Match(req, &match) {
        handler = match.Handler
        if !match.MethodErr {
            handler = http.HandlerFunc(MethodErrFunc)
        }
    } else {
        handler = http.HandlerFunc(NotFoundDefault)
    }
    defer ContextClear(req)
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

