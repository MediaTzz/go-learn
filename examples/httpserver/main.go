/*
 * @brief:
 * @Autor: tianzhize
 * @Date: 2022-04-11 21:06:56
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	_ "net/http/pprof"

	"github.com/golang/glog"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}

}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
	w.WriteHeader(200)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	// set response header
	for name, values := range r.Header {
		fmt.Println(name, "=", values)
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// get VERSION from env
	version := os.Getenv("VERSION")
	w.Header().Add("VERSION", version)
	fmt.Println("Version:", version)

	// get ip
	ip, err := GetIP(r)
	if err == nil {
		fmt.Println("client ip:", ip)
	}

	w.WriteHeader(200)
}

// GetIP returns request real ip.
func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}
