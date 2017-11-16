package yamlapigo

import (
    "fmt"
    "time"
    "net/http"
    "strings"
    "log"
    "gopkg.in/yaml.v2"
    "github.com/m-masataka/yamlapigo/mux"
)

func test(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w,"test\n")
}

func getstringmap(m interface{}, key interface {}, s string) (string, error) {
	var ret string
	if value, ok := m.(map[interface{}]interface {})[key].(map[interface {}]interface {})[s]; ok {
		ret = value.(string)
	} else {
		return "", fmt.Errorf("Error Not Found")
	}
	return ret, nil
}

func parseapi(m interface{}, fmap map[string]func(http.ResponseWriter, *http.Request), route *mux.Route, router *mux.Router, counter int) error {
	counter ++
	for key, _ := range m.(map[interface{}]interface {}) {
        path, _ := getstringmap(m, key, "path")
        function, err := getstringmap(m, key, "function")
		if err != nil {
			return fmt.Errorf("function is not defined in yamlfile")
		} else if _, ok := fmap[function]; !ok {
			return fmt.Errorf("%s is not defined in map", function)
		}
        methods, err := getstringmap(m, key, "methods")
		if err != nil {
			methods = "GET,PUT,POST,DELETE,HEAD,OPTIONS,TRACE,CONNECT"
		}
        methodarray := strings.Split(methods,",")
		var r *mux.Route
		if counter == 1 {
			r = router.HandleFunc(path, fmap[function]).Methods(methodarray)
		} else {
			r = route.Subroute(path,fmap[function]).Methods(methodarray)
		}
		if value, ok := m.(map[interface {}]interface {})[key].(map[interface {}]interface {})["children"]; ok {
			err = parseapi(value, fmap, r, router, counter)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func YamlApi(buf []byte, fmap map[string]func(http.ResponseWriter, *http.Request)) error {
    m := make(map[interface{}]interface{})
    err := yaml.Unmarshal(buf, &m)
    if err != nil {
        return err
    }
    router := mux.NewRouter()
    port := ":80"
    for serv, _ := range m {
        switch serv {
        case "server":
            for key, _ := range m[serv].(map[interface {}]interface {}) {
                switch key {
                case "port":
                    port = fmt.Sprintf(":%d",m[serv].(map[interface {}]interface {})[key].(int))
                case "notfound":
                    function := fmt.Sprintf("%s",m[serv].(map[interface {}]interface {})[key].(string))
                    router.NotFoundHandler = http.HandlerFunc(fmap[function])
                default:
                    continue
                }
            }
        case "api":
            err := parseapi(m[serv], fmap, nil, router, 0)
			if err != nil {
				return err
			}
        default:
            continue
        }
    }
    srv := &http.Server{
            Handler:      router,
            Addr:         "0.0.0.0"+port,
            WriteTimeout: 15 * time.Second,
            ReadTimeout:  15 * time.Second,
    }
    log.Fatal(srv.ListenAndServe())
    return nil
}

func GetVars(r *http.Request, s string) interface{} {
	vars := mux.ContextGet(r,s)
	return vars
}
