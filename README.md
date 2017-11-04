# yamlapigo

## Installtion

```
go get github.com/m-masataka/yamlapigo
```

## Usage

Define your API with yaml file.  

- ``server`` field represent global configration.(Now it only port)
- ``api`` field represent api details. 
  - ``path`` is api endpoint path
  - ``function`` is function that is called by api endpoint.

```
server:
  port: 9999
api:
  app1:
    path: "/api/func1"
    function: "f1"
  app2:
    path: "/api/func2"
    function: "f2"
```

Implement function that is linked with API endpoint.

```
fmap := map[string]func(http.ResponseWriter, *http.Request){"f1":func1, "f2":func2}

yamlapigo.YamlApi(yamlfile,c)
```


## Sample

sample.yml
```
service:
  port: 9999
api:
  app1:
    path: "/api/func1"
    function: "f1"
  app2:
    path: "/api/func2"
    function: "f2"
```

main.go
```
package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "github.com/m-masataka/yamlapigo"
)

func func1(w http.ResponseWriter, r *http.Request) {
   fmt.Fprintf(w,"This is func1\n")
}

func func2(w http.ResponseWriter, r *http.Request) {
   fmt.Fprintf(w,"This is func2\n")
}

func main() {
    buf, err := ioutil.ReadFile("./sample.yml")
    fmap := map[string]func(http.ResponseWriter, *http.Request){"f1":func1, "f2":func2}
    if err != nil {
        return
    }
    err = yamlapigo.YamlApi(buf,fmap)
    fmt.Println(err)
}
```
