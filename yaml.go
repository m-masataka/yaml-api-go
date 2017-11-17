package yamlapigo

import (
	"fmt"
	mux "github.com/m-masataka/yamlapigo/mux"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"strings"
	"time"

	gorilla "github.com/gorilla/mux"
)

const (
	//AllMethods define all http methods.
	allMethods = "GET,PUT,DELETE,PUT,HEAD,OPTIONS,TRACE,CONNECT,PATCH"
)

//Server include API server's imformation such as "port" "function" ,,,etc
type Server struct {
	port     string
	function string
}

//endpoint include API endpoint imformation
type endpoint struct {
	path     string
	function string
	methods  []string
}

type muxRouter struct {
	router *mux.Router
	route  *mux.Route
}

type gorillaRouter struct {
	router *gorilla.Router
	route  *gorilla.Route
}

//StartServer is the abstruct of server implementation
type StartServer interface {
	StartServer(map[interface{}]interface{}, map[string]func(http.ResponseWriter, *http.Request)) error
}

// NewServer returns a new Server instance.
func NewServer() *Server {
	return &Server{}
}

// newEndpoint returns a new Endpoint instance.
func newEndpoint() *endpoint {
	return &endpoint{}
}

// newMuxRouter returns a new muxRouter instance.
func newMuxRouter() *muxRouter {
	return &muxRouter{}
}

// newGorillaRouter returns a new muxgorillaRouter instance.
func newGorillaRouter() *gorillaRouter {
	return &gorillaRouter{}
}

func getstringmap(m interface{}, key interface{}, s string) (string, error) {
	var ret string
	if value, ok := m.(map[interface{}]interface{})[key].(map[interface{}]interface{})[s]; ok {
		ret = value.(string)
	} else {
		return "", fmt.Errorf("'%s' is Not Found in '%s'", s, key)
	}
	return ret, nil
}

func getmethodsarray(m interface{}, key interface{}, s string) ([]string, error) {
	var ret []string
	if value, ok := m.(map[interface{}]interface{})[key].(map[interface{}]interface{})[s]; ok {
		for _, v := range value.([]interface{}) {
			ret = append(ret, v.(string))
		}
	} else {
		return strings.Split(allMethods, ","), fmt.Errorf("'%s' is Not Found in '%s'", s, key)
	}
	return ret, nil
}

func (muxR *muxRouter) parseapi(m interface{}, fmap map[string]func(http.ResponseWriter, *http.Request), ep *endpoint, counter int) error {
	counter++
	var err error
	for key := range m.(map[interface{}]interface{}) {
		ep.path, err = getstringmap(m, key, "path")
		if err != nil {
			return err
		}
		ep.function, err = getstringmap(m, key, "function")
		if err != nil {
			return err
		} else if _, ok := fmap[ep.function]; !ok {
			return fmt.Errorf("%s is not defined in map", ep.function)
		}
		ep.methods, _ = getmethodsarray(m, key, "methods")
		if counter == 1 {
			muxR.route = muxR.router.HandleFunc(ep.path, fmap[ep.function]).Methods(ep.methods)
		} else {
			muxR.route = muxR.route.Subroute(ep.path, fmap[ep.function]).Methods(ep.methods)
		}
		if value, ok := m.(map[interface{}]interface{})[key].(map[interface{}]interface{})["children"]; ok {
			err = muxR.parseapi(value, fmap, ep, counter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) parseserv(m interface{}, fmap map[string]func(http.ResponseWriter, *http.Request)) error {
	s.port = ":80"
	for key := range m.(map[interface{}]interface{}) {
		switch key {
		case "port":
			s.port = fmt.Sprintf(":%d", m.(map[interface{}]interface{})[key].(int))
		case "notfound":
			s.function = fmt.Sprintf("%s", m.(map[interface{}]interface{})[key].(string))
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
		err = server.parseserv(m["server"], fmap)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("'server' is Not Found")
	}
	for serv := range m {
		switch serv {
		case "api":
			ep := newEndpoint()
			err := muxR.parseapi(m[serv], fmap, ep, 0)
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
		Addr:         "0.0.0.0" + server.port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
	return nil
}

func (gR *gorillaRouter) StartServer(m map[interface{}]interface{}, fmap map[string]func(http.ResponseWriter, *http.Request)) error {
	return fmt.Errorf("Sorry, gorilla is not implemented")
}

//YamlAPI unmarshal yaml data and Start API Server.
func YamlAPI(buf []byte, fmap map[string]func(http.ResponseWriter, *http.Request)) error {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(buf, &m)
	if err != nil {
		return err
	}
	if value, ok := m["multiplexer"]; ok && value == "gorilla" {
		router := newGorillaRouter()
		router.router = gorilla.NewRouter()
		return router.StartServer(m, fmap)
	}
	router := newMuxRouter()
	router.router = mux.NewRouter()
	return router.StartServer(m, fmap)
}

//GetVars return route variables for the current request.
func GetVars(r *http.Request, s string) interface{} {
	vars := mux.ContextGet(r, s)
	return vars
}
