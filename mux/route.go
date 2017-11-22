package mux

import (
	"net/http"
	"errors"
)

var(
	ErrHostMismatch = errors.New("Host is not allowed")
)

// Route stores information to match a request and build URLs.
type Route struct {
	handler  http.Handler
	path     string
	methods  []string
	matchers []matcher
	children []Route
}

type matcher interface {
	Match(*http.Request, *RouteMatch) bool
}

// addMatcher adds a matcher to the route.
func (r *Route) AddMatcher(m matcher) *Route {
	r.matchers = append(r.matchers, m)
	return r
}

// MatchVars store values for requrst.
func MatchVars(s string, req *http.Request) (bool, bool) {
	match, next, keys, values := MatchVarsRegexp(s, req.URL.Path)
	for i, key := range keys {
		ContextSet(req, key, values[i])
	}
	return match, next
}

type hostmatcher []string

// AddHost add matcher against rwquest host.
func (r *Route) AddHost(host []string) *Route {
	var hm hostmatcher
	hm = hostmatcher(host)
	return r.AddMatcher(&hm)
}

// Match matches the route headers against the request.
func (h *hostmatcher) Match(req *http.Request, match *RouteMatch) bool {
	for _, value := range *h {
		if Trimcoron(req.Host) == string(value) {
			match.MatchErr = nil
			return true
		}
	}
	match.MatchErr = ErrHostMismatch
	return false
}

// Match matches the route against the request.
func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
	if !r.MethodMatch(req) {
		match.MethodErr = false
	} else {
		match.MethodErr = true
	}
	for _, m := range r.matchers {
		if matched := m.Match(req, match); !matched {
			return false;
		}
	}
	m, n := MatchVars(r.path, req)
	if m {
		if n {
			var ret bool
			for i := 0; i < len(r.children); i++ {
				ret = r.children[i].Match(req, match)
				if ret {
					return true
				}
			}
			return ret
		}
		match.Handler = r.handler
		match.Route = r
		return true
	}
	return false
}

// Handler sets a handler for the route.
func (r *Route) Handler(handler http.Handler) *Route {
	r.handler = handler
	return r
}

// HandlerFunc sets a handler function for the route.
func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
	return r.Handler(http.HandlerFunc(f))
}

// Methods adds a matcher for HTTP methods.
func (r *Route) Methods(methods []string) *Route {
	for _, method := range methods {
		r.methods = append(r.methods, method)
	}
	return r
}

// MethodMatch judge the method in request mutch or unmutch.
func (r *Route) MethodMatch(req *http.Request) bool {
	for _, method := range r.methods {
		if req.Method == method {
			return true
		}
	}
	return false
}

// RouteConf set the path to Route instance.
func (r *Route) RouteConf(tpl string) *Route {
	r.path = tpl
	return r
}

// Subroute set child Route to current Route instance.
func (r *Route) Subroute(tpl string, f func(http.ResponseWriter, *http.Request)) *Route {
	route := Route{}
	r.children = append(r.children, route)
	r.children[len(r.children)-1].path = r.path + tpl
	return r.children[len(r.children)-1].Handler(http.HandlerFunc(f))
}
