package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type PortResponse struct {
	PID         int    `json:"pid"`
	ImageName   string `json:"imageName"`
	SessionName string `json:"sessionName"`
}

func portToPID(port int) []int {
	cmd := exec.Command("cmd", "/c", fmt.Sprintf("netstat -ano | findstr :%d", port))
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()
	var PIDs []int

	lines := strings.Split(out.String(), "\n")
	if len(lines) > 1 {
		for _, line := range lines {
			fields := strings.Fields(line)
			PIDstr := fields[len(fields)-1]
			PID, PIDerr := strconv.Atoi(PIDstr)
			if PIDerr != nil {
				continue
			}
			PIDs = append(PIDs, PID)
		}
	}

	return PIDs
}

func tasksFromPIDs(PIDs []int) []PortResponse {
	var tasks []PortResponse

	for i := 0; i < len(PIDs); i++ {
		pid := PIDs[i]
		cmd := exec.Command("cmd", "/c", fmt.Sprintf("tasklist /FO CSV | findstr :%d", pid))
		var out bytes.Buffer
		cmd.Stdout = &out
		_ = cmd.Run()

		lines := strings.Split(out.String(), "\n")
		if len(lines) > 1 {
			for _, line := range lines {
				var task PortResponse
				reader := csv.NewReader(strings.NewReader(line))
				fields, readError := reader.Read()
				if readError != nil {
					continue
				}
				task.PID = pid
				task.ImageName = fields[0]
				task.SessionName = fields[2]
				tasks = append(tasks, task)
			}
		}
	}

	return tasks
}

func handleStatus(writer http.ResponseWriter, request *http.Request) {
	portStr := request.URL.Query().Get("port")
	if portStr == "" {
		http.Error(writer, `{"error": "Missing 'port' parameter"}`, http.StatusBadRequest)
		return
	}

	portInt, portError := strconv.Atoi(portStr)
	if portError != nil {
		http.Error(writer, `{"error": "Invalid 'port' parameter"}`, http.StatusBadRequest)
		return
	}

	foundPIDs := portToPID(portInt)
	writer.Header().Set("Content-Type", "application/json")
	returnTasks := []PortResponse{}
	if len(foundPIDs) >= 1 {
		tasks := tasksFromPIDs(foundPIDs)
		if len(tasks) >= 1 {
			returnTasks = tasks
		}
	}
	json.NewEncoder(writer).Encode(returnTasks)
}

func main() {
	http.HandleFunc("/status", handleStatus)
	fmt.Println("PortLeash Daemon listening on port 4848...")
	http.ListenAndServe(":4848", nil)
}
