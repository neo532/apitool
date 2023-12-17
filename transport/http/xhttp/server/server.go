package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var counter = make(map[string]int, 100)
var reqChan = make(chan string, 100)

func main() {
	graceClose()
	go count()

	port := flag.Int("port", 8500, "")
	flag.Parse()

	http.HandleFunc("/demo/v1/resource", handleGetById)
	http.HandleFunc("/demo/resource", handleGet)
	log(fmt.Sprintf("listen:%d,pid:%d", *port, os.Getpid()))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}
}

func count() {
	for s := range reqChan {
		counter[s]++
	}
}

func log(msg ...string) {
	m := make([]interface{}, 0, len(msg)+1)
	m = append(m, fmt.Sprintf(`"%+v"`, time.Now().Format(time.DateTime)))
	for _, s := range msg {
		m = append(m, s)
	}
	fmt.Println(m...)
}

func graceClose() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, os.Interrupt)
	go func() {
		<-c
		close(reqChan)
		time.Sleep(time.Second)
		fmt.Println(fmt.Sprintf("counter:\t%+v", counter))
		os.Exit(0)
	}()
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	log(r.RemoteAddr)
	reqChan <- r.RemoteAddr
	w.Write([]byte(fmt.Sprintf("reply: %s", r.RequestURI)))
}
func handleGetById(w http.ResponseWriter, r *http.Request) {
	log(r.RemoteAddr)
	reqChan <- r.RemoteAddr
	w.Write([]byte(fmt.Sprintf("reply: %s", r.RequestURI)))
}
