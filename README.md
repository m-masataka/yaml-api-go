# yamlapigo

## Installtion

```
go get github.com/m-masataka/yamlapigo
```

## Usage

Define your API with yaml file.  

- ``server`` field represent global configration.
  - ``port``
  - ``notfound`` see [Not found](#notfound)
- ``api`` field represent api details. 
  - ``path`` is api endpoint path
  - ``function`` is function that is called by api endpoint.
  - ``apitype`` is type of api. see [Vars](#vars)

```
server:
  port: 9999
api:
  app1:
    path: "/api/func1"
    function: "f1"
    apitype: normal
  app2:
    path: "/api/func2/{var1}"
    function: "f2"
    apitype: vars
```

Implement function that is linked with API endpoint.

```
fmap := map[string]func(http.ResponseWriter, *http.Request){"f1":func1, "f2":func2}

yamlapigo.YamlApi(yamlfile, fmap)
```

## <a name="notfound"> Not found
You can set Not Found response in your program.
```
server:
  port: 9999
  notfound: notfoundfunc 
...
```
In your program...
```
func notfound(w http.ResponseWriter, r *http.Request) {
   fmt.Fprintf(w,"Notfound\n")
}
...
...
    c := map[string]func(http.ResponseWriter, *http.Request){
        "notfound": notfound,
    }
    err = yamlapigo.YamlApi(buf,c)
```

## <a name="vars"> Vars
You can use some valiables with API.  
You define apitype = ``vars`` and {valiable} in path.  
For example
```
...
api:
  app1:
    path: "/api/func1"
    function: "f1"
    apitype: normal
  app2:
    path: "/api/func2/{var1}"
    function: "f2"
    apitype: vars
  app3:
    path: "/api/func3/{var1}/var/{var2}"
    function: "f3"
    apitype: vars
```
You can get valiables by use ``yamlapigo.ContextGet()``.

```
func func1(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w,"func1\n")
}

func func2(w http.ResponseWriter, r *http.Request) {
   var1 := yamlapigo.ContextGet(r,"var1").(string)
   fmt.Fprintf(w,"func2: var1="+var1+"\n")
}

func func3(w http.ResponseWriter, r *http.Request) {
   var1 := yamlapigo.ContextGet(r,"var1").(string)
   var2 := yamlapigo.ContextGet(r,"var2").(string)
   fmt.Fprintf(w,"func3: var1=" + var1 +", var2="+var2+"\n")
}
```
result

```
$ root@ip-172-31-24-86:~# curl http://localhost:9999/api/func1
func1
$ curl http://localhost:9999/api/func2/ssssss
func2: var1=ssssss
$ curl http://localhost:9999/api/func3/12r/var/sdf
func3: var1=12r, var2=sdf
```

## Sample

sample.yml
```
server:
  port: 9999
  notfound: notfound
api:
  app1:
    path: "/api/func1"
    function: func1
    apitype: normal
  app2:
    path: "/api/func2/{var1}"
    function: func2
    apitype: vars
  app3:
    path: "/api/func3/{var1}/var/{var2}"
    function: func3
    apitype: vars
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
    fmt.Fprintf(w,"func1\n")
}

func func2(w http.ResponseWriter, r *http.Request) {
   var1 := yamlapigo.ContextGet(r,"var1").(string)
   fmt.Fprintf(w,"func2: var1="+var1+"\n")
}

func func3(w http.ResponseWriter, r *http.Request) {
   var1 := yamlapigo.ContextGet(r,"var1").(string)
   var2 := yamlapigo.ContextGet(r,"var1").(string)
   fmt.Fprintf(w,"func3: var1=" + var1 +", var2="+var2+"\n")
}

func notfound(w http.ResponseWriter, r *http.Request) {
   fmt.Fprintf(w,"Notfound\n")
}

func main() {
    buf, err := ioutil.ReadFile("./yamlapigo/sample.yml")
    if err != nil {
        fmt.Println(err)
        return
    }
    c := map[string]func(http.ResponseWriter, *http.Request){
        "func1":func1,
        "func2":func2,
        "func3":func3,
        "notfound": notfound,
    }
    err = yamlapigo.YamlApi(buf,c)
    if err != nil {
        fmt.Println(err)
        return
    }
}
```
