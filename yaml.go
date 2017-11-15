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
            for key, _ := range m[serv].(map[interface {}]interface {}) {
                path := fmt.Sprintf(m[serv].(map[interface {}]interface {})[key].(map[interface {}]interface {})["path"].(string))
                function := fmt.Sprintf(m[serv].(map[interface {}]interface {})[key].(map[interface {}]interface {})["function"].(string))
                apitype := fmt.Sprintf(m[serv].(map[interface {}]interface {})[key].(map[interface {}]interface {})["apitype"].(string))
                methods := "GET,PUT,POST,DELETE,HEAD,OPTIONS,TRACE,CONNECT"
                if value, ok := m[serv].(map[interface {}]interface {})[key].(map[interface {}]interface {})["methods"]; ok {
                    methods = value.(string)
                }
                methodarray := strings.Split(methods,",")
                router.HandleFunc(path, fmap[function],apitype).Methods(methodarray)
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
