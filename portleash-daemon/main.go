package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type PortResponse struct {
	Port int    `json:"port"`
	PID  int    `json:"pid"`
	Name string `json:"name"`
}

func handleStatus(writer http.ResponseWriter, request *http.Request) {
	portStr := request.URL.Query().Get("port")
	if portStr == "" {
		http.Error(writer, `{"error": "Missing 'port' parameter"}`, http.StatusBadRequest)
		return
	}

	portInt, error := strconv.Atoi(portStr)
	if error != nil {
		http.Error(writer, `{"error": "Invalid 'port' parameter"}`, http.StatusBadRequest)
		return
	}

	response := PortResponse{Port: portInt, PID: 12140, Name: "portleash-daemon.exe"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

func main() {
	http.HandleFunc("/status", handleStatus)
	fmt.Println("PortLeash Daemon listening on port 4848...")
	http.ListenAndServe(":4848", nil)
}
