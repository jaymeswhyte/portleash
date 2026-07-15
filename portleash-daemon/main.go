package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PortResponse struct {
	Port int    `json:"port"`
	PID  int    `json:"pid"`
	Name string `json:"name"`
}

func handleStatus(writer http.ResponseWriter, request *http.Request) {
	response := PortResponse{Port: 4848, PID: 12140, Name: "portleash-daemon.exe"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

func main() {
	http.HandleFunc("/status", handleStatus)
	fmt.Println("PortLeash Daemon listening on port 4848...")
	http.ListenAndServe(":4848", nil)
}
