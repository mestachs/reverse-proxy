package main

import (
	"encoding/json"
	"fmt"
	guuid "github.com/google/uuid"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"
)

var globalRequests sync.Map

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Get the port to listen on
func getListenAddress() string {
	port := getEnv("PORT", "3000")
	return "127.0.0.1:" + port
}

func handleStats(res http.ResponseWriter, req *http.Request) {
	stats := make(map[string]interface{})
	globalRequests.Range(func(k interface{}, v interface{}) bool {
		stats[k.(string)] = v
		return true
	})

	j, err := json.Marshal(&stats)
	if err != nil {
		panic(err)
	}
	res.Write(j)
}

type myTransport struct {
}

func (t *myTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	start := time.Now()
	key := guuid.New().String()
	globalRequests.Store(key, request.URL)

	response, err := http.DefaultTransport.RoundTrip(request)
	elapsed := time.Since(start)

	globalRequests.Delete(key)

	if err != nil {
		return nil, err //Server is not reachable. Server not working
	}

	log.Println("Response Time:", elapsed, request.URL)
	return response, err
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &myTransport{}
	proxy.ModifyResponse = func(r *http.Response) error {
		if r.StatusCode != http.StatusOK {
			log.Printf("request for %s got %v", r.Request.URL.String(), r.StatusCode)
		}
		if location, err := r.Location(); err == nil {

			// Turn it into a relative URL
			location.Scheme = ""
			location.Host = ""
			r.Header.Set("Location", location.String())
			return nil
		}
		return nil
	}

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	serveReverseProxy(getEnv("UPSTREAM", ""), res, req)
}

func handlerStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {

	// start server
	http.HandleFunc("/stats", handleStats)
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		panic(err)
	}
}
