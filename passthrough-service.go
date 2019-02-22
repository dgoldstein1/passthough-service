package main

import (
    "fmt"
    "net/http"
    "log"
    "io/ioutil"
    "github.com/davecgh/go-spew/spew"
)

func get(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    log.Println(spew.Sprint(r.Form))

    if r.Form.Get("url") == "" {
        fmt.Fprintf(w, "Error: must pass 'url'")
        return
    }

    res, err := http.Get(r.Form.Get("url"))
    if err != nil {
        fmt.Fprintf(w, err.Error())
        return
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Fprintf(w, err.Error())
        return
    }    
    fmt.Fprintf(w, string(body))
}

func main() {
    http.HandleFunc("/get", get) // set router

	log.Println("Serving on port 8080")
    err := http.ListenAndServe(":8080", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}