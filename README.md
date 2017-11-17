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
  - ``methods`` is restriction of method. see [Methods](#methods)
  - ``children`` define Hierachical api structure.

```
server:
  port: 9999
  notfound: notfound
api:
  app1:
    path: "/api/func1"
    function: "f1"
    methods:
      - GET
      - POST
  app2:
    path: "/api/func2/{var1}"
    function: "f2"
    methods:
      - POST
    children:
      app3:
        path: "/{var3}"
        function: "f3"
        methods:
          - GET
```

Implement function that is linked with API endpoint.

```
fmap := map[string]func(http.ResponseWriter, *http.Request){"f1":func1, "f2":func2}

yamlapigo.YamlApi(yamlfile, fmap)
```

## <a name="methods"> Methods
You can ristrict the method by this field.  
```
api:
  ...
  ...
  methods: [PUT, POST]
```
You can set ALL method in this field.
if you access this api by PUT method, you receive following message.
```
$ curl -X GET http://localhost:9999/api/func2/sss
Method not match
```
By default (If you don't define this field), ALL method is allowed.

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
For example
```
...
api:
  app1:
    path: "/api/func1"
    function: "f1"
  app2:
    path: "/api/func2/{var1}"
    function: "f2"
  app3:
    path: "/api/func3/{var1}/var/{var2}"
    function: "f3"
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

If you want to use regexp, set the vars field as ``{var:[regexp]}``.  
```
  app4:
    path: "/api/func4/{id:^[0-9]+$}
    function: "f4"
```

## Sample

sample.yml
```
multiplexer:
server:
  host: "www.example.com"
  port: 9999
  notfound: notfound
api:
  app1:
    path: "/api/func1"
    function: func1
    methods:
      - POST
      - GET
  app2:
    path: "/api/func2/{var1}"
    function: func2
    methods:
      - POST
      - GET
  app3:
    path: "/api/func3/{var1}/var/{var2}"
    function: func3
    methods:
      - PUT
  app4:
    path: "/api/func4/{id:^[0-9]+$}"
    function: func4
    methods:
      - GET
    children:
      app5:
        path: "/pic"
        function: func5
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
   var1 := yamlapigo.GetVars(r,"var1").(string)
   fmt.Fprintf(w,"func2: var1="+var1+"\n")
}

func func3(w http.ResponseWriter, r *http.Request) {
   var1 := yamlapigo.GetVars(r,"var1").(string)
   var2 := yamlapigo.GetVars(r,"var1").(string)
   fmt.Fprintf(w,"func3: var1=" + var1 +", var2="+var2+"\n")
}

func func4(w http.ResponseWriter, r *http.Request) {
   id := yamlapigo.GetVars(r,"id").(string)
   fmt.Fprintf(w,"func4: id="+id+"\n")
}

func func5(w http.ResponseWriter, r *http.Request) {
   id := yamlapigo.GetVars(r,"id").(string)
   fmt.Fprintf(w,"func5: id="+id+"\n")
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
        "func4":func4,
        "func5":func5,
        "notfound": notfound,
    }
    err = yamlapigo.YamlApi(buf,c)
    if err != nil {
        fmt.Println(err)
		return
    }
}
```
