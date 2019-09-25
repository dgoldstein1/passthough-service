package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received ping from %v. Connection type: %v\n", r.Host, r.Proto)
	fmt.Fprintf(w, "Pong from Mesh: %s", os.Getenv("MESH_ID"))
}

func get(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(spew.Sprint(r.Form))

	if r.Form.Get("url") == "" {
		fmt.Println("Error: must pass 'url'")
		fmt.Fprintf(w, "Error: must pass 'url'")
		return
	}

	pause := r.Form.Get("pause")
	if pause != "" {
		n, err := strconv.Atoi(pause)
		if err != nil {
			fmt.Printf("Error getting pause %s: %s \n", pause, err.Error())
			fmt.Fprintf(w, "bad pause: %s", err.Error())
			return
		}
		fmt.Printf("sleeping for %d s\n", n)
		time.Sleep(time.Duration(int32(n)) * time.Second)
		// time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
	}

	res, err := http.Get(r.Form.Get("url"))
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response: %s \n", err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Printf("Received response from '%s': %s \n", r.Form.Get("url"), string(body))
	fmt.Fprintf(w, string(body))
}

func main() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/ping", pong)

	// parse port
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("no port set")
		os.Exit(1)
	}

	log.Printf("Serving on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
