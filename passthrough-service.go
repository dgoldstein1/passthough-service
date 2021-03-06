package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var client = &http.Client{}

type errorResponse struct {
	Msg   string
	RCode int
}

// returns error in query, otherwise 500 error
func error(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nReceived request from %v. Connection type: %v\n", r.Host, r.Proto)
	r.ParseForm()
	rCode := r.Form.Get("rCode")
	n, err := strconv.Atoi(rCode)
	if err != nil {
		fmt.Printf("Error getting rCode %s: %s \n", rCode, err.Error())
		fmt.Fprintf(w, "bad rCode: %s \n", err.Error())
		return
	}
	// return response with respones code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(n)
	json.NewEncoder(w).Encode(errorResponse{"example-error", n})

}

// makes ping request to PING_RESPONES_URL
func serve(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving to %s\n", os.Getenv("PING_RESPONSE_URL"))
	go func() {
		_, err := client.Get(os.Getenv("PING_RESPONSE_URL"))
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
	if os.Getenv("LOG_HEADERS") == "true" {
		fmt.Printf("request headers %v \n", spew.Sdump(r.Header))
	}
	if os.Getenv("LOG_BODY") == "true" {
		fmt.Printf("request body %v \n", spew.Sprint(r.Form))
	}
	// write response
	fmt.Fprintf(w, "Pong. mesh=%s \n", os.Getenv("MESH_ID"))
	// don't do anything if url is empty is false
	if os.Getenv("PING_RESPONSE_URL") == "" {
		return
	}
	// make new get request to PING_RESPONSE_URL
	if rand.Intn(20) == 10 {
		fmt.Println("Ball of out bounds")
		fmt.Fprintln(w, "Ball out of bounds")
		return
	}
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
		res, err := client.Get(os.Getenv("PING_RESPONSE_URL"))
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
	request, err := http.NewRequest("GET", r.Form.Get("url"), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if r.Form.Get("user_dn") != "" {
		fmt.Println("adding user_dn to header")
		request.Header.Set("user_dn", r.Form.Get("user_dn"))
	}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
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
	fmt.Fprintf(w, "%s\n", string(body))
}

// write base64 env variables to file, panic on error
func writebase64File(filename string, b64 string) {
	os.Remove(filename)
	fmt.Printf("writing file %s\n", filename)
	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/serve", serve)
	http.HandleFunc("/error", error)

	// parse port
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("no port set")
		os.Exit(1)
	}

	// load in certs
	if os.Getenv("USE_TLS") == "true" {
		certFile := os.Getenv("SERVER_CERT")
		keyFile := os.Getenv("SERVER_KEY")
		caFile := os.Getenv("SERVER_CA")

		// if read from env is true, write files from env variables
		if os.Getenv("READ_TLS_FROM_ENV") == "true" {
			fmt.Println("reading tls from env variables")
			certFile = "server.crt"
			keyFile = "server.key"
			caFile = "server-ca.crt"
			writebase64File(certFile, os.Getenv("SERVER_CERT"))
			writebase64File(keyFile, os.Getenv("SERVER_KEY"))
			writebase64File(caFile, os.Getenv("SERVER_CA"))
		}

		// Load certs
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client = &http.Client{Transport: transport}
	}

	log.Printf("Serving on port %s", port)

	addr := fmt.Sprintf(":%s", port)
	if os.Getenv("USE_HTTP2") == "true" {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("incoming request: %v\n", r.Proto)
			if !r.ProtoAtLeast(2, 0) {
				fmt.Fprint(w, "Http connections below 2.0 are not accepted\n")
			} else {
				fmt.Fprint(w, "Hello from the world of http2 <3\n")
			}
		})
		h2s := &http2.Server{}
		h1s := &http.Server{
			Addr:    addr,
			Handler: h2c.NewHandler(handler, h2s),
		}
		log.Fatal(h1s.ListenAndServe())
	} else {
		err := http.ListenAndServe(addr, nil) // set listen port
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}

}
