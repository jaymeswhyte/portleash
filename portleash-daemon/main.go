package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type PortResponse struct {
	Port int    `json:"port"`
	PID  int    `json:"pid"`
	Name string `json:"name"`
}

func portToPID(port int) int {
	cmd := exec.Command("cmd", "/c", fmt.Sprintf("netstat -ano | findstr :%d", port))
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()

	lines := strings.Split(out.String(), "\n")
	if len(lines) > 1 {
		for _, line := range lines {
			fields := strings.Fields(line)
			PIDstr := fields[len(fields)-1]
			println(PIDstr)
			PID, PIDerr := strconv.Atoi(PIDstr)
			if PIDerr != nil {
				return 0
			}
			return PID
		}
	}

	return 0
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

	found := portToPID(portInt)

	if found >= 1 {
		response := PortResponse{Port: portInt, PID: 12140, Name: "portleash-daemon.exe"}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
	}
}

func main() {
	http.HandleFunc("/status", handleStatus)
	fmt.Println("PortLeash Daemon listening on port 4848...")
	http.ListenAndServe(":4848", nil)
}
