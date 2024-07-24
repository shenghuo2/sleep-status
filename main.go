package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.String("port", "8000", "port to serve on")
	host := flag.String("host", "0.0.0.0", "host to serve on")
	flag.Parse()

	if err := LoadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	http.HandleFunc("/status", StatusHandler)
	http.HandleFunc("/change", ChangeHandler)

	addr := fmt.Sprintf("%s:%s", *host, *port)
	fmt.Printf("Starting server on %s\n", addr)
	fmt.Printf("服务器启动于 %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		fmt.Printf("启动服务器错误: %v\n", err)
	}
}
