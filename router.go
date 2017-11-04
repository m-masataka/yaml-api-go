package yamlapi

import (
    "fmt"
    "net/http"
)

type Router struct {
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

func APINotFound(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "API Not Found")
}

func (r *Router) ServeHTTP( w http.ResponseWriter, req *http.Request ) {
    var handler http.Handler
    var match RouteMatch
    if r.Match(req, &match) {
        handler = match.Handler
    } else {
        handler = http.HandlerFunc(APINotFound)
    }
    handler.ServeHTTP(w, req)
}

func (r *Router) Match (req *http.Request, match *RouteMatch ) bool {
    for _, route := range r.routes {
        if route.Match(req, match) {
            return true
        }
    }
    return false
}

func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route {
    return r.NewRoute().Path(path).HandlerFunc(f)
}

func (r *Router) NewRoute() *Route {
    route := &Route{}
    r.routes = append(r.routes, route)
    return route
}

func (r *Router) Path(tpl string) *Route {
    return r.NewRoute().Path(tpl)
}
