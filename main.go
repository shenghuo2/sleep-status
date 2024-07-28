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

	// Output the current key at startup
	// 启动时输出当前的 key
	fmt.Printf("Current key: %s\n", ConfigData.Key)
	fmt.Printf("当前的 key: %s\n", ConfigData.Key)

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
