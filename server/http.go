package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 打印請求的方法、路徑和頭部訊息
	fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
	fmt.Println("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	// 根據請求方法進行不同處理
	if r.Method == http.MethodOptions {
		// 處理 OPTIONS 預檢查請求
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		// 讀取並顯示請求體
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		fmt.Printf("Body: %s\n", body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST request received and logged.\n"))
		return
	}

	// 處理其他請求方法
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed.\n"))
}

func main() {
	http.HandleFunc("/", handler)

	port := ":8080" // 設置伺服器端口
	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
