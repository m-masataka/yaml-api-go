package yamlapigo

import (
    "fmt"
    "time"
    "net/http"
    "log"
    "gopkg.in/yaml.v2"
)


func YamlApi(buf []byte, fmap map[string]func(http.ResponseWriter, *http.Request)) error {
    m := make(map[interface{}]interface{})
    err := yaml.Unmarshal(buf, &m)
    if err != nil {
        return err
    }
    router := NewRouter()
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
                router.HandleFunc(path, fmap[function],apitype)
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
