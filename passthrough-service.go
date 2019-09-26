package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// makes ping request to PING_RESPONES_URL
func serve(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving to %s\n", os.Getenv("PING_RESPONSE_URL"))
	go func() {
		_, err := http.Get(os.Getenv("PING_RESPONSE_URL"))
		if err != nil {
			fmt.Printf("Error: %s \n", err.Error())
			fmt.Fprintf(w, "error making get request to %s: %s \n", os.Getenv("PING_RESPONSE_URL"), err.Error())
			return
		}
	}()

	fmt.Fprintf(w, "Served to %s", os.Getenv("PING_RESPONSE_URL"))
}

func pong(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("\nReceived ping from %v. Connection type: %v\n", r.Host, r.Proto)
	// make new get request to PING_RESPONSE_URL
	if rand.Intn(20) == 10 {
		fmt.Println("Ball of out bounds")
		fmt.Fprintln(w, "Ball out of bounds")
		return
	}
	// write response
	fmt.Fprintf(w, "Pong. mesh=%s \n", os.Getenv("MESH_ID"))
	// hit back in go routine
	go func() {
		fmt.Printf("hitting back to: %s\n", os.Getenv("PING_RESPONSE_URL"))
		pause := r.Form.Get("pause")
		if pause != "" {
			n, err := strconv.Atoi(pause)
			if err != nil {
				fmt.Printf("Error getting pause %s: %s \n", pause, err.Error())
				fmt.Fprintf(w, "bad pause: %s \n", err.Error())
				return
			}
			fmt.Printf("sleeping for %d s\n", n)
			time.Sleep(time.Duration(int32(n)) * time.Second)
		}
		res, err := http.Get(os.Getenv("PING_RESPONSE_URL"))
		if err != nil {
			fmt.Printf("Error: %s \n", err.Error())
			fmt.Fprintf(w, "error making get request to %s: %s \n", os.Getenv("PING_RESPONSE_URL"), err.Error())
			return
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response: %s \n", err.Error())
			fmt.Fprintf(w, "Error reading response %s \n", err.Error())
			return
		}
		fmt.Printf("Response: %s \n", string(body))
	}()
}

func get(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(spew.Sprint(r.Form))

	if r.Form.Get("url") == "" {
		fmt.Println("Error: must pass 'url'")
		fmt.Fprintf(w, "Error: must pass 'url' \n")
		return
	}

	pause := r.Form.Get("pause")
	if pause != "" {
		n, err := strconv.Atoi(pause)
		if err != nil {
			fmt.Printf("Error getting pause %s: %s \n", pause, err.Error())
			fmt.Fprintf(w, "bad pause: %s \n", err.Error())
			return
		}
		fmt.Printf("sleeping for %d s\n", n)
		time.Sleep(time.Duration(int32(n)) * time.Second)
		// time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
	}

	res, err := http.Get(r.Form.Get("url"))
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
		fmt.Fprintf(w, "error making get request to %s: %s \n", r.Form.Get("url"), err.Error())
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response: %s \n", err.Error())
		fmt.Fprintf(w, "Error reading response %s \n", err.Error())
		return
	}
	fmt.Printf("Received response from '%s': %s \n", r.Form.Get("url"), string(body))
	fmt.Fprintf(w, string(body)+"\n")
}

func main() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/serve", serve)

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
