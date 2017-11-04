package yamlapi

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

func YamlApi(buf []byte, fmap map[string]func(http.ResponseWriter, *http.Request)) {
    m := make(map[interface{}]interface{})
    err := yaml.Unmarshal(buf, &m)
    if err != nil {
        panic(err)
    }
    for key, _ := range m{
        path := fmt.Sprintf(m[key].(map[interface {}]interface {})["path"].(string))
        http.HandleFunc(path,fmap["cube"])
    }
    http.ListenAndServe(":9999",nil)
}
