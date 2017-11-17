package yamlapigo

import (
    "fmt"
    "time"
    "net/http"
    "log"
	"strings"
    "gopkg.in/yaml.v2"
    mux "github.com/m-masataka/yamlapigo/mux"

	gorilla "github.com/gorilla/mux"
)

const (
	ALL_METHODS = "GET,PUT,DELETE,PUT,HEAD,OPTIONS,TRACE,CONNECT,PATCH"
)
type Server struct {
	port string
	function string
}

type Api struct {
	path string
	function string
	methods []string
}

type muxRouter struct {
	router *mux.Router
	route  *mux.Route
}

type gorillaRouter struct {
	router *gorilla.Router
	route  *gorilla.Route
}

type AbstructStart interface {
	StartServer(map[interface{}]interface{}, map[string]func(http.ResponseWriter, *http.Request)) error
}

func NewServer() *Server {
	return &Server{}
}

func NewApi() *Api {
	return &Api{}
}

func NewmuxRouter() *muxRouter {
	return &muxRouter{}
}

func NewgorillaRouter() *gorillaRouter {
	return &gorillaRouter{}
}

func getstringmap(m interface{}, key interface {}, s string) (string, error) {
	var ret string
	if value, ok := m.(map[interface{}]interface {})[key].(map[interface {}]interface {})[s]; ok {
		ret = value.(string)
	} else {
		return "", fmt.Errorf("'%s' is Not Found in '%s'", s, key)
	}
	return ret, nil
}

func getmethodsarray(m interface{}, key interface {}, s string) ([]string, error) {
	var ret []string
	if value, ok := m.(map[interface{}]interface {})[key].(map[interface {}]interface {})[s]; ok {
		for _, v := range value.([]interface{}) {
			ret = append(ret, v.(string))
		}
	} else {
		return strings.Split(ALL_METHODS,","), fmt.Errorf("'%s' is Not Found in '%s'", s, key)
	}
	return ret, nil
}

func (muxR *muxRouter) parseapi(m interface{}, fmap map[string]func(http.ResponseWriter, *http.Request), api *Api, counter int) error {
	counter ++
	var err error
	for key, _ := range m.(map[interface{}]interface {}) {
		api.path, err = getstringmap(m, key, "path")
		if err != nil {
			return err
		}
		api.function, err = getstringmap(m, key, "function")
		if err != nil {
			return err
		} else if _, ok := fmap[api.function]; !ok {
			return fmt.Errorf("%s is not defined in map", api.function)
		}
		api.methods, _ = getmethodsarray(m, key, "methods")
		if counter == 1 {
			muxR.route = muxR.router.HandleFunc(api.path, fmap[api.function]).Methods(api.methods)
		} else {
			muxR.route = muxR.route.Subroute(api.path,fmap[api.function]).Methods(api.methods)
		}
		if value, ok := m.(map[interface {}]interface {})[key].(map[interface {}]interface {})["children"]; ok {
			err = muxR.parseapi(value, fmap, api, counter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) Parseserv(m interface{}, fmap map[string]func(http.ResponseWriter, *http.Request)) error {
    s.port = ":80"
	for key, _ := range m.(map[interface {}]interface {}) {
		switch key {
        case "port":
			s.port = fmt.Sprintf(":%d",m.(map[interface {}]interface {})[key].(int))
        case "notfound":
			s.function = fmt.Sprintf("%s",m.(map[interface {}]interface {})[key].(string))
        default:
            continue
        }
    }
	return nil
}

func (muxR *muxRouter) StartServer(m map[interface{}]interface{}, fmap map[string]func(http.ResponseWriter, *http.Request)) error {
	var err error
	server := NewServer()
	if _, ok := m["server"]; ok {
        err  = server.Parseserv(m["server"], fmap)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("'server' is Not Found")
	}
	for serv, _ := range m {
        switch serv {
        case "api":
			api := NewApi()
            err := muxR.parseapi(m[serv], fmap, api, 0)
			if err != nil {
				return err
			}
        default:
            continue
        }
    }
	if server.function != "" {
		muxR.router.NotFoundHandler = http.HandlerFunc(fmap[server.function])
	}
    srv := &http.Server{
            Handler:      muxR.router,
            Addr:         "0.0.0.0"+server.port,
            WriteTimeout: 15 * time.Second,
            ReadTimeout:  15 * time.Second,
    }
    log.Fatal(srv.ListenAndServe())
    return nil
}

func (gR *gorillaRouter) StartServer(m map[interface{}]interface{}, fmap map[string]func(http.ResponseWriter, *http.Request)) error {
	return fmt.Errorf("Sorry, gorilla is not implemented")
}

func YamlApi(buf []byte, fmap map[string]func(http.ResponseWriter, *http.Request)) error {
    m := make(map[interface{}]interface{})
    err := yaml.Unmarshal(buf, &m)
    if err != nil {
        return err
    }
	if value, ok := m["multiplexer"]; ok && value == "gorilla"{
		router := NewgorillaRouter()
		router.router = gorilla.NewRouter()
		return router.StartServer(m, fmap)
	} else {
		router := NewmuxRouter()
		router.router = mux.NewRouter()
		return router.StartServer(m, fmap)
	}
	return nil
}

func GetVars(r *http.Request, s string) interface{} {
	vars := mux.ContextGet(r,s)
	return vars
}
