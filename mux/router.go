package mux

import (
	"net/http"
)

// Router registers routes to be matched and dispatches a handler.
type Router struct {
	NotFoundHandler http.Handler
	routes          []*Route
}

// RouteMatch stores information about a matched route.
type RouteMatch struct {
	Route     *Route
	Handler   http.Handler
	Vars      map[string]string
	MethodErr bool
}

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{}
}

// NotFoundDefault is set to Not Found function(that retrun if API not found) if user don't define own Not Found function.
func NotFoundDefault(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

// MethodErrFunc is set to Method not mutch function(that retrun if http method not mutch).
func MethodErrFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
}

// ServeHTTP dispatches the handler registered in the matched route.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

// Match attempts to match the given request against the router's registered routes.
func (r *Router) Match(req *http.Request, match *RouteMatch) bool {
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

// HandleFunc registers a new route with a matcher for the URL path.
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().RouteConf(path).HandlerFunc(f)
}

// NewRoute set Route instanse and return one.
func (r *Router) NewRoute() *Route {
	route := &Route{}
	r.routes = append(r.routes, route)
	return route
}

// RouteConf return Route instance.
func (r *Router) RouteConf(tpl string) *Route {
	return r.NewRoute().RouteConf(tpl)
}
