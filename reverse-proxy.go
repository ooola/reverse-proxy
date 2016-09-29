package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func NewReverseProxy(scheme, host string) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = scheme
		req.URL.Host = host
		req.Header.Set("X-SecretHeader", "GoProxy")
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	scheme := flag.String("scheme", "http", "scheme of proxy target - http or https")
	host := flag.String("url", "10.100.42.46:7171", "target url to proxy to")
	secret := flag.String("secret", "", "shared secret set in X-success-proxy-secret-token header sent to target")
	port := flag.String("port", ":8080", "port this program listens to")

	// environment variables override command line flags
	if len(os.Getenv("SUCCESS_SCHEME")) != 0 {
		*scheme = os.Getenv("SUCCESS_SCHEME")
	}
	if len(os.Getenv("SUCCESS_HOST")) != 0 {
		*host = os.Getenv("SUCCESS_HOST")
	}
	if len(os.Getenv("SUCCESS_SHARED_SECRET")) != 0 {
		*secret = os.Getenv("SUCCESS_SHARED_SECRET")
	}
	if len(os.Getenv("SUCCESS_PORT")) != 0 {
		*port = os.Getenv("SUCCESS_PORT")
	}

	flag.Parse()

	fmt.Println("scheme: %s", *scheme)
	fmt.Println("redirecting to: %s", *host)
	fmt.Println("server will listen on port: %s", *port)

	proxy := NewReverseProxy("http", "10.100.42.46:7171", *secret)
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
