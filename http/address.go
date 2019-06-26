package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/golang-collections/collections/set"
)

//var All []string
var All = set.New()
var mu sync.Mutex

func AllHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var m []interface{}

	m = append(m, "备注: 本页面是填入地址引导.")
	m = append(m, "添加方式: 通过链接[http://ip:8000/add?addr=selfaddr]填入地址, ip填写本页面地址.")
	m = append(m, getAll())

	byt, _ := json.Marshal(m)
	fmt.Fprint(w, string(byt))

}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	addr := r.FormValue("addr")
	All.Insert(addr)
	//All = append(All, addr)

	byt, _ := json.Marshal(getAll())
	fmt.Fprint(w, string(byt))
	fmt.Printf("addr[%s] add\n", addr)
}
func getAll() []string {
	var r = []string{}
	All.Do(func(h interface{}) {
		r = append(r, h.(string))
	})
	return r
}

func main() {
	http.HandleFunc("/", AllHandler)
	http.HandleFunc("/add", AddHandler)
	http.ListenAndServe("0.0.0.0:8000", nil)
}
